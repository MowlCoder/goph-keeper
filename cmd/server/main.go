package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/MowlCoder/goph-keeper/internal/config"
	"github.com/MowlCoder/goph-keeper/internal/handlers"
	"github.com/MowlCoder/goph-keeper/internal/middleware"
	dbRepositories "github.com/MowlCoder/goph-keeper/internal/repositories/postgresql"
	serverServices "github.com/MowlCoder/goph-keeper/internal/services/server"
	"github.com/MowlCoder/goph-keeper/internal/storage/postgresql"
	"github.com/MowlCoder/goph-keeper/internal/utils/cryptor"
	"github.com/MowlCoder/goph-keeper/internal/utils/password"
	"github.com/MowlCoder/goph-keeper/internal/utils/token"
)

func main() {
	err := godotenv.Load(".env.server")
	if err != nil {
		log.Println("No .env.server provided")
	}

	serverConfig := &config.Server{}
	serverConfig.Parse()

	dbPool, err := postgresql.InitPool(serverConfig.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	if err := postgresql.RunMigrations(serverConfig.DatabaseDSN); err != nil {
		log.Fatal(err)
	}

	passwordHasher := password.NewHasher()
	tokenGenerator := token.NewGenerator()
	tokenParser := token.NewParser()

	dataCryptor := cryptor.New(serverConfig.DataSecretKey)

	authMiddleware := middleware.NewAuthMiddleware(tokenParser)

	userRepository := dbRepositories.NewUserRepository(dbPool)
	userStoredDataRepository := dbRepositories.NewUserStoredDataRepository(dbPool)

	userService := serverServices.NewUserService(
		userRepository,
		passwordHasher,
	)
	userStoredDataService := serverServices.NewUserStoredDataService(userStoredDataRepository, dataCryptor)

	userHandler := handlers.NewUserHandler(userService, tokenGenerator)
	userStoredDataHandler := handlers.NewUserStoredDataHandler(userStoredDataService)

	server := &http.Server{
		Addr: serverConfig.HTTPAddr,
		Handler: makeHTTPRouter(
			authMiddleware,
			userHandler,
			userStoredDataHandler,
		),
	}

	log.Println("goph-keeper server is running on", serverConfig.HTTPAddr)

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	log.Println("goph-keeper server started shutdown process")

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCtxCancel()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Fatal("graceful shutdown timed out... forcing exit")
		}
	}()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}

	log.Println("goph-keeper server shutdown process successfully completed")
}

func makeHTTPRouter(
	authMiddleware *middleware.AuthMiddleware,

	userHandler *handlers.UserHandler,
	userStoredDataHandler *handlers.UserStoredDataHandler,
) http.Handler {
	router := chi.NewRouter()

	router.Route("/api/v1", func(apiRouter chi.Router) {
		apiRouter.Route("/user", func(userRouter chi.Router) {
			userRouter.Post("/register", userHandler.Register)
			userRouter.Post("/authorize", userHandler.Authorize)
		})

		apiRouter.Route("/data", func(dataRouter chi.Router) {
			dataRouter.Use(authMiddleware.Middleware)
			dataRouter.Get("/{type}", userStoredDataHandler.GetOfType)
			dataRouter.Post("/{type}", userStoredDataHandler.Add)
			dataRouter.Get("/", userStoredDataHandler.GetUserAll)
			dataRouter.Delete("/", userStoredDataHandler.DeleteBatch)
		})
	})

	return router
}
