package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/GritNg/simple-bank-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccId int64  `json:"from_acc_id" binding:"required,min=1"`
	ToAccId   int64  `json:"to_acc_id" binding:"required,min=1"`
	Amount    int64  `json:"amount" binding:"required,gt=0"`
	Currency  string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccId, req.Currency) {
		return
	}

	if server.validAccount(ctx, req.ToAccId, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccId,
		ToAccountId:   req.ToAccId,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch: %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	return true
}
