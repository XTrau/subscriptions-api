package types

import (
	"strings"
	"time"
)

type MonthYear time.Time

func (m *MonthYear) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)

	t, err := time.Parse("01-2006", str)
	if err != nil {
		return err
	}

	*m = MonthYear(t)
	return nil
}

type SubscriptionRequest struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   MonthYear `json:"start_date"`
}

type SubscriptionResponse struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_date"`
}

type TotalStatsResponse struct {
	Total int `json:"total"`
}
