package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/token"
	"net/http"
)

type CreateTransferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id"  binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var rq CreateTransferRequest
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	var valid bool
	_, valid = server.validateAccountForOwner(ctx, rq.ToAccountId, rq.Currency)

	if !valid {
		return
	}
	_, valid = server.validateAccount(ctx, rq.ToAccountId, rq.Currency)

	if !valid {
		return
	}
	arg := db.TransferTxParams{
		FromAccountId: rq.FromAccountId,
		ToAccountId:   rq.ToAccountId,
		Amount:        rq.Amount,
	}

	transferTxResult, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, transferTxResult)

}
func (server *Server) validateAccount(ctx *gin.Context, accountId int64, currency string) (account db.Account, valid bool) {

	account, err := server.store.GetAccount(ctx, accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account.Currency != currency {
		errMessage := fmt.Errorf("account (%d) currency mismatch : %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(errMessage))
		return account, false

	}

	return account, true
}
func (server *Server) validateAccountForOwner(ctx *gin.Context, accountId int64, currency string) (account db.Account, valid bool) {

	authPayload := ctx.MustGet(authorisationPayloadKey).(*token.Payload)

	arg := db.GetAccountForOwnerParams{
		ID:    accountId,
		Owner: authPayload.Username,
	}
	account, err := server.store.GetAccountForOwner(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account.Currency != currency {
		errMessage := fmt.Errorf("account (%d) currency mismatch : %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(errMessage))
		return account, false

	}

	return account, true
}
