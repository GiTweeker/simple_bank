package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/token"
	"net/http"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency"  binding:"required,currency`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var rq CreateAccountRequest
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	authPayload := ctx.MustGet(authorisationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: rq.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var rq getAccountRequest
	if err := ctx.ShouldBindUri(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}
	authPayload := ctx.MustGet(authorisationPayloadKey).(*token.Payload)

	arg := db.GetAccountForOwnerParams{
		ID:    rq.ID,
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

	/*if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to authenticated user")
		ctx.JSON(http.StatusOK, errorResponse(err))
		return
	}*/

	ctx.JSON(http.StatusOK, account)

}

type listAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var rq listAccountsRequest
	if err := ctx.ShouldBindQuery(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}
	authPayload := ctx.MustGet(authorisationPayloadKey).(*token.Payload)
	arg := db.ListAccountsForOwnerParams{
		Owner:  authPayload.Username,
		Limit:  rq.PageSize,
		Offset: (rq.PageId - 1) * rq.PageId,
	}

	accounts, err := server.store.ListAccountsForOwner(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)

}
