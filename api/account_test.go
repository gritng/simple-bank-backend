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

	mockdb "github.com/GritNg/simple-bank-backend/db/mock"
	db "github.com/GritNg/simple-bank-backend/db/sqlc"
	"github.com/GritNg/simple-bank-backend/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountId     int64
		buildstubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountId: account.ID,
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountId: account.ID,
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountId: account.ID,
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidId",
			accountId: 0,
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// add more testcase
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildstubs(store)

			// start test
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", tc.accountId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
