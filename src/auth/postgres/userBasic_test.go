package postgres

import (
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"testing"
)

func TestSetUserBasic(t *testing.T) {
	var (
		user1 = &model.UserBasic{}
		user2 = &model.UserBasic{}
	)
	user1.Email = &emailValid1
	user1.EncryptedPass = &encryptedPass
	fname := "Denis"
	user1.Fname = &fname
	lname := "Globchansky"
	user1.Lname = &lname
	user1.Username = username
	user2.Email = &emailValid2
	user2.EncryptedPass = &encryptedPass
	user2.Username = username

	initTest(t)

	defer func(t *testing.T, user1, user2 *model.UserBasic) {
		t.Run("Delete test user #1", func(t_ *testing.T) {
			if Err := UserDeleteBasic(user1); Err != nil {
				t_.Errorf("%sError: cannot delete test user - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			} else {
				t_.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
		t.Run("Delete test user #2", func(t_ *testing.T) {
			if Err := UserDeleteBasic(user2); Err != nil {
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
		if Err := UserSetBasic(user1); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid update basic user #1", func(t_ *testing.T) {
		if Err := UserUpdateBasic(user1); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was updated successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid update password basic user #1", func(t_ *testing.T) {
		if Err := UserUpdateEncryptedPassBasic(user1); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was updated successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid get user #1 by email", func(t_ *testing.T) {
		if user1.Email == nil {
			t_.Errorf("%sFail: email is nil. Skip test%s", logger.RED_BG, logger.NO_COLOR)
			t_.FailNow()
		}
		newUser, Err := UserGetBasicByEmail(*user1.Email)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			if user1.UserId != newUser.UserId {
				t_.Errorf("%sError: received users Id differs. Expected %d got %d%s", logger.RED_BG,
					user1.UserId, newUser.UserId, logger.NO_COLOR)
			}
			if newUser.Email == nil {
				t_.Errorf("%sFail: email of new user is nil. Skip test%s", logger.RED_BG, logger.NO_COLOR)
				t_.FailNow()
			}
			if *user1.Email != *newUser.Email {
				t_.Errorf("%sError: received users Email differs. Expected %s got %s%s", logger.RED_BG,
					*user1.Email, *newUser.Email, logger.NO_COLOR)
			}
			if (user1.EncryptedPass != nil && newUser.EncryptedPass == nil) ||
				(user1.EncryptedPass == nil && newUser.EncryptedPass != nil) {
				t_.Errorf("%sError: received users EncryptedPass differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.EncryptedPass, newUser.EncryptedPass, logger.NO_COLOR)
			} else if user1.EncryptedPass != nil && newUser.EncryptedPass != nil &&
				*user1.EncryptedPass != *newUser.EncryptedPass {
				t_.Errorf("%sError: received users EncryptedPass differs. Expected %s got %s%s", logger.RED_BG,
					*user1.EncryptedPass, *newUser.EncryptedPass, logger.NO_COLOR)
			}
			if (user1.Fname != nil && newUser.Fname == nil) ||
				(user1.Fname == nil && newUser.Fname != nil) {
				t_.Errorf("%sError: received users Fname differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.Fname, newUser.Fname, logger.NO_COLOR)
			} else if user1.Fname != nil && newUser.Fname != nil && *user1.Fname != *newUser.Fname {
				t_.Errorf("%sError: received users Fname differs. Expected %s got %s%s", logger.RED_BG,
					*user1.Fname, *newUser.Fname, logger.NO_COLOR)
			}
			if (user1.Lname != nil && newUser.Lname == nil) ||
				(user1.Lname == nil && newUser.Lname != nil) {
				t_.Errorf("%sError: received users Lname differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.Lname, newUser.Lname, logger.NO_COLOR)
			} else if user1.Lname != nil && newUser.Lname != nil && *user1.Lname != *newUser.Lname {
				t_.Errorf("%sError: received users Lname differs. Expected %s got %s%s", logger.RED_BG,
					*user1.Lname, *newUser.Lname, logger.NO_COLOR)
			}
			if (user1.ImageBody != nil && newUser.ImageBody == nil) ||
				(user1.ImageBody == nil && newUser.ImageBody != nil) {
				t_.Errorf("%sError: received users ImageBody differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.ImageBody, newUser.ImageBody, logger.NO_COLOR)
			} else if user1.ImageBody != nil && newUser.ImageBody != nil && *user1.ImageBody != *newUser.ImageBody {
				t_.Errorf("%sError: received users ImageBody differs. Expected %s got %s%s", logger.RED_BG,
					*user1.ImageBody, *newUser.ImageBody, logger.NO_COLOR)
			}
			if user1.Username != newUser.Username {
				t_.Errorf("%sError: received users Username differs. Expected %s got %s%s", logger.RED_BG,
					user1.Username, newUser.Username, logger.NO_COLOR)
			}
			if (user1.NewEmail != nil && newUser.NewEmail == nil) ||
				(user1.NewEmail == nil && newUser.NewEmail != nil) {
				t_.Errorf("%sError: received users NewEmail differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.NewEmail, newUser.NewEmail, logger.NO_COLOR)
			} else if user1.NewEmail != nil && newUser.NewEmail != nil && *user1.NewEmail != *newUser.NewEmail {
				t_.Errorf("%sError: received users NewEmail differs. Expected %s got %s%s", logger.RED_BG,
					*user1.NewEmail, *newUser.NewEmail, logger.NO_COLOR)
			}
		}
		if !t_.Failed() {
			t_.Logf("%sSuccess: user was received successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid get user #1 by email", func(t_ *testing.T) {
		_, Err := UserGetBasicByEmail("not existing email")
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

	t.Run("invalid get user #1 by id", func(t_ *testing.T) {
		_, Err := UserGetBasicById(0)
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
		newUser, Err := UserGetBasicById(user1.UserId)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			if user1.UserId != newUser.UserId {
				t_.Errorf("%sError: received users Id differs. Expected %d got %d%s", logger.RED_BG,
					user1.UserId, newUser.UserId, logger.NO_COLOR)
			}
			if newUser.Email == nil || user1.Email == nil {
				t_.Errorf("%sFail: email is nil. Skip test%s", logger.RED_BG, logger.NO_COLOR)
				t_.FailNow()
			}
			if *user1.Email != *newUser.Email {
				t_.Errorf("%sError: received users Email differs. Expected %s got %s%s", logger.RED_BG,
					*user1.Email, *newUser.Email, logger.NO_COLOR)
			}
			if (user1.EncryptedPass != nil && newUser.EncryptedPass == nil) ||
				(user1.EncryptedPass == nil && newUser.EncryptedPass != nil) {
				t_.Errorf("%sError: received users EncryptedPass differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.EncryptedPass, newUser.EncryptedPass, logger.NO_COLOR)
			} else if user1.EncryptedPass != nil && newUser.EncryptedPass != nil &&
				*user1.EncryptedPass != *newUser.EncryptedPass {
				t_.Errorf("%sError: received users EncryptedPass differs. Expected %s got %s%s", logger.RED_BG,
					*user1.EncryptedPass, *newUser.EncryptedPass, logger.NO_COLOR)
			}
			if (user1.Fname != nil && newUser.Fname == nil) ||
				(user1.Fname == nil && newUser.Fname != nil) {
				t_.Errorf("%sError: received users Fname differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.Fname, newUser.Fname, logger.NO_COLOR)
			} else if user1.Fname != nil && newUser.Fname != nil && *user1.Fname != *newUser.Fname {
				t_.Errorf("%sError: received users Fname differs. Expected %s got %s%s", logger.RED_BG,
					*user1.Fname, *newUser.Fname, logger.NO_COLOR)
			}
			if (user1.Lname != nil && newUser.Lname == nil) ||
				(user1.Lname == nil && newUser.Lname != nil) {
				t_.Errorf("%sError: received users Lname differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.Lname, newUser.Lname, logger.NO_COLOR)
			} else if user1.Lname != nil && newUser.Lname != nil && *user1.Lname != *newUser.Lname {
				t_.Errorf("%sError: received users Lname differs. Expected %s got %s%s", logger.RED_BG,
					*user1.Lname, *newUser.Lname, logger.NO_COLOR)
			}
			if (user1.ImageBody != nil && newUser.ImageBody == nil) ||
				(user1.ImageBody == nil && newUser.ImageBody != nil) {
				t_.Errorf("%sError: received users ImageBody differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.ImageBody, newUser.ImageBody, logger.NO_COLOR)
			} else if user1.ImageBody != nil && newUser.ImageBody != nil && *user1.ImageBody != *newUser.ImageBody {
				t_.Errorf("%sError: received users ImageBody differs. Expected %s got %s%s", logger.RED_BG,
					*user1.ImageBody, *newUser.ImageBody, logger.NO_COLOR)
			}
			if user1.Username != newUser.Username {
				t_.Errorf("%sError: received users Username differs. Expected %s got %s%s", logger.RED_BG,
					user1.Username, newUser.Username, logger.NO_COLOR)
			}
			if (user1.NewEmail != nil && newUser.NewEmail == nil) ||
				(user1.NewEmail == nil && newUser.NewEmail != nil) {
				t_.Errorf("%sError: received users NewEmail differs. Expected %#v got %#v%s", logger.RED_BG,
					user1.NewEmail, newUser.NewEmail, logger.NO_COLOR)
			} else if user1.NewEmail != nil && newUser.NewEmail != nil && *user1.NewEmail != *newUser.NewEmail {
				t_.Errorf("%sError: received users NewEmail differs. Expected %s got %s%s", logger.RED_BG,
					*user1.NewEmail, *newUser.NewEmail, logger.NO_COLOR)
			}
		}
		if !t_.Failed() {
			t_.Logf("%sSuccess: user was received successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid create user #2", func(t_ *testing.T) {
		if Err := UserSetBasic(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid create user #1", func(t_ *testing.T) {
		if Err := UserSetBasic(user1); Err != nil {
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
		if Err := UserDeleteBasic(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was deleted successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid user delete #2", func(t_ *testing.T) {
		if Err := UserDeleteBasic(user2); Err != nil {
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
		if Err := UserSetBasic(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid recreate user #2", func(t_ *testing.T) {
		if Err := UserSetBasic(user2); Err != nil {
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

	t.Run("valid user confirm email #2", func(t_ *testing.T) {
		if Err := UserConfirmEmailBasic(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was updated successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid get user #2 by id, check update status", func(t_ *testing.T) {
		newUser, Err := UserGetBasicById(user2.UserId)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else if newUser.IsEmailConfirmed != true {
			t_.Errorf("%sError: user email confirm status didnt change%s", logger.RED_BG, logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was received and checked email confirm status%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	if !t.Failed() {
		t.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
	}
}
