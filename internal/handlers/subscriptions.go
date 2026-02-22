package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"subscriptions-api/internal/errors"
	"subscriptions-api/internal/responses"
	"subscriptions-api/internal/types"
	"subscriptions-api/internal/usecases"

	"github.com/go-chi/chi/v5"
)

type SubscriptionsRoutes struct {
	uc     usecases.SubscriptionUseCases
	logger *slog.Logger
}

func NewSubscriptionsRoutes(uc usecases.SubscriptionUseCases, logger *slog.Logger) SubscriptionsRoutes {
	return SubscriptionsRoutes{uc, logger}
}

func (sr *SubscriptionsRoutes) RegisterRoutes(r chi.Router) {
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", sr.CreateSubscription)
		r.Get("/", sr.GetSubscriptions)
		r.Get("/{id}", sr.GetSubscription)
		r.Put("/{id}", sr.UpdateSubscription)
		r.Delete("/{id}", sr.DeleteSubscription)
		r.Get("/total", sr.GetTotalStats)
	})
}

func (sr *SubscriptionsRoutes) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var subReq types.SubscriptionRequest
	err = json.Unmarshal(body, &subReq)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sub, err := sr.uc.SaveSubscription(subReq)
	if err != nil {
		sr.logger.Error("Repo failed on create", slog.Any("obj", subReq), slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = responses.SetJsonBody(w, sub)

	if err != nil {
		sr.logger.Error("Json set body", slog.Any("obj", sub), slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (sr *SubscriptionsRoutes) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, err := strconv.Atoi(q.Get("page"))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(q.Get("count"))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	subs, err := sr.uc.GetSubscriptions(page, count)

	if err != nil {
		sr.logger.Error("Repo Get subs", slog.Int("page", page), slog.Int("count", count), slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = responses.SetJsonBody(w, subs)

	if err != nil {
		sr.logger.Error("Json set body", slog.Any("obj", subs), slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (sr *SubscriptionsRoutes) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sub, err := sr.uc.GetSubscription(id)

	if err != nil {
		if err == errors.SubscriptionNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			sr.logger.Error("Repo Get sub", slog.Int("id", id), slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = responses.SetJsonBody(w, sub)

	if err != nil {
		sr.logger.Error("Json set body", slog.Any("obj", sub), slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (sr *SubscriptionsRoutes) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var subReq types.SubscriptionRequest
	err = json.Unmarshal(body, &subReq)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sub, err := sr.uc.UpdateSubscription(id, subReq)

	if err != nil {
		if err == errors.SubscriptionNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			sr.logger.Error("Repo Update sub", slog.Int("id", id), slog.Any("obj", sub), slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = responses.SetJsonBody(w, sub)

	if err != nil {
		sr.logger.Error("Json set body", slog.Any("obj", sub), slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (sr *SubscriptionsRoutes) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sub, err := sr.uc.DeleteSubscriptions(id)

	if err != nil {
		sr.logger.Error("Repo Update sub", slog.Int("id", id), slog.Any("obj", sub), slog.Any("err", err))
		if err == errors.SubscriptionNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = responses.SetJsonBody(w, sub)

	if err != nil {
		sr.logger.Error("Json set body", slog.Any("obj", sub), slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (sr *SubscriptionsRoutes) GetTotalStats(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	userID := q.Get("user_id")
	serviceName := q.Get("service_name")
	startDateQuery := q.Get("start_date")
	endDateQuery := q.Get("end_date")

	var startDatePtr, endDatePtr *time.Time

	if len(startDateQuery) != 0 {
		startDate, err := time.Parse("01-2006", startDateQuery)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		startDatePtr = &startDate
	}

	if len(endDateQuery) != 0 {
		endDate, err := time.Parse("01-2006", endDateQuery)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		endDatePtr = &endDate
	}

	total, err := sr.uc.GetTotalStats(
		serviceName,
		userID,
		startDatePtr,
		endDatePtr,
	)

	err = responses.SetJsonBody(w, total)

	if err != nil {
		sr.logger.Error("Json set body", slog.Any("obj", total), slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
