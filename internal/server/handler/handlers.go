package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/enchik0reo/sup-back/internal/models"
)

type getItemsRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func (h *CustomRouter) getItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		req := getItemsRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if !strings.Contains(err.Error(), "EOF") {
				h.log.Debug("Can't decode body from get items request", h.log.Attr("error", err))

				err = responseJSONError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				if err != nil {
					h.log.Error("Can't make response", h.log.Attr("error", err))
				}
				return
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()

		sups, err := h.storager.GetReserved(ctx, req.From, req.To)
		if err != nil {
			h.log.Error("Can't create new command", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		respBody := getItemsRespBodyOK{
			Sups: sups,
		}

		if err = idRespJSONOk(w, http.StatusOK, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}

func (h *CustomRouter) makeReservation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		approve := models.Approve{}

		if err := json.NewDecoder(r.Body).Decode(&approve); err != nil {
			if !strings.Contains(err.Error(), "EOF") {
				h.log.Debug("Can't decode body from get items request", h.log.Attr("error", err))

				err = responseJSONError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				if err != nil {
					h.log.Error("Can't make response", h.log.Attr("error", err))
				}
				return
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()

		_, err := h.storager.CreateApprove(ctx, approve)

		if err != nil {
			h.log.Error("Can't create new command", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		respBody := makeReservationRespBodyOK{
			Created: true,
		}

		if err = makeReservationRespJSONOk(w, http.StatusOK, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}
