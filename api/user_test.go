package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/leemingen1227/couple-server/db/mock"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/techschool/simplebank/util"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.PasswordDigest)
	if err != nil {
		return false
	}

	e.arg.PasswordDigest = arg.PasswordDigest
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					ID: 		 user.ID,
					Name: 	  user.Name,
					Email: 	 user.Email,

				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(db.CreateUserTxResult{
						User: user,
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUserID",
			body: gin.H{
				"id":        "invalid-uuid",
				"name":      user.Name,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"id":        user.ID,
				"name":      user.Name,
				"email":     "invalid email",
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"password":  "123",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/users/signup"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// func TestLoginUserAPI(t *testing.T) {
// 	user, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"username": user.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(user.Username)).
// 					Times(1).
// 					Return(user, nil)
// 				store.EXPECT().
// 					CreateSession(gomock.Any(), gomock.Any()).
// 					Times(1)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "UserNotFound",
// 			body: gin.H{
// 				"username": "NotFound",
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.User{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "IncorrectPassword",
// 			body: gin.H{
// 				"username": user.Username,
// 				"password": "incorrect",
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(user.Username)).
// 					Times(1).
// 					Return(user, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			body: gin.H{
// 				"username": user.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.User{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidUsername",
// 			body: gin.H{
// 				"username":  "invalid-user#1",
// 				"password":  password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/login"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	uid, err := uuid.NewRandom()

	user = db.User{
		ID:             uid,
		Email:          util.RandomEmail(),
		Name:           util.RandomOwner(),
		PasswordDigest: hashedPassword,
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.ID, gotUser.ID)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.PasswordDigest)
}