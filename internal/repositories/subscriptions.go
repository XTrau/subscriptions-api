package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"subscriptions-api/internal/apperrors"
	"subscriptions-api/internal/types"
	"time"
)

type SubscriptionsRepository interface {
	SaveSubscription(sub types.SubscriptionRequest) (types.SubscriptionResponse, error)
	GetSubscription(id int) (types.SubscriptionResponse, error)
	GetSubscriptions(offset int, count int) ([]types.SubscriptionResponse, error)
	GetSubscriptionsByFilter(serviceName, userID string, startDate, endDate *time.Time) ([]types.SubscriptionResponse, error)
	UpdateSubscription(id int, sub types.SubscriptionRequest) (types.SubscriptionResponse, error)
	DeleteSubscription(id int) (types.SubscriptionResponse, error)
}

type SubscriptionsPostgresRepository struct {
	db *sql.DB
}

func NewSubscriptionsPostgresRepository(db *sql.DB) SubscriptionsPostgresRepository {
	return SubscriptionsPostgresRepository{db}
}

func (sr SubscriptionsPostgresRepository) SaveSubscription(sub types.SubscriptionRequest) (types.SubscriptionResponse, error) {
	query := `
		INSERT INTO subscriptions 
		(ServiceName, Price, UserID, StartDate) 
		VALUES ($1, $2, $3, $4) 
		RETURNING ID, ServiceName, Price, UserID, StartDate
	`

	r := sr.db.QueryRow(
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		time.Time(sub.StartDate),
	)

	var id, price int
	var serviceName, userID string
	var startDate time.Time

	if err := r.Scan(&id, &serviceName, &price, &userID, &startDate); err != nil {
		return types.SubscriptionResponse{}, err
	}

	res := types.SubscriptionResponse{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
	}

	return res, nil
}

func (sr SubscriptionsPostgresRepository) GetSubscription(id int) (types.SubscriptionResponse, error) {
	query := `
		SELECT ServiceName, Price, UserID, StartDate 
		FROM subscriptions
		WHERE id = $1
	`

	r := sr.db.QueryRow(query, id)

	var serviceName, userID string
	var price int
	var startDate time.Time

	if err := r.Scan(&serviceName, &price, &userID, &startDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.SubscriptionResponse{}, apperrors.SubscriptionNotFound
		}
		return types.SubscriptionResponse{}, err
	}

	return types.SubscriptionResponse{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
	}, nil
}

func (sr SubscriptionsPostgresRepository) GetSubscriptions(offset int, limit int) ([]types.SubscriptionResponse, error) {
	query := `
		SELECT id, ServiceName, Price, UserID, StartDate 
		FROM subscriptions
		OFFSET $1 LIMIT $2
	`

	rows, err := sr.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]types.SubscriptionResponse, 0)

	for rows.Next() {
		var id, price int
		var serviceName, userID string
		var startDate time.Time

		err := rows.Scan(&id, &serviceName, &price, &userID, &startDate)
		if err != nil {
			return nil, err
		}

		result = append(
			result,
			types.SubscriptionResponse{
				ID:          id,
				ServiceName: serviceName,
				Price:       price,
				UserID:      userID,
				StartDate:   startDate,
			},
		)
	}

	return result, nil

}

func (sr SubscriptionsPostgresRepository) UpdateSubscription(id int, sub types.SubscriptionRequest) (types.SubscriptionResponse, error) {
	query := `
		UPDATE subscriptions
		SET ServiceName=$2, Price=$3, UserID=$4, StartDate=$5
		WHERE id=$1
		RETURNING ServiceName, Price, UserID, StartDate
	`

	r := sr.db.QueryRow(query, id, sub.ServiceName, sub.Price, sub.UserID, time.Time(sub.StartDate))

	var price int
	var serviceName, userID string
	var startDate time.Time

	if err := r.Scan(&serviceName, &price, &userID, &startDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.SubscriptionResponse{}, apperrors.SubscriptionNotFound
		}
		return types.SubscriptionResponse{}, err
	}

	return types.SubscriptionResponse{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
	}, nil
}

func (sr SubscriptionsPostgresRepository) DeleteSubscription(id int) (types.SubscriptionResponse, error) {
	query := `
		DELETE FROM subscriptions
		WHERE id=$1
		RETURNING ServiceName, Price, UserID, StartDate
	`

	r := sr.db.QueryRow(query, id)

	var price int
	var serviceName, userID string
	var startDate time.Time

	if err := r.Scan(&serviceName, &price, &userID, &startDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.SubscriptionResponse{}, apperrors.SubscriptionNotFound
		}
		return types.SubscriptionResponse{}, err
	}

	return types.SubscriptionResponse{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
	}, nil
}

func (sr SubscriptionsPostgresRepository) GetSubscriptionsByFilter(serviceName, userID string, startDate, endDate *time.Time) ([]types.SubscriptionResponse, error) {
	query := `
		SELECT id, ServiceName, Price, userID, StartDate
		FROM subscriptions WHERE 1=1
	`

	args := make([]interface{}, 0, 3)
	argc := 1

	if len(serviceName) != 0 {
		query += fmt.Sprintf(" AND ServiceName=$%d", argc)
		args = append(args, serviceName)
		argc++
	}

	if len(userID) != 0 {
		query += fmt.Sprintf(" AND UserID=$%d", argc)
		args = append(args, userID)
		argc++
	}

	if startDate != nil {
		query += fmt.Sprintf(" AND StartDate >= $%d", argc)
		args = append(args, *startDate)
		argc++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND StartDate <= $%d", argc)
		args = append(args, *endDate)
		argc++
	}

	rows, err := sr.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]types.SubscriptionResponse, 0)

	for rows.Next() {
		var ID, Price int
		var ServiceName, UserID string
		var StartDate time.Time

		err := rows.Scan(&ID, &ServiceName, &Price, &UserID, &StartDate)

		if err != nil {
			return nil, err
		}

		result = append(result, types.SubscriptionResponse{
			ID:          ID,
			ServiceName: ServiceName,
			UserID:      UserID,
			Price:       Price,
			StartDate:   StartDate,
		})
	}

	return result, nil
}
