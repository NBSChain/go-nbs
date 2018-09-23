package utils

import (
	"fmt"
)

const (ErrCommon  	= -1
	ErrSystem 	= -2)



type NbsError struct {
	errCode	int
	errMsg	string
}

func (e *NbsError) Error() string {
	return fmt.Sprintf("%d - %s", e.errCode, e.errMsg)
}

func New(text string) error {
	return &NbsError{errCode:ErrCommon, errMsg:text}
}

func NewWithCode(code int, text string) error {
	return &NbsError{errCode:code, errMsg:text}
}