package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAuthBasic(t *testing.T) {
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

	var user = &model.UserBasic{}
	email := "test@gmail.com"
	user.Email = &email
	user.Passwd = "qweRTY123@"
	user.EncryptedPass, Err = hash.PasswdHash(user.Passwd)
	if Err != nil {
		t.Errorf("%sError: не смог создать хэш пароля %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	user.Username = "displayname"

	if Err = postgres.UserSetBasic(user); Err != nil {
		t.Errorf("%sError: не смог создать тестового юзера %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	if Err = postgres.UserConfirmEmailBasic(user); Err != nil {
		t.Errorf("%sError: не смог подтвердить почту тестовому юзеру %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}

	defer func(user *model.UserBasic, t *testing.T) {
		if Err = postgres.UserDeleteBasic(user); Err != nil {
			t.Errorf("%sError: не смог удалить тестового юзера %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}
	}(user, t)

	/*
	**	Test cases. Main part of testing
	 */
	testCases := []struct {
		name           string
		email          string
		passwd         string
		expectedStatus int
	}{
		{
			name:           "wrong mail",
			email:          "email",
			passwd:         user.Passwd,
			expectedStatus: errors.AuthFail.GetHttpStatus(),
		}, {
			name:           "wrong passwd",
			email:          *user.Email,
			passwd:         "passwd",
			expectedStatus: errors.AuthFail.GetHttpStatus(),
		}, {
			name:           "valid",
			email:          *user.Email,
			passwd:         user.Passwd,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t_ *testing.T) {
			var (
				url = "http://localhost:" + strconv.Itoa(int(conf.ServerPort)) + "/api/auth/basic"
				rec = httptest.NewRecorder()
				req *http.Request
			)

			req = httptest.NewRequest("GET", url, nil)
			req.SetBasicAuth(tc.email, tc.passwd)
			/*
			**	start test
			 */
			authBasic(rec, req)
			if rec.Code != tc.expectedStatus {
				t_.Errorf("%sERROR: wrong StatusCode: got %d, expected %d%s", logger.RED_BG,
					rec.Code, tc.expectedStatus, logger.NO_COLOR)
			} else if tc.expectedStatus != http.StatusOK {
				t_.Logf("%sSUCCESS: user auth was failed as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Logf("%sSUCCESS: user was authenticated%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
	}
}
