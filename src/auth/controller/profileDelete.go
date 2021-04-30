package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/postgres"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
**	/api/profile/delete
**	Удаление пользователя
**	В запросе должно содержаться поле passwd
**	авторизация в авторизационном хидере accessToken
**	-- Не протестировано
 */
func profileDelete(w http.ResponseWriter, r *http.Request) {
	passwd, Err := parsePasswordFromRequest(r)
	if Err != nil {
		logger.Error(r, Err)
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

	user, Err := postgres.UserGetBasicById(header.UserId)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	encryptedPass, Err := hash.PasswdHash(passwd)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if user.EncryptedPass == nil || *encryptedPass != *user.EncryptedPass {
		logger.Warning(r, "Хэши паролей не совпали. Ожидалось "+logger.BLUE+*user.EncryptedPass+logger.NO_COLOR+
			" получили "+logger.BLUE+*encryptedPass+logger.NO_COLOR)
		errorResponse(w, errors.ImpossibleToExecute.SetArgs("Пароль неверен", "Incorrect password"))
		return
	}

	if Err = postgres.UserDeleteBasic(user); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	successResponse(w, nil)
	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" успешно удален")
}

func parsePasswordFromRequest(r *http.Request) (string, *errors.Error) {
	type Password struct {
		Passwd *string `json:"passwd"`
	}
	var pass = Password{}
	if err := json.NewDecoder(r.Body).Decode(&pass); err != nil {
		return "", errors.InvalidRequestBody.SetOrigin(err)
	}
	if pass.Passwd == nil {
		return "", errors.NoArgument.SetArgs("passwd", "passwd")
	}
	return *pass.Passwd, nil
}
