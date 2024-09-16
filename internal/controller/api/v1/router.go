package controller

import (
	"net/http"

	"avito2024/internal/app/core/service"
	"avito2024/internal/controller/api/v1/handler/bid"
	"avito2024/internal/controller/api/v1/handler/tender"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type parentRouter struct {
	tenderService *service.TenderService
	bidService    *service.BidService
	logger        *zap.Logger
}

func (r *parentRouter) TenderService() *service.TenderService {
	return r.tenderService
}

func (r *parentRouter) BidService() *service.BidService {
	return r.bidService
}

func (r *parentRouter) Logger() *zap.Logger {
	return r.logger
}

func NewAPI(
	tenderService *service.TenderService,
	bidService *service.BidService,
	logger *zap.Logger,
) *gin.Engine {
	router := gin.New()

	api := router.Group("/api")

	pr := &parentRouter{
		tenderService: tenderService,
		bidService:    bidService,
		logger:        logger.Named("api"),
	}

	api.GET("/ping", func(ctx *gin.Context) { ctx.String(http.StatusOK, "ok") })

	bid.AttachToGroup(pr, api.Group("/bids"))
	tender.AttachToGroup(pr, api.Group("/tenders"))

	return router
}
