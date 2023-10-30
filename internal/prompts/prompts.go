package prompts

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alex-guoba/tbd/pkg/logger"
)

type Prompts struct {
	ctx context.Context
}

func New(ctx context.Context) *Prompts {
	return &Prompts{
		ctx: ctx,
	}
}

// Read prompt from pipe(like other file or shell cmd outputs), used to statics and anaylyze data
func (*Prompts) ReadFromPipe() (bool, string) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		logger.Error(err)
		return false, ""
	}

	if stat.Mode()&os.ModeNamedPipe == 0 {
		// not pipe
		return false, ""
	} else {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			logger.Error(err)
			return true, ""
		}
		str := string(stdin)
		return true, strings.TrimSuffix(str, "\n")
	}
}

// TODO: support multiline prompt
func (*Prompts) ReadFromStdin() string {
	fmt.Print("% ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		logger.Errorf("read failed. %v", scanner.Err())
	}
	input := scanner.Text()

	return input
}

func (p *Prompts) ShouldStop(msg string) bool {
	if msg == "" || msg == "exit" {
		return true
	}

	if p.ctx.Err() != nil {
		return true
	}

	return false
}

// Converge the chat completion request from pipe and args
func (*Prompts) Converge(prompt, piped string) string {
	if len(piped) > 0 && len(prompt) > 0 {
		return prompt + ".\n" + piped
	} else if len(piped) > 0 {
		return piped
	}
	return prompt
}
