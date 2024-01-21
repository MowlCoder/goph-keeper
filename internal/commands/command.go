package commands

import (
	"errors"
	"fmt"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type command struct {
	Name        string
	Tag         string
	Description string
	Usage       string
	Exec        func(args []string) error
}

type CommandManager struct {
	commands      map[string]command
	commandsByTag map[string][]command
}

func NewCommandManager() *CommandManager {
	manager := CommandManager{
		commands:      make(map[string]command),
		commandsByTag: make(map[string][]command),
	}
	manager.initAppCommands()

	return &manager
}

func (m *CommandManager) RegisterCommand(
	name string,
	description string,
	tag string,
	usage string,
	exec func(args []string) error,
) {
	cmd := command{
		Name:        name,
		Tag:         tag,
		Description: description,
		Usage:       usage,
		Exec:        exec,
	}

	m.commands[name] = cmd
	m.commandsByTag[tag] = append(m.commandsByTag[tag], cmd)
}

func (m *CommandManager) ExecCommandWithName(name string, args []string) error {
	cmd, ok := m.commands[name]
	if !ok {
		fmt.Println("invalid command, type 'help' to get list of commands")
		return nil
	}

	err := cmd.Exec(args)

	if errors.Is(err, domain.ErrInvalidCommandUsage) {
		fmt.Println("Usage:", cmd.Usage)
		return nil
	}

	return err
}

func (m *CommandManager) initAppCommands() {
	m.RegisterCommand(
		"help",
		"list all of the app commands",
		"system",
		"help",
		m.helpCommand,
	)

	m.RegisterCommand(
		"quit",
		"exit from the app",
		"system",
		"quit",
		m.quitCommand,
	)
}

func (m *CommandManager) quitCommand(args []string) error {
	return domain.ErrQuitApp
}

func (m *CommandManager) helpCommand(args []string) error {
	fmt.Println("================================")

	for tag, commands := range m.commandsByTag {
		fmt.Printf("%s commands:\n", tag)
		for _, val := range commands {
			fmt.Println(fmt.Sprintf("    %s - %s. Usage: %s", val.Name, val.Description, val.Usage))
		}
	}

	fmt.Println("================================")

	return nil
}
