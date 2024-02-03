package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var consoleReader *bufio.Reader

func init() {
	consoleReader = bufio.NewReader(os.Stdin)
}

// GetConsoleInput - display given placeholder and wait for input from user from standard input
func GetConsoleInput(placeholder string, defaultValue string) string {
	fmt.Print(placeholder)

	text, err := consoleReader.ReadString('\n')
	if err != nil {
		return defaultValue
	}

	text = strings.Trim(text, "\n\r")

	if text == "" {
		return defaultValue
	}

	return text
}
