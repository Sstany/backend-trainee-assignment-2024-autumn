package tender

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"avito2024/internal/app/core/entity"
	"avito2024/internal/app/core/service"
)

func (r *tenderRouter) create(ctx *gin.Context) {
	var tender entity.RequestTender

	if err := ctx.Bind(&tender); err != nil {
		r.logger.Error("bind failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := r.tenderService.Create(ctx, &tender.Tender, tender.Username)
	if err != nil {
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

		r.logger.Error("failed to create tender", zap.Any("reqBody", tender), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{
			Reason: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, tender)
}
