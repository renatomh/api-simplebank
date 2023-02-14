package api

import (
	"database/sql"
	"net/http"

	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	// Getting the request data
	var req createAccountRequest
	// Here, we'll use the body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If something goes wrong, we send a response to the user
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting args provided on the request body
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	// Creating the account
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	// Here, we'll use the URL params
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// If no item was found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest
	// Here, we'll use the query params
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Here, we'll get the query params
	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		// Calculating the offset from page number and size
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
