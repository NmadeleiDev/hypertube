package hash

import (
	"auth_backend/logger"
	"auth_backend/model"
	"testing"
)

func TestAccessToken(t *testing.T) {
	initializePackageForTest(t)

	t.Run("check for correct signature", func(t *testing.T) {
		var user = &model.UserBasic{}
		user.UserId = 42
		user.Email = "school21@gmail.com"
		accessToken, Err := CreateAccessToken(user)
		if Err != nil {
			t.Errorf("%sError during creating token - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}

		if Err = CheckAccessTokenSignature(accessToken); Err != nil {
			t.Errorf("%sError cannot unmarshal token - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}
		t.Logf("%sSuccess: token is valid%s", logger.GREEN_BG, logger.NO_COLOR)
	})

	t.Run("check token header data validity", func(t *testing.T) {
		var user = &model.UserBasic{}
		user.UserId = 42
		user.Email = "school21@gmail.com"
		imageBody := "image_body"
		user.ImageBody = &imageBody
		user.Username = "skinnyman"
		fname := "Den"
		user.Fname = &fname
		lname := "QWERTY"
		user.Lname = &lname
		accessToken, Err := CreateAccessToken(user)
		if Err != nil {
			t.Errorf("%sError during creating token - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}

		header, Err := GetHeaderFromAccessToken(accessToken)
		if Err != nil {
			t.Errorf("%sError cannot unmarshal token - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}

		if header.UserId != user.UserId {
			t.Errorf("%sError: UserId are incorrect after decoding. Expected %d Got %d%s", logger.RED_BG,
				user.UserId, header.UserId, logger.NO_COLOR)
		}

		if !t.Failed() {
			t.Logf("%sSuccess: token is valid%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})
}
