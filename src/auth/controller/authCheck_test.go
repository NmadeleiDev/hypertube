package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAuthCheck(t *testing.T) {
	/*
	**	Initialize server
	 */
	initTest(t)
	defer postgres.Close()

	conf, Err := getConfig()
	if Err != nil {
		t.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}

	/*
	**	Создаем тестового юзера от имени которого будет запрашивать данные
	 */
	var user = &model.UserBasic{}
	email := "test@gmail.com"
	user.Email = &email
	user.Passwd = "qweRTY123@"
	user.EncryptedPass, Err = hash.PasswdHash(user.Passwd)
	if Err != nil {
		t.Errorf("%sError: не смог создать хэш пароля %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	user.Username = "user1"

	token, Err := hash.CreateAccessToken(user)
	if Err != nil {
		t.Errorf("%sError: не смог создать тестовый авторизационный токен %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}

	/*
	**	Test cases. Main part of testing
	 */
	testCases := []struct {
		name           string
		token          string
		serverPasswd   string
		expectedStatus int
	}{
		{
			name:           "valid",
			token:          token,
			serverPasswd:   conf.ServerPasswd,
			expectedStatus: http.StatusOK,
		}, {
			name:           "invalid token",
			token:          "invalid",
			serverPasswd:   conf.ServerPasswd,
			expectedStatus: errors.InvalidToken.GetHttpStatus(),
		}, {
			name:           "invalid server password",
			token:          token,
			serverPasswd:   "invalid",
			expectedStatus: errors.UserNotLogged.GetHttpStatus(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t_ *testing.T) {
			var (
				url = "http://localhost:" + strconv.Itoa(int(conf.ServerPort)) + "/api/auth/check"
				rec = httptest.NewRecorder()
				req *http.Request
			)

			requestToken := &model.Token{
				AccessToken:  tc.token,
				ServerPasswd: tc.serverPasswd,
			}

			payloadJson, err := json.Marshal(requestToken)
			if err != nil {
				t_.Errorf("%sError: cannot marshal test case payload - %s%s", logger.RED_BG, err.Error(), logger.NO_COLOR)
				t_.FailNow()
			}

			req = httptest.NewRequest("GET", url, bytes.NewBufferString(string(payloadJson)))

			/*
			**	start test
			 */
			authCheck(rec, req)
			if rec.Code != tc.expectedStatus {
				t_.Errorf("%sERROR: wrong StatusCode: got %d, expected %d%s", logger.RED_BG,
					rec.Code, tc.expectedStatus, logger.NO_COLOR)
			} else if tc.expectedStatus != http.StatusOK {
				t_.Logf("%sSUCCESS: authorization check was failed as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Logf("%sSUCCESS: authorization checked well%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
	}
}
