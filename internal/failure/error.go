package failure

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FailureCode uint32

const (
	InternalError   FailureCode = 500
	ValidationError FailureCode = 400
)

type Failure struct {
	code    FailureCode
	message string
}

func New(code FailureCode, message string) *Failure {
	return &Failure{
		code:    code,
		message: message,
	}
}
func Err(code FailureCode, message string) error {
	return &Failure{
		code:    code,
		message: message,
	}
}

func (f *Failure) Context() FailureCode {
	return f.code
}
func (f *Failure) Message() string {
	return f.message
}

func (f *Failure) Error() string {
	return fmt.Sprintf(
		"context=%s;message=%s",
		f.code.ToString(),
		f.message,
	)
}

func (f FailureCode) ToString() string {
	switch f {
	case InternalError:
		return "InternalError"
	case ValidationError:
		return "ValidationError"
	}
	return "UnknownError"
}

func (f FailureCode) ToGrpcCode() codes.Code {
	switch f {
	case InternalError:
		return codes.Internal
	case ValidationError:
		return codes.InvalidArgument
	}
	return codes.Unknown
}

func (f *Failure) ToGrpcError() error {
	if f == nil {
		return nil
	}
	return status.Error(
		f.code.ToGrpcCode(),
		f.Message(),
	)
}
