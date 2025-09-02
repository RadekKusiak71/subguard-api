package utils

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ReadParamFromPathAsInt(r *http.Request, key string) (int, error) {
	valueAsString := mux.Vars(r)[key]

	value, err := strconv.Atoi(valueAsString)

	if err != nil {
		return 0, err
	}

	return value, nil
}
