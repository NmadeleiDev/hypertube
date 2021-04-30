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
**	/api/email/resend
**	Повторная отправка кода подтверждения почты на почту
**	Пользователь не залогинен (так как для этого нужно сначала подтвердить почту).
**	Письмо отсылается только в случае если пользователь еще не подтвердил почту.
**	-- Проверено
 */
func emailResend(w http.ResponseWriter, r *http.Request) {
	email, Err := parseEmailFromRequest(r)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user, Err := postgres.UserGetBasicByEmail(email)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if user.IsEmailConfirmed == true {
		logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
			" уже подтвердил свою почту")
		successResponse(w, nil)
		return
	}

	emailToken, Err := hash.CreateEmailToken(user)
	if Err != nil {
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

	go func(user *model.UserBasic, emailToken, serverIp string, serverPort uint) {
		if Err := mailer.SendEmailConfirmMessage(user, emailToken, serverIp, serverPort); Err != nil {
			logger.Error(r, Err)
		} else {
			logger.Success(r, "Писмьмо для подтверждения почты пользователя #"+
				logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+" успешно отправлено")
		}
	}(user, emailToken, conf.ServerIp, conf.ServerPort)

	logger.Success(r, "Повторное письмо пользователя #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" обработано и поставлено в очередь на отправку")
	successResponse(w, nil)
}

func parseEmailFromRequest(r *http.Request) (string, *errors.Error) {
	type Request struct {
		Email *string `json:"email"`
	}
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return "", errors.InvalidRequestBody.SetOrigin(err)
	}
	if request.Email == nil || *request.Email == "" {
		return "", errors.NoArgument.SetArgs("не указана почта", "email expected")
	}
	return *request.Email, nil
}
