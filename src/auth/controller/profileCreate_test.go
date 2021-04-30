package controller

import (
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestProfileCreate(t *testing.T) {
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
	**	Test cases. Main part of testing
	 */
	testCases := []struct {
		name           string
		payload        map[string]interface{}
		requestBody    *io.Reader
		expectedStatus int
	}{
		{
			name: "invalid mail",
			payload: map[string]interface{}{
				"email":    "email",
				"passwd":   "passWD123@",
				"username": "user",
			},
			expectedStatus: errors.InvalidArgument.GetHttpStatus(),
		}, {
			name: "invalid password",
			payload: map[string]interface{}{
				"email":    "email@gmail.com",
				"passwd":   "pass",
				"username": "user",
			},
			expectedStatus: errors.InvalidArgument.GetHttpStatus(),
		}, {
			name: "invalid displayname",
			payload: map[string]interface{}{
				"email":    "email@gmail.com",
				"passwd":   "passWD123@",
				"username": "",
			},
			expectedStatus: errors.InvalidArgument.GetHttpStatus(),
		}, {
			name: "valid",
			payload: map[string]interface{}{
				"email":    "email@gmail.com",
				"passwd":   "passWD123@",
				"username": "user",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t_ *testing.T) {
			var (
				url = "http://localhost:" + strconv.Itoa(int(conf.ServerPort)) + "/api/profile/create"
				rec = httptest.NewRecorder()
				req *http.Request
			)

			payloadJson, err := json.Marshal(tc.payload)
			if err != nil {
				t_.Errorf("%sError: cannot marshal test case payload - %s%s", logger.RED_BG, err.Error(), logger.NO_COLOR)
				t_.FailNow()
			}

			// tc.requestBody = bytes.NewBufferString(tc.payload)
			req = httptest.NewRequest("PUT", url, bytes.NewBufferString(string(payloadJson)))

			/*
			**	start test
			 */
			profileCreate(rec, req)
			if rec.Code != tc.expectedStatus {
				t_.Errorf("%sERROR: wrong StatusCode: got %d, expected %d%s", logger.RED_BG,
					rec.Code, tc.expectedStatus, logger.NO_COLOR)
			} else if tc.expectedStatus != http.StatusOK {
				t_.Logf("%sSUCCESS: user create was failed as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Logf("%sSUCCESS: user was created%s", logger.GREEN_BG, logger.NO_COLOR)
			}
			/*
			**	Delete user if it was created
			 */
			if rec.Code == http.StatusOK {
				/*
				**	Handle response
				 */
				var user *model.UserBasic
				err := json.NewDecoder(rec.Body).Decode(&user)
				if err != nil {
					t_.Errorf("%sERROR: Cannot read response body - %s%s", logger.RED_BG, err.Error(), logger.NO_COLOR)
					t_.FailNow()
				}
				if user != nil && user.UserId != 0 {
					if Err := postgres.UserDeleteBasic(user); Err != nil {
						t_.Errorf("%sERROR: Cannot delete test user - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
						t_.FailNow()
					}
				} else {
					t_.Errorf("%sError: user from response is invalid%s\n%#v", logger.RED_BG, logger.NO_COLOR, user)
				}
			}

			if !t_.Failed() {
				t_.Logf("%sSuccess - test user was deleted%s", logger.GREEN, logger.NO_COLOR)
			}
		})
	}
}
