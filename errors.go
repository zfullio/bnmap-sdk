package bnmap

import "errors"

var ErrDecodeToken = errors.New("failed to read token")

type DecodeError struct {
	msg string
}

func (e *DecodeError) Error() string {
	return e.msg
}
