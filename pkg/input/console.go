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

func GetConsoleInput(placeholder string, defaultValue string) (string, error) {
	fmt.Print(placeholder)

	text, err := consoleReader.ReadString('\n')
	if err != nil {
		return defaultValue, nil
	}

	text = strings.Trim(text, "\n\r")

	if text == "" {
		return defaultValue, nil
	}

	return text, nil
}
