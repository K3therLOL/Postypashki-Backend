package errorfmt

import (
	"encoding/json"
)

type errorJson struct {
	Err string `json:"error"`
}

func newErrorJson(err string) errorJson {
	return errorJson{Err: err}
}

func Jsonize(err error) string {
	errStruct := newErrorJson(err.Error())
	errJson, _ := json.Marshal(errStruct)
	return string(errJson)
}
