package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"subscriptions-api/internal/apperrors"
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

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body types.SubscriptionRequest true "Subscription data"
// @Success 201 {object} types.SubscriptionResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions [post]
func (sr *SubscriptionsRoutes) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var subReq types.SubscriptionRequest
	err = json.Unmarshal(body, &subReq)
	if err != nil || subReq.Price < 0 {
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

// GetSubscriptions godoc
// @Summary Get subscriptions list
// @Description Get paginated list of subscriptions
// @Tags subscriptions
// @Produce json
// @Param page query int true "Page number"
// @Param count query int true "Items per page"
// @Success 200 {array} types.SubscriptionResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions [get]
func (sr *SubscriptionsRoutes) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, err := strconv.Atoi(q.Get("page"))

	if err != nil || page <= 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(q.Get("count"))

	if err != nil || count <= 0 {
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

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} types.SubscriptionResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [get]
func (sr *SubscriptionsRoutes) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil || id <= 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sub, err := sr.uc.GetSubscription(id)

	if err != nil {
		if errors.Is(err, apperrors.SubscriptionNotFound) {
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

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body types.SubscriptionRequest true "Updated subscription data"
// @Success 200 {object} types.SubscriptionResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [put]
func (sr *SubscriptionsRoutes) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
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

	if err != nil || subReq.Price < 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sub, err := sr.uc.UpdateSubscription(id, subReq)

	if err != nil {
		if errors.Is(err, apperrors.SubscriptionNotFound) {
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

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} types.SubscriptionResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [delete]
func (sr *SubscriptionsRoutes) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sub, err := sr.uc.DeleteSubscriptions(id)

	if err != nil {
		sr.logger.Error("Repo Update sub", slog.Int("id", id), slog.Any("obj", sub), slog.Any("err", err))
		if errors.Is(err, apperrors.SubscriptionNotFound) {
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

// GetTotalStats godoc
// @Summary Get total subscription stats
// @Description Get total subscription price with optional filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Param service_name query string false "Service name"
// @Param start_date query string false "Start date (MM-YYYY)"
// @Param end_date query string false "End date (MM-YYYY)"
// @Success 200 {object} types.TotalStatsResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/total [get]
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
