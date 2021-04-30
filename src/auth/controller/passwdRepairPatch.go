package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/controller/validator"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/postgres"
	"net/http"
	"strconv"
)

/*
**	/api/passwd/repair/patch
**	восстановление пароля в случае его утраты. Последний эндпоинт (из трех)
**	В теле запроса ожидается поле пароля (passwd)
**	Ожидается заголовок repairToken
**	-- Не протестировано
 */
func passwdRepairPatch(w http.ResponseWriter, r *http.Request) {
	repairToken := r.Header.Get("repairToken")
	if repairToken == "" {
		logger.Error(r, errors.UserNotLogged.SetArgs("отсутствует токен восстановления пароля", "password repair token expected"))
		errorResponse(w, errors.UserNotLogged)
		return
	}

	repairHeader, Err := hash.GetHeaderFromRepairToken(repairToken)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	passwd, Err := parsePasswordFromRequest(r)
	if Err != nil {
		logger.Error(r, Err) // Warning?
		errorResponse(w, Err)
		return
	}

	if Err = validator.ValidatePassword(passwd); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	user, Err := postgres.UserGetBasicById(repairHeader.UserId)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	if user.EncryptedPass, Err = hash.PasswdHash(passwd); Err != nil {
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
		" успешно восстановил пароль задав новый")
}
