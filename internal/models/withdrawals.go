package models

import (
	"time"

	"github.com/google/uuid"
)

type DBWithdraw struct {
	UserID      uuid.UUID
	OrderID     string
	Sum         float64
	ProcessedAt time.Time
}

type ViewWithdraw struct {
	OrderID     string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawRequest struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}

func ConvertWithdrawDBView(src []*DBWithdraw) []*ViewWithdraw {
	dst := make([]*ViewWithdraw, 0, len(src))
	for _, dbWithdraw := range src {
		viewWithdraw := &ViewWithdraw{}
		viewWithdraw.OrderID = dbWithdraw.OrderID
		viewWithdraw.Sum = dbWithdraw.Sum
		viewWithdraw.ProcessedAt = dbWithdraw.ProcessedAt
		dst = append(dst, viewWithdraw)
	}
	return dst
}
