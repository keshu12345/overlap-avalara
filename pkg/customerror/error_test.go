package customerror

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/keshu12345/overlap-avalara/constants"
)

func TestNewCustomError_Defaults(t *testing.T) {
	code := constants.Code("TEST_CODE")
	msg := "something went wrong"

	c := NewCustomError(code, msg)

	if !c.Exists() {
		t.Error("Expected Exists() to be true")
	}
	if c.Retryable() {
		t.Error("Expected Retryable() to be false by default")
	}
	if !c.ShouldNotify() {
		t.Error("Expected ShouldNotify() to be true by default")
	}
	if c.ErrorCode() != code {
		t.Errorf("Expected ErrorCode() %v; got %v", code, c.ErrorCode())
	}
	if c.ErrorMessage() != msg {
		t.Errorf("Expected ErrorMessage() %q; got %q", msg, c.ErrorMessage())
	}

	err := c.ToError()
	if err == nil {
		t.Fatal("Expected ToError() non‑nil error")
	}
	expectedErr := fmt.Sprintf("Code: %s | %s", code, msg)
	if err.Error() != expectedErr {
		t.Errorf("Expected err.Error() %q; got %q", expectedErr, err.Error())
	}

	if c.ErrorString() != msg {
		t.Errorf("Expected ErrorString() %q; got %q", msg, c.ErrorString())
	}
	if c.UserMessage() != msg {
		t.Errorf("Expected UserMessage() %q; got %q", msg, c.UserMessage())
	}

	if data := c.ErrorData(); data != nil {
		t.Errorf("Expected ErrorData() nil; got %v", data)
	}
	if m := c.ErrorMap(); len(m) != 0 {
		t.Errorf("Expected empty ErrorMap(); got %v", m)
	}
	if p := c.LoggingParams(); len(p) != 0 {
		t.Errorf("Expected empty LoggingParams(); got %v", p)
	}
}

func TestNewCustomError_WithOptions(t *testing.T) {
	code := constants.Code("OPT_CODE")
	msg := "opt msg"

	c := NewCustomError(
		code, msg,
		WithRetryable(true),
		WithShouldNotify(false),
	)

	if !c.Retryable() {
		t.Error("Expected Retryable() to be true")
	}
	if c.ShouldNotify() {
		t.Error("Expected ShouldNotify() to be false")
	}
}

func TestNewCustomErrorWithPayload(t *testing.T) {
	code := constants.Code("PAYLOAD_CODE")
	msg := "payload msg"
	payload := map[string]int{"x": 1}

	c := NewCustomErrorWithPayload(code, msg, payload)

	if data := c.ErrorData(); !reflect.DeepEqual(data, payload) {
		t.Errorf("Expected ErrorData() %v; got %v", payload, data)
	}
}

func TestRequestInvalidError_Default(t *testing.T) {
	msg := "invalid request"
	c := RequestInvalidError(msg)

	if c.ErrorCode() != constants.RequestInvalid {
		t.Errorf("Expected RequestInvalid code; got %v", c.ErrorCode())
	}
	if !c.Exists() {
		t.Error("Expected Exists() to be true")
	}
	if c.Retryable() {
		t.Error("Expected Retryable() to be false")
	}
	if !c.ShouldNotify() {
		t.Error("Expected ShouldNotify() to be true")
	}
	if c.ErrorMessage() != msg {
		t.Errorf("Expected ErrorMessage() %q; got %q", msg, c.ErrorMessage())
	}
}

func TestWithParam_ToString(t *testing.T) {
	code := constants.Code("P_CODE")
	msg := "param message"

	c := NewCustomError(code, msg).WithParam("foo", 42)

	params := c.LoggingParams()
	v, ok := params["foo"]
	if !ok || v != 42 {
		t.Errorf("Expected LoggingParams()[\"foo\"]=42; got %v", params)
	}

	got := c.ToString()
	want := fmt.Sprintf("Code: %s, Msg: %s, Params: [FOO: {42}]", code, msg)
	if got != want {
		t.Errorf("Expected ToString() %q; got %q", want, got)
	}
}

func TestLog_DoesNotPanic(t *testing.T) {
	c := NewCustomError(constants.Code("L_CODE"), "log test")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c.Log()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "Code: L_CODE, Msg: log test") {
		t.Errorf("Expected log output to contain message; got %q", out)
	}
}
