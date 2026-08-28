package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Input struct {
	Reader *bufio.Reader
}

func NewInput(reader io.Reader) *Input {
	return &Input{Reader: bufio.NewReader(reader)}
}

func (i *Input) ReadLine(prompt string) (string, error) {
	if _, err := fmt.Fprint(io.Discard, prompt); err != nil {
		return "", err
	}
	line, err := i.Reader.ReadString('\n')
	return strings.TrimSpace(line), err
}

func (i *Input) ReadRequired(prompt string) (string, error) {
	value, err := i.ReadLine(prompt)
	if err != nil && err != io.EOF {
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("%s is required", prompt)
	}
	return value, nil
}

func SplitFields(line string) []string {
	parts := strings.Split(line, "|")
	for index := range parts {
		parts[index] = strings.TrimSpace(parts[index])
	}
	return parts
}
