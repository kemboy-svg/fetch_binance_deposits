package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var output io.Writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

var Logger = zerolog.New(output).
	Level(zerolog.TraceLevel).
	With().
	Timestamp().
	Logger()

func ErrLog(err error, message string) error {
	output = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	Logger.Output(output)
	Logger.Error().
		Err(err).
		Msg(message)
	return errors.New(message)

}
func LogDebg(message string) string {
	Logger.Debug().Msg(message)
	return message
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", " "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
func LogInterfce(fields interface{}) interface{} {
	Logger.Debug().Fields(fields).Interface("fields", fmt.Sprintf("%s", fields))
	return fields
}
