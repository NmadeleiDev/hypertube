package hash

import (
	"auth_backend/logger"
	"auth_backend/model"
	"testing"
)

func TestEmailToken(t *testing.T) {
	initializePackageForTest(t)

	t.Run("check for correct signature", func(t *testing.T) {
		headerString := "some string to test"
		signature1, Err := createEmailTokenSignature(headerString)
		if Err != nil {
			t.Errorf("%sError during creating signature - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}

		signature2, Err := createEmailTokenSignature(headerString)
		if Err != nil {
			t.Errorf("%sError during creating signature - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}

		if len(signature1) < 7 {
			t.Errorf("%sToo short signature: length %d '%s'%s", logger.RED_BG, len(signature1), signature1, logger.NO_COLOR)
		}

		if len(signature2) < 7 {
			t.Errorf("%sToo short signature: length %d '%s'%s", logger.RED_BG, len(signature2), signature2, logger.NO_COLOR)
		}

		if signature1 != signature2 {
			t.Errorf("%sSignature missmatch: %s != %s%s", logger.RED_BG, signature1, signature2, logger.NO_COLOR)
		}

		if Err = checkEmailTokenPartsSignature(headerString, signature1); Err != nil {
			t.Errorf("%sCannot check token parts - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}

		if !t.Failed() {
			t.Logf("%sSuccess: email token is valid%s", logger.GREEN_BG, logger.NO_COLOR)
		}
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
		newEmail := "newEmail@gmail.com"
		user.NewEmail = &newEmail
		emailToken, Err := CreateEmailToken(user)
		if Err != nil {
			t.Errorf("%sError during creating token - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}

		header, Err := GetHeaderFromEmailToken(emailToken)
		if Err != nil {
			t.Errorf("%sError cannot unmarshal token - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t.FailNow()
		}

		if header.UserId != user.UserId {
			t.Errorf("%sError: UserId are incorrect after decoding. Expected %d Got %d%s", logger.RED_BG,
				user.UserId, header.UserId, logger.NO_COLOR)
		}

		if header.NewEmail != *user.NewEmail {
			t.Errorf("%sError: NewEmail are incorrect after decoding. Expected %s Got %s%s", logger.RED_BG,
				*user.NewEmail, header.NewEmail, logger.NO_COLOR)
		}

		if !t.Failed() {
			t.Logf("%sSuccess: email token is valid%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})
}
