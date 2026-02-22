package usecases

import (
	"subscriptions-api/internal/repositories"
	"subscriptions-api/internal/types"
	"time"
)

type SubscriptionUseCases struct {
	repo repositories.SubscriptionsRepository
}

func NewSubscriptionUseCases(repo repositories.SubscriptionsRepository) SubscriptionUseCases {
	return SubscriptionUseCases{repo}
}

func (uc *SubscriptionUseCases) SaveSubscription(sub types.SubscriptionRequest) (types.SubscriptionResponse, error) {
	return uc.repo.SaveSubscription(sub)
}

func (uc *SubscriptionUseCases) GetSubscription(id int) (types.SubscriptionResponse, error) {
	return uc.repo.GetSubscription(id)
}

func (uc *SubscriptionUseCases) GetSubscriptions(page int, count int) ([]types.SubscriptionResponse, error) {
	return uc.repo.GetSubscriptions((page-1)*count, count)
}

func (uc *SubscriptionUseCases) DeleteSubscriptions(id int) (types.SubscriptionResponse, error) {
	return uc.repo.DeleteSubscription(id)
}

func (uc *SubscriptionUseCases) UpdateSubscription(id int, subscription types.SubscriptionRequest) (types.SubscriptionResponse, error) {
	return uc.repo.UpdateSubscription(id, subscription)
}

func (uc *SubscriptionUseCases) GetTotalStats(serviceName, userID string, startDate, endDate *time.Time) (types.TotalStatsResponse, error) {
	subs, err := uc.repo.GetSubscriptionsByFilter(serviceName, userID, startDate, endDate)

	if err != nil {
		return types.TotalStatsResponse{}, err
	}

	var total int
	for _, sub := range subs {
		total += sub.Price
	}

	return types.TotalStatsResponse{Total: total}, nil
}
