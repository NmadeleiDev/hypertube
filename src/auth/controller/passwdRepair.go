package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/controller/mailer"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"net/http"
	"strconv"
)

/*
**	/api/passwd/repair
**	восстановление пароля в случае его утраты. Первый эндпоинт (из трех)
**	остальные - /api/passwd/repair/confirm /api/passwd/repair/patch
**	-- Проверено
 */
func passwdRepair(w http.ResponseWriter, r *http.Request) {
	email, Err := parseEmailFromRequest(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		errorResponse(w, Err)
		return
	}

	user, Err := postgres.UserGetBasicByEmail(email)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	repairToken, Err := hash.CreateRepairToken(user)
	if Err != nil {
		logger.Warning(r, "cannot get password hash - "+Err.Error())
		errorResponse(w, Err)
		return
	}

	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	go func(user *model.UserBasic, repairToken, serverIp string, serverPort uint) {
		if Err := mailer.SendEmailPasswdRepair(user, repairToken, serverIp, serverPort); Err != nil {
			logger.Error(r, Err)
		} else {
			logger.Success(r, "Писмьмо для восстановления пароля пользователя #"+
				logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+" успешно отправлено")
		}
	}(user, repairToken, conf.ServerIp, conf.ServerPort)

	successResponse(w, nil)
	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" подал успешную заявку на восстановление пароля")
}
