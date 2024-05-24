package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enchik0reo/sup-back/internal/models"
)

type getItemsRespOK struct {
	Status int                `json:"status"`
	Body   getItemsRespBodyOK `json:"body"`
}

type getItemsRespBodyOK struct {
	Sups []models.Sup `json:"sups,omitempty"`
}

func getItemsRespJSONOk(w http.ResponseWriter, status int, body getItemsRespBodyOK) error {
	resp := getItemsRespOK{
		Status: status,
		Body:   body,
	}

	w.Header().Add("Content-Type", "application/json")

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(respJSON)
	if err != nil {
		return err
	}

	return nil
}

type makeReservationRespOK struct {
	Status int                       `json:"status"`
	Body   makeReservationRespBodyOK `json:"body"`
}

type makeReservationRespBodyOK struct {
	Created bool `json:"created,omitempty"`
}

func makeReservationRespJSONOk(w http.ResponseWriter, status int, body makeReservationRespBodyOK) error {
	resp := makeReservationRespOK{
		Status: status,
		Body:   body,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(respJSON)
	if err != nil {
		return err
	}

	return nil
}

type responseErr struct {
	Status int         `json:"status"`
	Body   respBodyErr `json:"body"`
}

type respBodyErr struct {
	Error string `json:"error,omitempty"`
}

func responseJSONError(w http.ResponseWriter, status int, error string) error {
	resp := responseErr{
		Status: status,
	}

	if error != "" {
		resp.Body.Error = error
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(respJSON)
	if err != nil {
		return err
	}

	return nil
}
