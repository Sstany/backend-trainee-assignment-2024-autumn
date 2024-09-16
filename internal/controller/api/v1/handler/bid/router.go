package bid

import (
	"avito2024/internal/app/core/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bidRouter struct {
	bidService *service.BidService
	logger     *zap.Logger
}

type serviceProvider interface {
	BidService() *service.BidService
	Logger() *zap.Logger
}

func AttachToGroup(sp serviceProvider, group *gin.RouterGroup) {
	br := &bidRouter{
		bidService: sp.BidService(),
		logger:     sp.Logger().Named("bid"),
	}

	group.POST("/new", br.create)
	// group.GET("/my", br.read)
	group.GET("/:id/list", br.list)
	group.GET("/:id/status", br.status)

	//group.PUT("/:id", br.update)
}
