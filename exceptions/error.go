package exceptions

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Code    string
	Message string
	Err     error
}

// func NewError(code string, err error) error {
// 	return &Error{Code: code, Message: constants.ResponseMessage(code, ""), Err: err}
// }

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s err:%v", e.Code, e.Message, e.Err)
}

type PkgErrorEntry struct {
	*logrus.Entry
	Depth int
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (e *PkgErrorEntry) WithError(err error) *logrus.Entry {
	out := e.Entry

	common := func(pError stackTracer) {
		st := pError.StackTrace()
		depth := 3
		if e.Depth != 0 {
			depth = e.Depth
		}
		valued := fmt.Sprintf("%+v", st[0:depth])
		valued = strings.Replace(valued, "\t", "", -1)
		stack := strings.Split(valued, "\n")
		out = out.WithField("stack", stack[2:])
	}

	if err2, ok := err.(stackTracer); ok {
		common(err2)
	}

	if err2, ok := errors.Cause(err).(stackTracer); ok {
		common(err2)
	}

	return out.WithError(err)
}

// func ErrorResponse(ctx *fiber.Ctx, errorCode string, err error, parameter string, snippetCodeError string) error {
// 	response := helpers.GenerateResponse(ctx, errorCode, parameter)

// 	ctx.SendStatus(200)
// 	ctx.JSON(response)

// 	if err != nil {
// 		ctx.Locals("errorMessage", err.Error())
// 	}

// 	if len(snippetCodeError) > 0 {
// 		ctx.Locals("snippetCodeError", snippetCodeError)
// 	}

// 	return err
// }
