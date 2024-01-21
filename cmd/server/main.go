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
	"github.com/MowlCoder/goph-keeper/internal/services"
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
	logPassCryptor := cryptor.New(serverConfig.LogPassSecret)
	cardCryptor := cryptor.New(serverConfig.LogPassSecret)

	authMiddleware := middleware.NewAuthMiddleware(tokenParser)

	userRepository := dbRepositories.NewUserRepository(dbPool)
	logPassRepository := dbRepositories.NewLogPassRepository(dbPool)
	cardRepository := dbRepositories.NewCardRepository(dbPool)

	userService := services.NewUserService(
		userRepository,
		passwordHasher,
	)
	logPassService := services.NewLogPassService(logPassRepository, logPassCryptor)
	cardService := services.NewCardService(cardRepository, cardCryptor)

	userHandler := handlers.NewUserHandler(userService, tokenGenerator)
	logPassHandler := handlers.NewLogPassHandler(logPassService)
	cardsHandler := handlers.NewCardHandler(cardService)

	server := &http.Server{
		Addr: serverConfig.HTTPAddr,
		Handler: makeHTTPRouter(
			authMiddleware,
			userHandler,
			logPassHandler,
			cardsHandler,
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
	logPassHandler *handlers.LogPassHandler,
	cardsHandler *handlers.CardHandler,
) http.Handler {
	router := chi.NewRouter()

	router.Route("/api/v1", func(apiRouter chi.Router) {
		apiRouter.Route("/user", func(userRouter chi.Router) {
			userRouter.Post("/register", userHandler.Register)
			userRouter.Post("/authorize", userHandler.Authorize)
		})

		apiRouter.Route("/logpass", func(logPassRouter chi.Router) {
			logPassRouter.Use(authMiddleware.Middleware)
			logPassRouter.Post("/", logPassHandler.AddNewPair)
			logPassRouter.Get("/", logPassHandler.GetMyPairs)
			logPassRouter.Delete("/", logPassHandler.DeleteBatchPairs)
		})

		apiRouter.Route("/cards", func(cardsRouter chi.Router) {
			cardsRouter.Use(authMiddleware.Middleware)
			cardsRouter.Post("/", cardsHandler.AddNewCard)
			cardsRouter.Get("/", cardsHandler.GetMyCards)
			cardsRouter.Delete("/", cardsHandler.DeleteBatchCards)
		})
	})

	return router
}
