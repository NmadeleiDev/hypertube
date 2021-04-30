package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/controller/validator"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/postgres"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
**	/api/passwd/patch
**	Обновление пароля пользователя
**	В запросе должны содержаться поля passwd, newPasswd
**	авторизация в авторизационном хидере accessToken
**	-- Проверено
 */
func passwdPatch(w http.ResponseWriter, r *http.Request) {
	passwd, newPasswd, Err := parsePasswordsFromRequest(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		errorResponse(w, Err)
		return
	}

	accessToken := r.Header.Get("accessToken")
	if accessToken == "" {
		logger.Error(r, errors.UserNotLogged.SetArgs("отсутствует токен доступа", "access token expected"))
		errorResponse(w, errors.UserNotLogged)
		return
	}

	header, Err := hash.GetHeaderFromAccessToken(accessToken)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if Err = validator.ValidatePassword(newPasswd); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user, Err := postgres.UserGetBasicById(header.UserId)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if user.EncryptedPass == nil {
		logger.Log(r, "У пользователя #"+strconv.Itoa(int(header.UserId))+
			" не задан пароль, поэтому не делаю проверку на совпадение старого пароля")
	} else {
		if encryptedPass, Err := hash.PasswdHash(passwd); Err != nil {
			logger.Error(r, Err)
			errorResponse(w, Err)
			return
		} else if *encryptedPass != *user.EncryptedPass {
			logger.Warning(r, "Хэши паролей не совпали. Ожидалось "+logger.BLUE+*user.EncryptedPass+logger.NO_COLOR+
				" получили "+logger.BLUE+*encryptedPass+logger.NO_COLOR)
			errorResponse(w, errors.ImpossibleToExecute.SetArgs("Пароль неверен", "Incorrect password"))
			return
		} else {
			logger.Log(r, "Хэши паролей совпали. Продолжаю выполнение")
		}
	}

	user.EncryptedPass, Err = hash.PasswdHash(newPasswd)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if Err = postgres.UserUpdateEncryptedPassBasic(user); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	successResponse(w, nil)
	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" успешно обновил свои поля")
}

func parsePasswordsFromRequest(r *http.Request) (string, string, *errors.Error) {
	type Passwords struct {
		NewPasswd *string `json:"newPasswd"`
		Passwd    *string `json:"passwd"`
	}
	var pass = Passwords{}
	if err := json.NewDecoder(r.Body).Decode(&pass); err != nil {
		return "", "", errors.InvalidRequestBody.SetOrigin(err)
	}
	if pass.NewPasswd == nil {
		return "", "", errors.NoArgument.SetArgs("newPasswd", "newPasswd")
	}
	if pass.Passwd == nil {
		return "", "", errors.NoArgument.SetArgs("passwd", "passwd")
	}
	return *pass.Passwd, *pass.NewPasswd, nil
}
