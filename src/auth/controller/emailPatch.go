package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/controller/mailer"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"net/http"
	"strconv"
)

/*
**	/api/email/patch
**	изменение почты. Отправляет письмо для подтверждения на новую почту
**	Первый эндпоинт из двух. Первый дергает пользователь с сайта,
**	второй - с почты для подтверждения
**	-- Проверено
 */
func emailPatch(w http.ResponseWriter, r *http.Request) {
	newEmail, Err := parseEmailFromRequest(r)
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

	accessHeader, Err := hash.GetHeaderFromAccessToken(accessToken)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user, Err := postgres.UserGetBasicById(accessHeader.UserId)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user.NewEmail = &newEmail

	emailPatchToken, Err := hash.CreateEmailToken(user)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if Err = postgres.UserUpdateBasic(user); Err != nil {
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

	go func(user *model.UserBasic, emailPatchToken, serverIp string, serverPort uint) {
		if Err := mailer.SendEmailPatchMailAddress(user, emailPatchToken, serverIp, serverPort); Err != nil {
			logger.Error(r, Err)
		} else {
			logger.Success(r, "Писмьмо для подтверждения новой почты пользователя #"+
				logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+" успешно отправлено")
		}
	}(user, emailPatchToken, conf.ServerIp, conf.ServerPort)

	successResponse(w, nil)
	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" подал успешную заявку на изменение почтового адреса "+logger.BLUE+newEmail+logger.NO_COLOR)
}
