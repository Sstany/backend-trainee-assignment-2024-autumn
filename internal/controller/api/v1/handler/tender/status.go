package tender

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"avito2024/internal/app/core/entity"
	"avito2024/internal/app/core/service"
)

func (r *tenderRouter) status(ctx *gin.Context) {
	tenderID := ctx.Param("tenderId")

	userName := ctx.Query("username")

	status, err := r.tenderService.GetStatus(ctx, entity.TenderID(tenderID), userName)
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

		if errors.Is(err, service.ErrTenderNotFound) {
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				entity.ResponseError{Reason: service.ErrTenderNotFound.Error()},
			)
			return
		}
	}

	ctx.JSON(http.StatusOK, status)
}
