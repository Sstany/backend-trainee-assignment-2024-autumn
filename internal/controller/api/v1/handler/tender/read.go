package tender

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"avito2024/internal/app/core/entity"
	"avito2024/internal/app/core/service"
)

func (r *tenderRouter) listMy(ctx *gin.Context) {
	userName := ctx.Query("username")
	limit := ctx.Query("limit")
	offset := ctx.Query("offset")

	if userName == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{
			Reason: "username query param is empty",
		})
		return
	}

	limitOffset := entity.ParseRequestLimitOffset(limit, offset)
	if limitOffset == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{
			Reason: "cannot parse limit or offset",
		})
		return
	}

	tenders, err := r.tenderService.ListMy(ctx, userName, limitOffset)
	if err != nil {
		r.logger.Error("failed to list users tenders", zap.String("username", userName), zap.Error(err))

		if errors.Is(err, service.ErrUserNotExists) {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				entity.ResponseError{Reason: service.ErrUserNotExists.Error()},
			)
			return
		}

		if errors.Is(err, service.ErrTenderOrBidNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound,
				entity.ResponseError{
					Reason: err.Error(),
				})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{
			Reason: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, tenders)
}

func (r *tenderRouter) list(ctx *gin.Context) {
	serviceTypes, _ := ctx.GetQueryArray("service_type")

	limit := ctx.Query("limit")
	offset := ctx.Query("offset")

	limitOffset := entity.ParseRequestLimitOffset(limit, offset)
	if limitOffset == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ResponseError{
			Reason: "cannot parse limit or offset",
		})
		return
	}

	var tenderServiceType []entity.TenderServiceType

	for _, service := range serviceTypes {
		tenderServiceType = append(tenderServiceType, entity.TenderServiceType(service))
	}

	tenders, err := r.tenderService.List(ctx, tenderServiceType, limitOffset)
	if err != nil {
		if errors.Is(err, service.ErrWrongInputFormat) {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				entity.ResponseError{Reason: service.ErrWrongInputFormat.Error()},
			)
			return
		}

		if errors.Is(err, service.ErrTenderNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound,
				entity.ResponseError{
					Reason: err.Error(),
				})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{
			Reason: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, tenders)
}
