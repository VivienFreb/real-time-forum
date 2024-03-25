package real

import (
	"strconv"
	"strings"
)

type DataError struct {
	Code    []string
	Message string
}

func NewDataError(code int, message string) *DataError {
	str := strconv.Itoa(code)
	return &DataError{strings.Split(str, ""), message}
}
