package env

import (
	"fmt"

	"github.com/mattn/go-shellwords"
)

func Parse(env []string) ([]string, error) {
	parsedEnv := make([]string, 0, len(env))

	// Unquote any quoted .env vars.
	for _, line := range env {
		parsed, err := shellwords.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("parsing .env line: %w", err)
		}

		if len(parsed) == 0 {
			continue
		}

		if len(parsed) != 1 {
			return nil, fmt.Errorf("invalid .env line: %s", line)
		}

		parsedEnv = append(parsedEnv, parsed[0])
	}

	return parsedEnv, nil
}
