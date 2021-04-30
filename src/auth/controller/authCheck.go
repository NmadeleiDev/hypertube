package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"encoding/json"
	"net/http"
)

/*
**	Проверка авторизации - этот эндпоинт должен дергать только Гриша
**	-- Проверено
 */
func authCheck(w http.ResponseWriter, r *http.Request) {
	var token model.Token
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		logger.Warning(r, errors.InvalidRequestBody.SetOrigin(err).Error())
		errorResponse(w, errors.InvalidRequestBody)
		return
	}

	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if conf.ServerPasswd != token.ServerPasswd {
		logger.Error(r, errors.UserNotLogged.SetArgs("Пароль сервера не верен", "Server password missmatch"))
		errorResponse(w, errors.UserNotLogged)
		return
	}

	if Err := hash.CheckAccessTokenSignature(token.AccessToken); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	successResponse(w, nil)
	logger.Success(r, "token signature was checked")
}
