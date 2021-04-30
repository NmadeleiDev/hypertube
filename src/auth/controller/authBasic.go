package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
**	/api/auth/basic
**	-- Проверено
 */

func authBasic(w http.ResponseWriter, r *http.Request) {
	email, passwd, ok := r.BasicAuth()
	if !ok {
		logger.Warning(r, "authenticaion failed - email or password not found")
		errorResponse(w, errors.NoArgument.SetArgs("Отсутствует авторизационное поле", "Authorization field expected"))
		return
	}

	user, Err := postgres.UserGetBasicByEmail(email)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	encryptedPass, Err := hash.PasswdHash(passwd)
	if Err != nil {
		logger.Warning(r, "cannot get password hash - "+Err.Error())
		errorResponse(w, Err)
		return
	}

	if user.EncryptedPass == nil || encryptedPass == nil {
		logger.Warning(r, "authenticaion failed - password missmatch")
		errorResponse(w, errors.AuthFail)
		return
	}

	if *user.EncryptedPass != *encryptedPass {
		logger.Warning(r, "authenticaion failed - password missmatch. Expected "+*user.EncryptedPass+" got "+*encryptedPass)
		errorResponse(w, errors.AuthFail)
		return
	}

	if user.IsEmailConfirmed == false {
		logger.Warning(r, "authenticaion failed - email of user is not confirmed")
		errorResponse(w, errors.NotConfirmedMail)
		return
	}

	accessToken, Err := hash.CreateAccessToken(user)
	if Err != nil {
		logger.Warning(r, "cannot get password hash - "+Err.Error())
		errorResponse(w, Err)
		return
	}

	var responseToken = model.Token{AccessToken: accessToken}

	responseJson, err := json.Marshal(responseToken)
	if err != nil {
		logger.Error(r, errors.MarshalError.SetOrigin(err))
		errorResponse(w, errors.MarshalError)
		return
	}

	successResponse(w, responseJson)
	logger.Success(r, "user #"+strconv.Itoa(int(user.UserId))+" was authenticated")
}
