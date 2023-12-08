package errors

import "net/http"

type UsersError struct {
	Code int
	Err  string
}

func (e *UsersError) Error() string {
	return e.Err
}

func InternalServerErr() error {
	return &UsersError{Code: http.StatusInternalServerError, Err: "internal server error"}
}

func IncorrectOrderNumber() error {
	return &UsersError{Code: http.StatusUnprocessableEntity, Err: "incorrect order number: does not exist"}
}

func UploadOrderEarlier() error {
	return &UsersError{Code: http.StatusOK, Err: "you uploaded this order number earlier"}
}

func UploadAnotherUser() error {
	return &UsersError{Code: http.StatusConflict, Err: "another user has already uploaded this order number"}
}

func InvalidTokenErr() error {
	return &UsersError{Code: http.StatusUnauthorized, Err: "token is not valid"}
}

func ExpiredTokenErr() error {
	return &UsersError{Code: http.StatusUnauthorized, Err: "token is expired"}
}

func AuthErr() error {
	return &UsersError{Code: http.StatusUnauthorized, Err: "unknown login/password"}
}

func EnoughtBalance() error {
	return &UsersError{Code: http.StatusPaymentRequired, Err: "your balance is less than the request amount"}
}
