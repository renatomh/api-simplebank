package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) validAccountCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	// Getting the account by the provided ID
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	// Checking if currency matches
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		// Sending error response to the client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

func (server *Server) createTransfer(ctx *gin.Context) {
	// Reading the request body
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Checking if accounts have the required currency
	fromAccount, valid := server.validAccountCurrency(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("origin account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccountCurrency(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	// Creating data to be set for the transfer
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// Calling the transfer transaction function
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Returning data to the request if no error occurs
	ctx.JSON(http.StatusOK, result)
}
