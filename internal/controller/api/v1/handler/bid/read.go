package bid

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"avito2024/internal/app/core/entity"
	"avito2024/internal/app/core/service"
)

func (r *bidRouter) list(ctx *gin.Context) {
	tenderID := ctx.Param("id")

	userName := ctx.Query("username")

	bids, err := r.bidService.ListTenderBids(ctx, entity.TenderID(tenderID), userName)
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

		if errors.Is(err, service.ErrTenderOrBidNotFound) {
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				entity.ResponseError{Reason: service.ErrTenderOrBidNotFound.Error()},
			)
			return
		}
	}

	ctx.JSON(http.StatusOK, bids)
}

// func (r *bidRouter) listMy(ctx *gin.Context) {

// }

func (r *bidRouter) status(ctx *gin.Context) {
	bidID := ctx.Param("id")

	userName := ctx.Query("username")

	status, err := r.bidService.Status(ctx, entity.BidId(bidID), userName)
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

		if errors.Is(err, service.ErrBidNotFound) {
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				entity.ResponseError{Reason: service.ErrBidNotFound.Error()},
			)
			return
		}
	}

	ctx.JSON(http.StatusOK, status)
}
