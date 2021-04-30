package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestProfileGet(t *testing.T) {
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
	var user1 = &model.UserBasic{}
	email := "test@gmail.com"
	user1.Email = &email
	user1.Passwd = "qweRTY123@"
	user1.EncryptedPass, Err = hash.PasswdHash(user1.Passwd)
	if Err != nil {
		t.Errorf("%sError: не смог создать хэш пароля %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	user1.Username = "user1"

	if Err = postgres.UserSetBasic(user1); Err != nil {
		t.Errorf("%sError: не смог создать тестового юзера %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	if Err = postgres.UserConfirmEmailBasic(user1); Err != nil {
		t.Errorf("%sError: не смог подтвердить почту тестовому юзеру %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}

	token, Err := hash.CreateAccessToken(user1)
	if Err != nil {
		t.Errorf("%sError: не смог создать тестовый авторизационный токен %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}

	/*
	**	Создаем второго тестового юзера. Его поля мы будем запрашивать
	 */
	var user2 = &model.UserBasic{}
	email2 := "testUser2@gmail.com"
	user2.Email = &email2
	user2.Passwd = "qweRTY123@"
	user2.EncryptedPass, Err = hash.PasswdHash(user2.Passwd)
	if Err != nil {
		t.Errorf("%sError: не смог создать хэш пароля %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}
	user2.Username = "user2"

	if Err = postgres.UserSetBasic(user2); Err != nil {
		t.Errorf("%sError: не смог создать тестового юзера %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}

	/*
	**	Этот блок выполнится в конце работы программы. Удаление тестовых юзеров
	 */
	defer func(user1 *model.UserBasic, user2 *model.UserBasic, t *testing.T) {
		if Err = postgres.UserDeleteBasic(user1); Err != nil {
			t.Errorf("%sError: не смог удалить тестового юзера %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}
		if Err = postgres.UserDeleteBasic(user2); Err != nil {
			t.Errorf("%sError: не смог удалить тестового юзера %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}
	}(user1, user2, t)

	/*
	**	Test cases. Main part of testing
	 */
	testCases := []struct {
		name             string
		userId           uint
		expectedEmail    *string
		expectedUsername string
		expectedStatus   int
		token            string
	}{
		{
			name:             "invalid id (no such user)",
			userId:           user2.UserId * 2,
			expectedEmail:    nil,
			expectedUsername: "nil",
			expectedStatus:   errors.ImpossibleToExecute.GetHttpStatus(),
			token:            token,
		}, {
			name:             "valid profile of user #2",
			userId:           user2.UserId,
			expectedEmail:    nil,
			expectedUsername: user2.Username,
			expectedStatus:   http.StatusOK,
			token:            token,
		}, {
			name:             "valid self profile",
			userId:           0,
			expectedEmail:    user1.Email,
			expectedUsername: user1.Username,
			expectedStatus:   http.StatusOK,
			token:            token,
		}, {
			name:             "invalid - no token in header",
			userId:           0,
			expectedEmail:    nil,
			expectedUsername: "nil",
			expectedStatus:   errors.InvalidToken.GetHttpStatus(),
			token:            "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t_ *testing.T) {
			var (
				url = "http://localhost:" + strconv.Itoa(int(conf.ServerPort)) + "/api/profile/get"
				rec = httptest.NewRecorder()
				req *http.Request
			)

			if tc.userId != 0 {
				url += "?userId=" + strconv.Itoa(int(tc.userId))
			}
			req = httptest.NewRequest("GET", url, nil)
			req.Header.Add("accessToken", tc.token)

			/*
			**	start test
			 */
			profileGet(rec, req)
			if rec.Code != tc.expectedStatus {
				t_.Errorf("%sERROR: wrong StatusCode: got %d, expected %d%s", logger.RED_BG,
					rec.Code, tc.expectedStatus, logger.NO_COLOR)
			} else if tc.expectedStatus != http.StatusOK {
				t_.Logf("%sSUCCESS: getting profile failed as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				/*
				**	Handle response
				 */
				var responseUser *model.UserBasic
				err := json.NewDecoder(rec.Body).Decode(&responseUser)
				if err != nil {
					t_.Errorf("%sERROR: decoding response body error: %s%s", logger.RED_BG, err.Error(), logger.NO_COLOR)
					t_.FailNow()
				}
				if tc.expectedEmail != nil && responseUser.Email == nil {
					t_.Errorf("%sERROR: response user email is nil. Expected not nil%s", logger.RED_BG, logger.NO_COLOR)
				} else if tc.expectedEmail == nil && responseUser.Email != nil {
					t_.Errorf("%sERROR: response user email is '%s'. Expected nil%s", logger.RED_BG, *responseUser.Email, logger.NO_COLOR)
				} else if tc.expectedEmail != nil && responseUser.Email != nil && *responseUser.Email != *tc.expectedEmail {
					t_.Errorf("%sERROR: response user email differs. Expected %s got %s%s", logger.RED_BG,
						*tc.expectedEmail, *responseUser.Email, logger.NO_COLOR)
				}
				if responseUser.Username != tc.expectedUsername {
					t_.Errorf("%sERROR: response username differs. Expected %s got %s%s", logger.RED_BG,
						tc.expectedUsername, responseUser.Username, logger.NO_COLOR)
				}
				if !t_.Failed() {
					t_.Logf("%sSUCCESS: user profile was received%s", logger.GREEN_BG, logger.NO_COLOR)
				}
			}
		})
	}
}
