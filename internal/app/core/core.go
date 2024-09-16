package core

import (
	"context"

	"go.uber.org/zap"

	"avito2024/internal/adapter/repo"
	"avito2024/internal/app/core/service"
	"avito2024/internal/config"
	v1 "avito2024/internal/controller/api/v1"
)

func Run(cfg *config.Config) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	postgresRepo, err := repo.NewPostgresRepo(ctx, cfg.ConnectionString, logger)
	if err != nil {
		panic(err)
	}

	tenderRepo, err := postgresRepo.NewTenderRepo(ctx)
	if err != nil {
		panic(err)
	}

	bidRepo, err := postgresRepo.NewBidRepo(ctx)
	if err != nil {
		panic(err)
	}

	userRepo, err := postgresRepo.NewUserRepo(ctx, cfg.IsTest)
	if err != nil {
		panic(err)
	}

	orgRepo, err := postgresRepo.NewOrganizationRepo(ctx, cfg.IsTest)
	if err != nil {
		panic(err)
	}

	tenderService := service.NewTenderService(tenderRepo, userRepo, orgRepo)
	bidService := service.NewBidService(bidRepo, userRepo, orgRepo, tenderRepo)
	v1.NewAPI(tenderService, bidService, logger).Run(cfg.Host)

}
