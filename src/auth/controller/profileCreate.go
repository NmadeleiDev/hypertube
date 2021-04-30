package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/controller/mailer"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
**	/api/profile/create
**	Базовая регистрация пользователя
**	-- Проверено
 */
func profileCreate(w http.ResponseWriter, r *http.Request) {
	user, Err := parseUserBasicFromRequest(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		errorResponse(w, Err)
		return
	}

	if Err = user.Validate(); Err != nil {
		logger.Warning(r, Err.Error())
		errorResponse(w, Err)
		return
	}

	if user.EncryptedPass, Err = hash.PasswdHash(user.Passwd); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user.NewEmail = user.Email

	if Err = postgres.UserSetBasic(user); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	emailToken, Err := hash.CreateEmailToken(user)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	go func(user *model.UserBasic, emailToken, serverIp string, serverPort uint) {
		if Err := mailer.SendEmailConfirmMessage(user, emailToken, serverIp, serverPort); Err != nil {
			logger.Error(r, Err)
		} else {
			logger.Success(r, "Писмьмо для подтверждения почты пользователя #"+
				logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+" успешно отправлено")
		}
	}(user, emailToken, conf.ServerIp, conf.ServerPort)

	userJson, err := json.Marshal(user)
	if err != nil {
		logger.Error(r, errors.MarshalError.SetOrigin(Err))
		errorResponse(w, errors.MarshalError)
		return
	}

	successResponse(w, userJson)
	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+" успешно создан")
}

func parseUserBasicFromRequest(r *http.Request) (*model.UserBasic, *errors.Error) {
	type userWithPassword struct {
		model.UserBasicModel
		Passwd string `json:"passwd"`
	}
	var requestUser = userWithPassword{}
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		return nil, errors.InvalidRequestBody.SetOrigin(err)
	}
	var user = &model.UserBasic{UserBasicModel: requestUser.UserBasicModel}
	user.Passwd = requestUser.Passwd
	return user, nil
}
