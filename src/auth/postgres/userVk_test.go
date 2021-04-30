package postgres

import (
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"testing"
	"time"
)

func TestSetUserVk(t *testing.T) {
	var (
		user1      = &model.UserVk{}
		user2      = &model.UserVk{}
		userBasic1 *model.UserBasic
		userBasic2 *model.UserBasic
	)
	user1.UserVkId = 42
	user2.UserVkId = 21
	accessToken := "access_token"
	user1.AccessToken = &accessToken
	user2.AccessToken = nil
	t1 := time.Now()
	user1.ExpiresAt = &t1
	user2.ExpiresAt = nil

	initTest(t)

	defer func(t *testing.T, user1, user2 *model.UserVk) {
		t.Run("Delete test user #1", func(t_ *testing.T) {
			if userBasic1 == nil {
				t_.Logf("%sBasic user 1 not deleted because its nil%s", logger.YELLOW, logger.NO_COLOR)
				return
			}
			if Err := UserDeleteBasic(userBasic1); Err != nil {
				t_.Errorf("%sError: cannot delete test user - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			} else {
				t_.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
		t.Run("Delete test user #2", func(t_ *testing.T) {
			if userBasic2 == nil {
				t_.Logf("%sBasic user 2 not deleted because its nil%s", logger.YELLOW, logger.NO_COLOR)
				return
			}
			if Err := UserDeleteBasic(userBasic2); Err != nil {
				t_.Errorf("%sError: cannot delete test user - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			} else {
				t_.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
		t.Run("Close connection", func(t_ *testing.T) {
			if Err := Close(); Err != nil {
				t_.Errorf("%sError: cannot close connection - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			} else {
				t_.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
	}(t, user1, user2)

	t.Run("valid create user #1", func(t_ *testing.T) {
		if userB1, Err := UserSetVk(user1); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			userBasic1 = userB1
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid update user #1", func(t_ *testing.T) {
		newAccessToken := "new access_token"
		user1.AccessToken = &newAccessToken
		if Err := UserUpdateVk(user1); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid get user #1 by id", func(t_ *testing.T) {
		_, Err := UserGetVkById(0)
		if Err != nil {
			if errors.UserNotExist.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: user not exists as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid get user #1 by user id", func(t_ *testing.T) {
		newUser, Err := UserGetVkById(user1.UserVkId)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else if user1.UserId != newUser.UserId || newUser.AccessToken == nil ||
			newUser.ExpiresAt == nil || *user1.AccessToken != *newUser.AccessToken ||
			user1.ExpiresAt.Format(time.StampMilli) != newUser.ExpiresAt.Format(time.StampMilli) {
			t_.Errorf("%sError: received user differs from original%s\nexpected %#v\ngot %#v", logger.RED_BG, logger.NO_COLOR,
				user1.UserVkModel, newUser.UserVkModel)
			t_.Errorf("%s\n%s", user1.ExpiresAt.Format(time.StampMilli), newUser.ExpiresAt.Format(time.StampMilli))
		} else {
			t_.Logf("%sSuccess: user was received successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid create user #2", func(t_ *testing.T) {
		var Err *errors.Error
		userBasic2, Err = UserSetVk(user2)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid create user #1", func(t_ *testing.T) {
		if _, Err := UserSetVk(user1); Err != nil {
			if errors.ImpossibleToExecute.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: found error that was expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: expected %s found %s error%s", logger.RED_BG,
					errors.ImpossibleToExecute, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)

		}
	})

	t.Run("valid user delete #2", func(t_ *testing.T) {
		if userBasic2 == nil {
			t_.Errorf("%sError: cannot start test because basic user 2 is nil%s", logger.RED_BG, logger.NO_COLOR)
			t_.FailNow()
		}
		if Err := UserDeleteBasic(userBasic2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was deleted successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	// /!!!!!
	t.Run("invalid update user #2", func(t_ *testing.T) {
		newAccessToken := "new access_token"
		user2.AccessToken = &newAccessToken
		if Err := UserUpdateVk(user2); Err != nil {
			if errors.UserNotExist.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: user not exists as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid user delete #2", func(t_ *testing.T) {
		if userBasic2 == nil {
			t_.Errorf("%sError: cannot start test because basic user 2 is nil%s", logger.RED_BG, logger.NO_COLOR)
			t_.FailNow()
		}
		if Err := UserDeleteBasic(userBasic2); Err != nil {
			if errors.ImpossibleToExecute.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: found error that was expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: expected %s found %s error%s", logger.RED_BG,
					errors.ImpossibleToExecute, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid recreate user #2", func(t_ *testing.T) {
		var Err *errors.Error
		if userBasic2, Err = UserSetVk(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid recreate user #2", func(t_ *testing.T) {
		if _, Err := UserSetVk(user2); Err != nil {
			if errors.ImpossibleToExecute.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: found error that was expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: expected %s found %s error%s", logger.RED_BG,
					errors.ImpossibleToExecute, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)
		}
	})

	if !t.Failed() {
		t.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
	}
}
