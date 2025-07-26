package customerror

import (
	"fmt"
	"strings"

	"github.com/keshu12345/overlap-avalara/constants"
	errorsPkg "github.com/pkg/errors"
)

type CustomError struct {
	code      constants.Code
	message   string
	data      any
	errMap    map[string]string
	err       error
	exists    bool
	retryable bool
	notify    bool
	logParams map[string]any
}

func NewCustomError(errorCode constants.Code, error string, options ...func(*CustomError)) CustomError {
	c := CustomError{
		code:      errorCode,
		message:   error,
		data:      nil,
		exists:    true,
		retryable: false,
		notify:    true,
	}

	for _, option := range options {
		option(&c)
	}

	e := fmt.Errorf("Code: %s | %s", c.code, c.message)
	c.err = errorsPkg.WithStack(e)
	c.logParams = make(map[string]any, 0)
	c.errMap = make(map[string]string, 0)
	return c
}

func NewCustomErrorWithPayload(errorCode constants.Code, error string, data any, options ...func(*CustomError)) CustomError {
	c := NewCustomError(errorCode, error, options...)
	c.data = data
	return c
}

func WithRetryable(retryable bool) func(*CustomError) {
	return func(c *CustomError) {
		c.retryable = retryable
	}
}

func WithShouldNotify(shouldNotify bool) func(*CustomError) {
	return func(c *CustomError) {
		c.notify = shouldNotify
	}
}

func RequestInvalidError(message string, options ...func(*CustomError)) CustomError {
	c := CustomError{
		code:      constants.RequestInvalid,
		message:   message,
		data:      nil,
		exists:    true,
		retryable: false,
		notify:    true,
	}
	e := fmt.Errorf("Code: %s | %s", c.code, c.message)
	c.err = errorsPkg.WithStack(e)
	c.logParams = make(map[string]any, 0)
	c.errMap = make(map[string]string, 0)

	for _, option := range options {
		option(&c)
	}
	return c
}

func (c CustomError) Exists() bool {
	return c.exists
}

func (c CustomError) Log() {
	fmt.Println(c.ToString())
}

func (c CustomError) LoggingParams() map[string]any {
	return c.logParams
}

func (c CustomError) ErrorCode() constants.Code {
	return c.code
}

func (c CustomError) ToError() error {
	return c.err
}

func (c CustomError) Error() string {
	return c.err.Error()
}

func (c CustomError) ErrorMessage() string {
	return c.message
}

func (c CustomError) ShouldNotify() bool {
	return c.notify
}

func (c CustomError) Retryable() bool {
	return c.retryable
}

func (c CustomError) ToString() string {
	logMsg := fmt.Sprintf("Code: %s, Msg: %s", c.code, c.message)

	paramStrings := make([]string, 0)
	for key, val := range c.logParams {
		paramStrings = append(paramStrings, fmt.Sprintf("%s: {%+v}", strings.ToUpper(key), val))
	}
	return fmt.Sprintf("%s, Params: [%+v]", logMsg, strings.Join(paramStrings, " | "))
}

func (c CustomError) WithParam(key string, val any) CustomError {
	if c.logParams == nil {
		c.logParams = make(map[string]any, 0)
	}
	c.logParams[key] = val
	return c
}

func (c CustomError) ErrorString() string {
	return c.message
}

func (c CustomError) UserMessage() string {
	return c.message
}

func (c CustomError) ErrorData() any {
	return c.data
}

func (c CustomError) ErrorMap() map[string]string {
	return c.errMap
}

func WithErrors(errors map[string]string) func(*CustomError) {
	return func(c *CustomError) {
		c.errMap = errors
	}
}

func WithData(data any) func(*CustomError) {
	return func(c *CustomError) {
		c.data = data
	}
}
