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
**	/api/profile/get
**	Возвращает личные данные пользователя в случае отсутствия параметра id
**	В случае наличия параметра id - возвращает данные пользователя с данным id
**	в этом случае приватные поля (email) будут скрыты
**	-- Проверено
 */
func profileGet(w http.ResponseWriter, r *http.Request) {
	/*
	**	Получаю токен из заголовка
	 */
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

	var id uint
	var idString = r.FormValue("userId")
	if idString != "" {
		idInt, err := strconv.Atoi(idString)
		if err != nil {
			logger.Warning(r, "Id пользователя ("+idString+") содежит ошибку: "+err.Error())
			errorResponse(w, errors.InvalidRequestBody)
			return
		}
		if idInt < 1 {
			logger.Warning(r, "Id пользователя ("+idString+") меньше единицы")
			errorResponse(w, errors.InvalidRequestBody)
			return
		}
		id = uint(idInt)
		logger.Log(r, "Пользователь #"+strconv.Itoa(int(header.UserId))+" ищет профиль пользователя #"+idString)
	} else {
		id = header.UserId
		logger.Log(r, "Пользователь #"+strconv.Itoa(int(header.UserId))+" ищет свой профиль")
	}

	user, Err := postgres.UserGetBasicById(id)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}

	/*
	**	Если пользователь не мой - почистить приватные поля
	 */
	if user.UserId != header.UserId {
		user.Sanitize()
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		logger.Error(r, errors.MarshalError.SetOrigin(err))
		errorResponse(w, errors.MarshalError)
		return
	}

	successResponse(w, jsonUser)
	logger.Success(r, "Profile of user #"+strconv.Itoa(int(user.UserId))+" provided")
}
