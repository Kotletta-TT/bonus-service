package agent

import (
	"encoding/json"
	"time"

	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/go-resty/resty/v2"
)

type AgentRepository interface {
	UpdateOrderAndBalance(*models.DBOrders) error
	GetUnprocessedOrders() ([]*models.DBOrders, error)
}

type OrderAgent struct {
	client *resty.Client
	config *config.Config
	repo   AgentRepository
	jobsCh chan *models.DBOrders
	errCh  chan error
}

func NewAgent(cnf *config.Config, repo AgentRepository, jobs chan *models.DBOrders) OrderAgent {
	client := resty.New()
	client.SetRetryCount(5)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)
	return OrderAgent{client: client, repo: repo, jobsCh: jobs, config: cnf, errCh: make(chan error)}
}

func (a *OrderAgent) worker(dbOrder *models.DBOrders) {
	resp, err := a.client.R().Get(a.config.AccuralURL + "/api/orders/" + dbOrder.Number)
	if err != nil {
		a.errCh <- err
		return
	}
	if resp.StatusCode() != 200 {
		return
	}
	order := &models.AccrualOrders{}
	body := resp.Body()
	err = json.Unmarshal(body, order)
	logger.Info(string(body))
	if err != nil {
		a.errCh <- err
		return
	}
	if order.Accrual != nil {
		dbOrder.Accrual = *order.Accrual
	}
	dbOrder.Number = order.Number
	dbOrder.Status = order.Status
	if err := a.repo.UpdateOrderAndBalance(dbOrder); err != nil {
		a.errCh <- err
	}
}

func (a *OrderAgent) cronTakeUpdateOrder() {
	tic := time.NewTicker(1 * time.Second)
	for range tic.C {
		logger.Info("Run cron job")
		unprocessedOrders, err := a.repo.GetUnprocessedOrders()
		if err != nil {
			panic(err)
		}
		for _, order := range unprocessedOrders {
			a.jobsCh <- order
		}
	}
}

func (a *OrderAgent) run() {
	go a.cronTakeUpdateOrder()
	for {
		select {
		case order := <-a.jobsCh:
			a.worker(order)
		case err := <-a.errCh:
			logger.Error(err.Error())
		}
	}
}

func (a *OrderAgent) Run() {
	go a.run()
}
