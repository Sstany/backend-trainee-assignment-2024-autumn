package tender

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"avito2024/internal/app/core/entity"
	"avito2024/internal/app/core/service"
)

func (r *tenderRouter) updateStatus(ctx *gin.Context) {
	tenderID := ctx.Param("tenderId")

	userName := ctx.Query("username")

	status := ctx.Query("status")
	if status == "" {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			entity.ResponseError{Reason: "status malformed"},
		)
		return
	}

	tender, err := r.tenderService.SetStatus(ctx, entity.TenderID(tenderID), userName, entity.TenderStatus(status))
	if err != nil {
		r.logger.Error("set status failed", zap.Error(err))
		if errors.Is(err, service.ErrUserNotExists) {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				entity.ResponseError{Reason: service.ErrUserNotExists.Error()},
			)
			return
		}

		if errors.Is(err, service.ErrNotEnoughRights) {
			ctx.AbortWithStatusJSON(
				http.StatusForbidden,
				entity.ResponseError{Reason: service.ErrNotEnoughRights.Error()},
			)
			return
		}

		if errors.Is(err, service.ErrTenderNotFound) {
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				entity.ResponseError{Reason: service.ErrTenderNotFound.Error()},
			)
			return
		}
	}

	ctx.JSON(http.StatusOK, tender)
}
