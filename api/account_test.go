package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simple-bank/db/mock"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/token"
	"github.com/techschool/simple-bank/util"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAccountApi(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountId     int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountForOwner(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "Unauthorized",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountForOwner(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				//require.Empty(t, recorder.Body)
			},
		},
		{
			name:      "No authorized",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountForOwner(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountForOwner(gomock.Any(), gomock.Any()).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			name:      "InternalError",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountForOwner(gomock.Any(), gomock.Any()).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name:      "InvalidId",
			accountId: 0,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountForOwner(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}

	for i := range testCases {

		func() {
			testCase := testCases[i]
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			testCase.buildStubs(store)

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", testCase.accountId)

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)

			testCase.checkResponse(t, recorder)
		}()

	}

}

func randomAccount(owner string) db.Account {

	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {

	data, err := io.ReadAll(body)

	require.NoError(t, err)

	var getAccount db.Account

	err = json.Unmarshal(data, &getAccount)

	require.NoError(t, err)

	require.Equal(t, account, getAccount)
}
