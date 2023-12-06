package models

import (
	"time"

	"github.com/google/uuid"
)

type DBOrders struct {
	UserID     uuid.UUID
	Number     string
	Status     string
	UploadedAt time.Time
	Accrual    float64
}

type ViewOrders struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
	Accrual    *float64  `json:"accrual,omitempty"`
}

type AccrualOrders struct {
	Number  string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

func ConvertDBToView(src []*DBOrders) []*ViewOrders {
	dst := make([]*ViewOrders, 0, len(src))
	for _, dbOrder := range src {
		viewOrder := &ViewOrders{}
		viewOrder.Number = dbOrder.Number
		viewOrder.Status = dbOrder.Status
		viewOrder.UploadedAt = dbOrder.UploadedAt
		if dbOrder.Accrual > 0 {
			viewOrder.Accrual = &dbOrder.Accrual
		}
		dst = append(dst, viewOrder)
	}
	return dst
}
