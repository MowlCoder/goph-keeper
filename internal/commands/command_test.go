package commands

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandManager_RegisterCommand(t *testing.T) {
	commandManager := NewCommandManager()

	t.Run("add command", func(t *testing.T) {
		commandManager.RegisterCommand(
			"cmd",
			"cmd description",
			"tag",
			"cmd",
			func(args []string) error {
				return nil
			},
		)

		assert.Equal(t, "cmd", commandManager.commands["cmd"].Name)
		assert.Equal(t, "cmd", commandManager.commandsByTag["tag"][0].Name)
	})
}

func TestCommandManager_ExecCommandWithName(t *testing.T) {
	commandManager := NewCommandManager()
	err := errors.New("command error")
	cmdName := "cmd"

	commandManager.RegisterCommand(
		cmdName,
		"cmd description",
		"tag",
		"cmd",
		func(args []string) error {
			return err
		},
	)

	t.Run("exec command", func(t *testing.T) {
		assert.Equal(t, err, commandManager.ExecCommandWithName(cmdName, []string{}))
	})
}
