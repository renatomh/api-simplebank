package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetDeposit(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	deposit := randomDeposit(account.ID, user.Username)

	testCases := []struct {
		name          string
		depositID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// Success case
		{
			name:      "OK",
			depositID: deposit.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We expect the function to be called once
				store.EXPECT().
					GetDeposit(gomock.Any(), gomock.Eq(deposit.ID)).
					Times(1).
					// In this case, we expect the deposit to be returned and no errors
					Return(deposit, nil)
			},
			// Checking response body and status code
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchDeposit(t, recorder.Body, deposit)
			},
		},
		// Unauthorized user
		{
			name:      "UnauthorizedUser",
			depositID: deposit.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We expect the function to be called once
				store.EXPECT().
					GetDeposit(gomock.Any(), gomock.Eq(deposit.ID)).
					Times(1).
					// In this case, we expect the deposit to be returned and no errors
					Return(deposit, nil)
			},
			// Checking response body and status code
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// No authorization
		{
			name:      "NoAuthorization",
			depositID: deposit.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We expect the function to be called once
				store.EXPECT().
					GetDeposit(gomock.Any(), gomock.Any()).
					Times(0)
			},
			// Checking response body and status code
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// Item not found case
		{
			name:      "NotFound",
			depositID: deposit.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We expect the function to be called once
				store.EXPECT().
					GetDeposit(gomock.Any(), gomock.Eq(deposit.ID)).
					Times(1).
					// In this case, we expect no deposit to be returned and a NoRows sql error
					Return(db.Deposit{}, sql.ErrNoRows)
			},
			// Checking response status code
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		// Internal server error case
		{
			name:      "InternalError",
			depositID: deposit.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We expect the function to be called once
				store.EXPECT().
					GetDeposit(gomock.Any(), gomock.Eq(deposit.ID)).
					Times(1).
					// In this case, we expect no deposit to be returned and an connection sql error
					Return(db.Deposit{}, sql.ErrConnDone)
			},
			// Checking response status code
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		// Bad request error case
		{
			name:      "InvalidID",
			depositID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We expect the function not to be called
				store.EXPECT().
					GetDeposit(gomock.Any(), gomock.Any()).
					Times(0)
			},
			// Checking response status code
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Creating a controller from mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Creating a store with the controller
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// Starting the mock server, recording the responses
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Defining the API URL, making the request and checking if everything is ok
			url := fmt.Sprintf("/deposits/%d", tc.depositID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			// Checking response status code
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDepositAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	amount := int64(10)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DepositTxParams{
					AccountID: account.ID,
					Amount:    amount,
					User:      user.Username,
				}
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
				store.EXPECT().DepositTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().DepositTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "AccountNotFound",
			body: gin.H{
				"account_id": 99999,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().DepositTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/deposits"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomDeposit(accountId int64, user string) db.Deposit {
	return db.Deposit{
		ID:        util.RandomInt(1, 1000),
		AccountID: accountId,
		Amount:    util.RandomMoney(),
		User:      user,
	}
}

func requireBodyMatchDeposit(t *testing.T, body *bytes.Buffer, deposit db.Deposit) {
	// Reading the response body
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// Checking if data matches
	var gotDeposit db.Deposit
	err = json.Unmarshal(data, &gotDeposit)
	require.NoError(t, err)
	require.Equal(t, deposit, gotDeposit)
}

func requireBodyMatchDeposits(t *testing.T, body *bytes.Buffer, deposits []db.Deposit) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotDeposits []db.Deposit
	err = json.Unmarshal(data, &gotDeposits)
	require.NoError(t, err)
	require.Equal(t, deposits, gotDeposits)
}
