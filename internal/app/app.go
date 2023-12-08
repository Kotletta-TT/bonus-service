package app

import (
	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/agent"
	"github.com/Kotletta-TT/bonus-service/internal/controller"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/Kotletta-TT/bonus-service/internal/repository/postgres"
)

func Run(cnf *config.Config) {
	repo := postgres.WrapErrDBNew(cnf)
	jobsCh := make(chan *models.DBOrders, 10)
	defer repo.Close()
	router := controller.Router(cnf, repo)
	orderAgent := agent.NewAgent(cnf, repo, jobsCh)
	orderAgent.Run()
	logger.Error(router.Run(cnf.ServAddr).Error())
}
