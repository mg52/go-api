package helper

import (
	"encoding/json"
	"github.com/mg52/go-api/domain"
	"net/http"
)

func Resp(w http.ResponseWriter, r *http.Request, httpStatusCode int, msg string) {
	var resp domain.Response
	resp.Msg = msg
	jsonBytes, _ := json.Marshal(resp)
	w.WriteHeader(httpStatusCode)
	w.Write(jsonBytes)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	var resp domain.Response
	resp.Msg = "not found"
	jsonBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusNotFound)
	w.Write(jsonBytes)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	var resp domain.Response
	resp.Msg = err.Error()
	jsonBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonBytes)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	var resp domain.Response

	w.WriteHeader(http.StatusUnauthorized)
	if err != nil {
		resp.Msg = err.Error()
		jsonBytes, _ := json.Marshal(resp)
		w.Write(jsonBytes)
	} else {
		resp.Msg = "unauthorized"
		jsonBytes, _ := json.Marshal(resp)
		w.Write(jsonBytes)
	}

}
