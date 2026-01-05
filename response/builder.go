package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Msg    string
	Reason string
}

type responseBuilder struct {
	status int
	body   any
}

func NewResponseBuilder(status int) responseBuilder {
	return responseBuilder{
		status: status,
	}
}

func (r responseBuilder) SetBody(body any) responseBuilder {
	return responseBuilder{
		status: r.status,
		body:   body,
	}
}

func (r responseBuilder) SetError(msg string, reason error) responseBuilder {
	return responseBuilder{
		status: r.status,
		body:   ErrorResponse{Msg: msg, Reason: reason.Error()},
	}
}

func (r responseBuilder) SendAsJSON(w http.ResponseWriter) error {
	w.WriteHeader(r.status)
	bodyBytes, err := json.Marshal(r.body)
	if err != nil {
		return err
	}

	_, err = w.Write(bodyBytes)
	return err
}
