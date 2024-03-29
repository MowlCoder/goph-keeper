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
	clientServices "github.com/MowlCoder/goph-keeper/internal/services/client"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/internal/storage/file"
	"github.com/MowlCoder/goph-keeper/internal/utils/cryptor"
)

var (
	buildVersion string
	buildDate    string

	dataSecret string
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
		Timeout:   time.Second * time.Duration(clientConfig.ApiServerTimeout),
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

	userStoredDataStorage, err := file.InitFileStorage(path.Join(appDataDirPath, "user_stored_data.json"))
	if err != nil {
		log.Println(err)
		return
	}
	defer userStoredDataStorage.Close()

	sessionStorage, err := file.InitFileStorage(path.Join(appDataDirPath, "session.json"))
	if err != nil {
		log.Println(err)
		return
	}
	defer sessionStorage.Close()

	clientSession := session.NewClientSession(sessionStorage)
	userStoredDataAPI := api.NewUserStoredDataAPI(clientConfig.ServerBaseAddr, httpClient, clientSession)
	userAPI := api.NewUserAPI(clientConfig.ServerBaseAddr, httpClient, clientSession)

	userStoredDataRepository := fileRepositories.NewUserStoredDataRepository(userStoredDataStorage)

	dataCryptor := cryptor.New(dataSecret)

	userStoredDataService := clientServices.NewUserStoredDataService(userStoredDataRepository, dataCryptor)

	userHandler := handlers.NewUserHandler(httpClient, clientSession, userAPI)
	logPassHandler := handlers.NewLogPassHandler(clientSession, userStoredDataService)
	cardHandler := handlers.NewCardHandler(clientSession, userStoredDataService)
	textHandler := handlers.NewTextHandler(clientSession, userStoredDataService)
	fileHandler := handlers.NewFileHandler(clientSession, userStoredDataService)

	dataSyncer := clientsync.NewBaseSyncer(
		clientSession,
		userStoredDataAPI,
		userStoredDataService,
		userStoredDataRepository,
	)

	commandManager := commands.NewCommandManager()

	registerSystemCommands(commandManager, dataSyncer, appDataDirPath)
	registerUserCommands(commandManager, userHandler)
	registerLogPassCommands(commandManager, logPassHandler)
	registerCardCommands(commandManager, cardHandler)
	registerTextCommands(commandManager, textHandler)
	registerFileCommands(commandManager, fileHandler)

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
}

func registerSystemCommands(
	commandManager *commands.CommandManager,
	dataSyncer *clientsync.BaseSyncer,

	fileStoragePath string,
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
	commandManager.RegisterCommand(
		"storage",
		"show path to directory where data stores",
		"system",
		"storage",
		func(args []string) error {
			fmt.Println(fileStoragePath)
			return nil
		},
	)
	commandManager.RegisterCommand(
		"sync",
		"synchronize your data with server",
		"system",
		"sync [need auth]",
		dataSyncer.SyncCommandHandler,
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
		"login",
		userHandler.Authorize,
	)
	commandManager.RegisterCommand(
		"register",
		"create user",
		"user",
		"register",
		userHandler.Register,
	)
}

func registerLogPassCommands(
	commandManager *commands.CommandManager,
	logPassHandler *handlers.LogPassHandler,
) {
	commandManager.RegisterCommand(
		"lp-save",
		"save login password pair",
		"login password",
		"lp-save",
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
		"lp-upd",
		"update logpass pair by id",
		"login password",
		"lp-upd <id:int>",
		logPassHandler.UpdatePair,
	)
	commandManager.RegisterCommand(
		"lp-del",
		"delete login password pair by id",
		"login password",
		"lp-del <id:int>",
		logPassHandler.DeletePair,
	)
}

func registerCardCommands(
	commandManager *commands.CommandManager,
	cardHandler *handlers.CardHandler,
) {
	commandManager.RegisterCommand(
		"card-save",
		"save new card",
		"card",
		"card-save",
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
		"card-upd",
		"update card by id",
		"card",
		"card-upd <id:int>",
		cardHandler.UpdateCard,
	)
	commandManager.RegisterCommand(
		"card-del",
		"delete card by id",
		"card",
		"card-del <id:int>",
		cardHandler.DeleteCard,
	)
}

func registerTextCommands(
	commandManager *commands.CommandManager,
	textHandler *handlers.TextHandler,
) {
	commandManager.RegisterCommand(
		"text-save",
		"save new text",
		"text",
		"text-save",
		textHandler.AddText,
	)
	commandManager.RegisterCommand(
		"text-get",
		"get texts",
		"text",
		"text-get <page:int>",
		textHandler.GetTexts,
	)
	commandManager.RegisterCommand(
		"text-upd",
		"update text by id",
		"text",
		"text-upd <id:int>",
		textHandler.UpdateText,
	)
	commandManager.RegisterCommand(
		"text-del",
		"delete text by id",
		"text",
		"text-del <id:int>",
		textHandler.DeleteText,
	)
}

func registerFileCommands(
	commandManager *commands.CommandManager,
	fileHandler *handlers.FileHandler,
) {
	commandManager.RegisterCommand(
		"file-save",
		"save file",
		"file",
		"file-save",
		fileHandler.AddFile,
	)
	commandManager.RegisterCommand(
		"file-get",
		"get files",
		"file",
		"file-get <page:int>",
		fileHandler.GetFiles,
	)
	commandManager.RegisterCommand(
		"file-decrypt",
		"decrypt file to given directory",
		"file",
		"file-decrypt <id:int>",
		fileHandler.DecryptFile,
	)
	commandManager.RegisterCommand(
		"file-upd",
		"update file by id",
		"file",
		"file-upd <id:int>",
		fileHandler.UpdateFile,
	)
	commandManager.RegisterCommand(
		"file-del",
		"delete file by id",
		"file",
		"file-del <id:int>",
		fileHandler.DeleteFile,
	)
}
