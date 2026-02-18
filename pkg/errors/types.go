package errors

import "net/http"

type RequiredError string

func (e RequiredError) Error() string {
	return string(e)
}

func (e RequiredError) StatusCode() int { return http.StatusBadRequest }

type InternalServerError string

func (e InternalServerError) Error() string {
	return string(e)
}

func (e InternalServerError) StatusCode() int { return http.StatusInternalServerError }

type ExistError string

func (e ExistError) Error() string {
	return string(e)
}

func (e ExistError) StatusCode() int { return http.StatusConflict }

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

func (e NotFoundError) StatusCode() int { return http.StatusNotFound }

type InvalidIdentifierError string

func (e InvalidIdentifierError) Error() string {
	return string(e)
}

func (e InvalidIdentifierError) StatusCode() int { return http.StatusBadRequest }

type InvalidEmailError string

func (e InvalidEmailError) Error() string {
	return string(e)
}

func (e InvalidEmailError) StatusCode() int { return http.StatusBadRequest }

type InvalidPhoneError string

func (e InvalidPhoneError) Error() string {
	return string(e)
}

func (e InvalidPhoneError) StatusCode() int { return http.StatusBadRequest }

type BadRequestError string

func (e BadRequestError) Error() string {
	return string(e)
}

func (e BadRequestError) StatusCode() int { return http.StatusBadRequest }

type InvalidLicensePlateError string

func (e InvalidLicensePlateError) Error() string {
	return string(e)
}

func (e InvalidLicensePlateError) StatusCode() int { return http.StatusBadRequest }

type InvalidIDError string

func (e InvalidIDError) Error() string {
	return string(e)
}

func (e InvalidIDError) StatusCode() int { return http.StatusBadRequest }

type InvalidPriceError string

func (e InvalidPriceError) Error() string {
	return string(e)
}

func (e InvalidPriceError) StatusCode() int { return http.StatusBadRequest }

type ConflictError string

func (e ConflictError) Error() string {
	return string(e)
}

func (e ConflictError) StatusCode() int { return http.StatusConflict }
