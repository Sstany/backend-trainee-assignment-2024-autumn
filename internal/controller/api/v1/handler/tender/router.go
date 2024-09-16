package tender

import (
	"avito2024/internal/app/core/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type tenderRouter struct {
	tenderService *service.TenderService
	logger        *zap.Logger
}

type serviceProvider interface {
	TenderService() *service.TenderService
	Logger() *zap.Logger
}

func AttachToGroup(sp serviceProvider, group *gin.RouterGroup) {
	tr := &tenderRouter{
		tenderService: sp.TenderService(),
		logger:        sp.Logger().Named("tender"),
	}

	group.POST("/new", tr.create)
	group.GET("/", tr.list)
	group.GET("/my", tr.listMy)
	group.GET("/:tenderId/status", tr.status)
	group.PUT("/:tenderId/status", tr.updateStatus)
	group.PATCH("/:tenderId/edit", tr.edit)
	// group.PUT("/:tenderId/rollback/:version", tr.rollback)
}
