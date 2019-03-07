package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Error struct {
	Message string                 `json:"message"`
	Err     string                 `json:"error"`
	Config  map[string]interface{} `json:"config"`
}

func New(msg string, e error) *Error {

	err := &Error{
		Message: msg,

		Config: viper.AllSettings(),
	}
	if e == nil {
		err.Err = ""
	} else {
		err.Err = e.Error()
	}
	return err
}

func NewErr(msg string) error {
	return errors.New(msg)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (e *Error) String() string {
	output, _ := json.MarshalIndent(e, "", "  ")
	return fmt.Sprintf("%s", output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (e *Error) Error() string {
	output, _ := json.MarshalIndent(e, "", "  ")
	return fmt.Sprintf("%s", output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (e *Error) FailIfErr() {
	if e.Err != "" {
		logrus.Fatal(e.String())
	}
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func (e *Error) WarnIfErr() {
	if e.Err != "" {
		logrus.Warn(e.String())
	}
}
