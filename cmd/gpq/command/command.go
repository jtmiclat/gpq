package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/planetlabs/gpq/internal/storage"
)

var CLI struct {
	Convert  ConvertCmd  `cmd:"" help:"Convert data from one format to another."`
	Validate ValidateCmd `cmd:"" help:"Validate a GeoParquet file."`
	Describe DescribeCmd `cmd:"" help:"Describe a GeoParquet file."`
	Version  VersionCmd  `cmd:"" help:"Print the version of this program."`
}

type CommandError struct {
	err error
}

func NewCommandError(format string, a ...any) *CommandError {
	return &CommandError{err: fmt.Errorf(format, a...)}
}

func (e *CommandError) Error() string {
	return e.err.Error()
}

func (e *CommandError) Unwrap() error {
	return e.err
}

func readerFromInput(input string) (storage.ReaderAtSeeker, error) {
	if input == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("trouble reading from stdin: %w", err)
		}
		return bytes.NewReader(data), nil
	}

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return storage.NewHttpReader(input)
	}

	return os.Open(input)
}
