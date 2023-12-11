package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
)

type depositRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createDeposit(ctx *gin.Context) {
	// Reading the request body
	var req depositRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the account by the provided ID
	if _, err := server.store.GetAccount(ctx, req.AccountID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Creating data to be set for the deposit
	arg := db.DepositTxParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
		User:      authPayload.Username,
	}

	// Calling the deposit transaction function
	result, err := server.store.DepositTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Returning data to the request if no error occurs
	ctx.JSON(http.StatusOK, result)
}

type getDepositRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getDeposit(ctx *gin.Context) {
	var req getDepositRequest
	// Here, we'll use the URL params
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deposit, err := server.store.GetDeposit(ctx, req.ID)
	if err != nil {
		// If no item was found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// Checking if user made the deposit
	if deposit.User != authPayload.Username {
		err := errors.New("deposit wasn't made by the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, deposit)
}

type listDepositRequest struct {
	AccountId int64 `form:"account_id" binding:"required,min=1"`
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listDeposits(ctx *gin.Context) {
	var req listDepositRequest
	// Here, we'll use the query params
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Checking if account belongs to the user
	account, err := server.store.GetAccount(ctx, req.AccountId)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Here, we'll get the query params
	arg := db.ListDepositsParams{
		AccountID: req.AccountId,
		Limit:     req.PageSize,
		// Calculating the offset from page number and size
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListDeposits(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
