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
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

func (h *CustomRouter) getItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := r.Body.Close(); err != nil {
				h.log.Error("Can't close body", h.log.Attr("error", err))
			}
		}()

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

		from := time.Unix(req.From, 0).Format(time.DateOnly)
		to := time.Unix(req.To, 0).Format(time.DateOnly)

		sups, err := h.storage.GetReserved(ctx, from, to)
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

		if err = getItemsRespJSONOk(w, http.StatusOK, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}

type requestSup struct {
	ID   int64  `json:"id"`
	Name string `json:"model_name"`
	From int64  `json:"from"`
	To   int64  `json:"to"`
}

type makeReservationRequest struct {
	ID           int64        `json:"id"`
	ClientNumber string       `json:"client_phone"`
	ClientName   string       `json:"client_name"`
	SupsInfo     []requestSup `json:"sups_info"`
	FullPrice    int64        `json:"price"`
}

func (h *CustomRouter) makeReservation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := r.Body.Close(); err != nil {
				h.log.Error("Can't close body", h.log.Attr("error", err))
			}
		}()

		req := makeReservationRequest{}

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

		approve := models.Approve{
			SupsInfo: make([]models.ApproveSup, len(req.SupsInfo)),
		}

		approve.ID = req.ID
		approve.ClientNumber = req.ClientNumber
		approve.ClientName = req.ClientName
		approve.FullPrice = req.FullPrice
		for i, info := range req.SupsInfo {
			approve.SupsInfo[i].ID = info.ID
			approve.SupsInfo[i].Name = info.Name
			approve.SupsInfo[i].From = time.Unix(info.From, 0)
			approve.SupsInfo[i].To = time.Unix(info.To, 0)
		}

		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()

		_, err := h.storage.CreateApprove(ctx, approve)
		if err != nil {
			h.log.Error("Can't create new command", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		if err := h.notifier.PushNotice(); err != nil {
			h.log.Error("Can't push notification", h.log.Attr("error", err))
		}

		respBody := makeReservationRespBodyOK{
			Created: true,
		}

		if err = makeReservationRespJSONOk(w, http.StatusOK, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}
