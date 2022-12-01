package utils

import (
	"fmt"
	"math"
	"time"
)

type NotFoundError struct {
	subject string
	field   string
	filter  interface{}
}

func (nfe NotFoundError) Error() string {
	return fmt.Sprintf("%s with %s = %v  not found", nfe.subject, nfe.field, nfe.filter)
}

func NewNotFoundError(subject, field string, filter interface{}) NotFoundError {
	return NotFoundError{subject: subject, field: field, filter: filter}
}

type WrongOtpError struct {
	otp string
}

func (wor WrongOtpError) Error() string {
	return fmt.Sprintf("this token with value %s was wrong", wor.otp)
}

func NewWrongOtpError(otp string) WrongOtpError {
	return WrongOtpError{otp: otp}
}

type ExpireError struct {
	expiredAt time.Time
	subject   string
}

func (er ExpireError) Error() string {
	return fmt.Sprintf("this %s expired %f hour ago", er.subject, math.Abs(er.expiredAt.Sub(time.Now()).Hours()))
}

func NewExpireError(exp time.Time, subject string) ExpireError {
	return ExpireError{expiredAt: exp, subject: subject}
}

type MysqlInternalServerError struct {
	err error
}

func (mier MysqlInternalServerError) Error() string {
	return mier.err.Error()
}

func NewMysqlInternalServerError(err error) MysqlInternalServerError {
	return MysqlInternalServerError{err: err}
}
