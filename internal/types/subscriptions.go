package types

import (
	"strings"
	"time"

	"github.com/google/uuid"
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
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       int       `json:"price" example:"999"`
	UserID      uuid.UUID `json:"user_id" format:"uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   MonthYear `json:"start_date" swaggertype:"string" example:"01-2025"`
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
