package bid

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"avito2024/internal/app/core/entity"
	"avito2024/internal/app/core/service"
)

func (r *bidRouter) create(ctx *gin.Context) {
	var bid entity.Bid

	if err := ctx.Bind(&bid); err != nil {
		r.logger.Error("bind failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := r.bidService.Create(ctx, &bid)
	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				entity.ResponseError{Reason: service.ErrUserNotExists.Error()},
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
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				entity.ResponseError{Reason: err.Error()},
			)
			return
		}

		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			entity.ResponseError{Reason: err.Error()},
		)
		return
	}

	ctx.JSON(http.StatusOK, bid)

}
