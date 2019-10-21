package cliopts

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	cli "gopkg.in/urfave/cli.v2"
)

type EnsurePair struct {
	optName    string
	ensureFunc func() error
}

type EnsurePairs []*EnsurePair

func NewEnsure() EnsurePairs {
	return make([]*EnsurePair, 0, 8)
}

func (ensurePairs EnsurePairs) With(name string, f func() error) EnsurePairs {
	ensurePairs = append(ensurePairs, &EnsurePair{
		optName:    name,
		ensureFunc: f,
	})
	return ensurePairs
}

type ReadLineFunc func() (string, error)

func Interact(prompt string, in io.Reader, out io.Writer) ReadLineFunc {
	return func() (string, error) {
		if out != nil {
			fmt.Fprint(out, prompt)
		}
		if in == nil {
			return "", errors.New("in io.Reader is nil")
		}
		reader := bufio.NewReader(in)
		line, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}
		return string(line), nil
	}
}

func StdInteract(prompt string) ReadLineFunc {
	return Interact(prompt, os.Stdin, os.Stdout)
}

func (pre ReadLineFunc) After(f func(string) error) func() error {
	return func() error {
		read, err := pre()
		if err != nil {
			return err
		}
		err = f(read)
		if err != nil {
			return err
		}
		return nil
	}
}

func (ensurePairs EnsurePairs) Do(ctx *cli.Context) error {
	for _, pair := range ensurePairs {
		if ctx.IsSet(pair.optName) {
			continue
		}
		err := pair.ensureFunc()
		if err != nil {
			return err
		}
	}
	return nil
}
