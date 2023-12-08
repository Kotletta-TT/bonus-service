package main

import (
	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/app"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
)

func main() {
	config := config.New()
	logger.Init(config)
	app.Run(config)
}
