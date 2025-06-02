package utility

import (
	"strconv"
)

type StrErr string

func (e StrErr) Error() string {
	return string(e)
}

func Atoi(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	} else {
		panic(StrErr(err.Error()))
	}
}

func IsDateFormat(s string) bool {
	//TODO check date format is preserved in the POST request
	return true
}
