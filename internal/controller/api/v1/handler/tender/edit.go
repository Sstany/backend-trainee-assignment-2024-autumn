package tender

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"avito2024/internal/app/core/entity"
	"avito2024/internal/app/core/service"
)

func (r *tenderRouter) edit(ctx *gin.Context) {
	tenderID := ctx.Param("tenderId")

	userName := ctx.Query("username")

	var update entity.TenderUpdate

	if err := ctx.Bind(&update); err != nil {
		r.logger.Error("bind failed", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{Reason: service.ErrWrongInputFormat.Error()})
		return
	}

	updatedTender, err := r.tenderService.Edit(ctx, entity.TenderID(tenderID), &update, userName)
	if err != nil {
		r.logger.Error("update tender failed", zap.Error(err))

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

		if errors.Is(err, service.ErrWrongInputFormat) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{Reason: service.ErrWrongInputFormat.Error()})
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Reason: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTender)
}
