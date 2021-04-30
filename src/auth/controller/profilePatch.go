package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/controller/validator"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"net/http"
	"strconv"
)

/*
**	/api/profile/patch
**	Обновление полей пользователя username, first_name, last_name, image_body
**	авторизация в авторизационном хидере accessToken
**	-- Проверено
 */
func profilePatch(w http.ResponseWriter, r *http.Request) {
	user, Err := parseUserBasicFromRequest(r)
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

	user.UserId = header.UserId

	if user.ImageBody == nil && user.Username == "" && user.Fname == nil && user.Lname == nil {
		logger.Warning(r, "Нет ни одного поля для модификации")
		errorResponse(w, errors.NoArgument)
		return
	}

	if message, Err := validateProfilePatchFields(user); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	} else {
		logger.Log(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+message)
	}

	if Err = postgres.UserUpdateBasic(user); Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	successResponse(w, nil)
	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" успешно обновил свои поля")
}

/*
**	Проверка и логгирование полей из запроса фронта
 */
func validateProfilePatchFields(user *model.UserBasic) (string, *errors.Error) {
	var fieldsMessage string

	if user.ImageBody != nil {
		fieldsMessage += "ImageBody=hidden "
	}
	if user.Username != "" {
		if Err := validator.ValidateName(user.Username); Err != nil {
			return "", Err
		}
		fieldsMessage += "Username=" + logger.BLUE + user.Username + logger.NO_COLOR + " "
	}
	if user.Fname != nil {
		if Err := validator.ValidateName(*user.Fname); Err != nil {
			return "", Err
		}
		fieldsMessage += "Fname=" + logger.BLUE + *user.Fname + logger.NO_COLOR + " "
	}
	if user.Lname != nil {
		if Err := validator.ValidateName(*user.Lname); Err != nil {
			return "", Err
		}
		fieldsMessage += "Lname=" + logger.BLUE + *user.Lname + logger.NO_COLOR + " "
	}
	if fieldsMessage == "" {
		return "", errors.NoArgument.SetArgs("нет ни одного поля для модификации", "no fields for modification")
	}
	return ". Поля для изменения: " + fieldsMessage, nil
}
