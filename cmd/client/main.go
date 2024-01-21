package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/MowlCoder/goph-keeper/internal/api"
	"github.com/MowlCoder/goph-keeper/internal/clientsync"
	"github.com/MowlCoder/goph-keeper/internal/commands"
	"github.com/MowlCoder/goph-keeper/internal/commands/handlers"
	"github.com/MowlCoder/goph-keeper/internal/config"
	"github.com/MowlCoder/goph-keeper/internal/domain"
	fileRepositories "github.com/MowlCoder/goph-keeper/internal/repositories/file"
	"github.com/MowlCoder/goph-keeper/internal/services"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/internal/storage/file"
	"github.com/MowlCoder/goph-keeper/internal/utils/cryptor"
)

var (
	buildVersion string
	buildDate    string

	logPassSecret string
	cardSecret    string
)

func main() {
	err := godotenv.Load(".env.client")
	if err != nil {
		log.Println("No .env.client provided")
	}

	clientConfig := &config.Client{}
	clientConfig.Parse()

	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second * 20,
	}

	userCache, err := os.UserCacheDir()
	if err != nil {
		log.Println(err)
		return
	}

	appDataDirPath := path.Join(userCache, "goph-keeper")
	if err := os.Mkdir(appDataDirPath, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	clientSession := session.NewClientSession()
	logPassApi := api.NewLogPassAPI(clientConfig.ServerBaseAddr, httpClient, clientSession)
	cardApi := api.NewCardAPI(clientConfig.ServerBaseAddr, httpClient, clientSession)

	logPassFileStorage, err := file.InitFileStorage(path.Join(appDataDirPath, "logpass.json"))
	if err != nil {
		log.Println(err)
		return
	}

	cardFileStorage, err := file.InitFileStorage(path.Join(appDataDirPath, "card.json"))
	if err != nil {
		log.Println(err)
		return
	}

	logPassRepository := fileRepositories.NewLogPassRepository(logPassFileStorage)
	cardRepository := fileRepositories.NewCardRepository(cardFileStorage)

	logPassCryptor := cryptor.New(logPassSecret)
	cardCryptor := cryptor.New(cardSecret)

	logPassService := services.NewLogPassService(logPassRepository, logPassCryptor)
	cardService := services.NewCardService(cardRepository, cardCryptor)

	userHandler := handlers.NewUserHandler(httpClient, clientSession)
	logPassHandler := handlers.NewLogPassHandler(clientSession, logPassService)
	cardHandler := handlers.NewCardHandler(clientSession, cardService)

	logPassSyncer := clientsync.NewLogPassSyncer(
		clientSession,
		logPassApi,
		logPassService,
		logPassRepository,
	)
	cardSyncer := clientsync.NewCardSyncer(
		clientSession,
		cardApi,
		cardService,
		cardRepository,
	)

	commandManager := commands.NewCommandManager()

	registerSystemCommands(commandManager)
	registerUserCommands(commandManager, userHandler)
	registerLogPassCommands(commandManager, logPassHandler, logPassSyncer)
	registerCardCommands(commandManager, cardHandler, cardSyncer)

	go logPassSyncer.InfiniteSync(2 * time.Minute)
	go cardSyncer.InfiniteSync(2 * time.Minute)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Goph Keeper")
	fmt.Println("Type 'help' to get command list")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for {
			if !clientSession.IsAuth() {
				fmt.Print("(no auth) ")
			}

			fmt.Print("> ")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\n\r")
			parts := strings.Split(text, " ")

			err := commandManager.ExecCommandWithName(parts[0], parts[1:])

			if errors.Is(err, domain.ErrQuitApp) {
				sig <- syscall.SIGQUIT
				break
			}

			if err != nil {
				fmt.Println("executed with error -", err.Error())
			}
		}
	}()

	<-sig

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer shutdownCtxCancel()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			fmt.Println("Something goes wrong in exiting from app...forcing exit")
			os.Exit(1)
		}
	}()

	fmt.Println("Saving data...")

	if err := logPassSyncer.Sync(shutdownCtx); err != nil {
		fmt.Println("Error when saving data -", err.Error())
		return
	}

	fmt.Println("Successfully saved all data!")
}

func registerSystemCommands(
	commandManager *commands.CommandManager,
) {
	commandManager.RegisterCommand(
		"version",
		"get version of client binary",
		"system",
		"version",
		func(args []string) error {
			fmt.Println("==========================================")
			fmt.Printf("Version: %s\n", buildVersion)
			fmt.Printf("Build date: %s\n", buildDate)
			fmt.Println("==========================================")
			return nil
		},
	)
}

func registerUserCommands(
	commandManager *commands.CommandManager,
	userHandler *handlers.UserHandler,
) {
	commandManager.RegisterCommand(
		"login",
		"start user session",
		"user",
		"login <email:string> <password:string>",
		userHandler.Authorize,
	)
	commandManager.RegisterCommand(
		"register",
		"create user",
		"user",
		"register <email:string> <password:string>",
		userHandler.Register,
	)
}

func registerLogPassCommands(
	commandManager *commands.CommandManager,
	logPassHandler *handlers.LogPassHandler,
	logPassSyncer *clientsync.BaseSyncer,
) {
	commandManager.RegisterCommand(
		"lp-save",
		"save login password pair",
		"login password",
		"lp-save <login:string> <password:string> <source:string>",
		logPassHandler.AddPair,
	)
	commandManager.RegisterCommand(
		"lp-get",
		"get login password pairs",
		"login password",
		"lp-get <page:int>",
		logPassHandler.GetPairs,
	)
	commandManager.RegisterCommand(
		"lp-del",
		"delete login password pair by id",
		"login password",
		"lp-del <id:int>",
		logPassHandler.DeletePair,
	)
	commandManager.RegisterCommand(
		"lp-sync",
		"synchronize login password pairs with server",
		"login password",
		"lp-sync [need auth]",
		logPassSyncer.SyncCommandHandler,
	)
}

func registerCardCommands(
	commandManager *commands.CommandManager,
	cardHandler *handlers.CardHandler,
	cardSyncer *clientsync.BaseSyncer,
) {
	commandManager.RegisterCommand(
		"card-save",
		"save new card",
		"card",
		"card-save <number:string> <expiredAt:string> <cvv:string> <meta:string>",
		cardHandler.AddCard,
	)
	commandManager.RegisterCommand(
		"card-get",
		"get cards",
		"card",
		"card-get <page:int>",
		cardHandler.GetCards,
	)
	commandManager.RegisterCommand(
		"card-del",
		"delete card by id",
		"card",
		"card-del <id:int>",
		cardHandler.DeleteCard,
	)
	commandManager.RegisterCommand(
		"card-sync",
		"synchronize cards with server",
		"card",
		"card-sync [need auth]",
		cardSyncer.SyncCommandHandler,
	)
}
