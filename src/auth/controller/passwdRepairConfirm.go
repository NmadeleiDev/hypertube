package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/logger"
	"net/http"
	"net/url"
	"strconv"
)

/*
**	/api/passwd/repair/confirm
**	Восстановление утраченого пароля. Второй эндпоинт из трех.
**	Должен активизироваться из письма на почте. Отвечает редиректом на страницу предлагающую установить новый пароль.
**	-- Не проверено
 */
func passwdRepairConfirm(w http.ResponseWriter, r *http.Request) {
	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())), //base64.StdEncoding.EncodeToString(Err.ToJson())
			http.StatusTemporaryRedirect)
		return
	}

	repairToken, Err := parseCodeFromRequest(r)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	// Делаю это как для того чтобы узнать ID пользователя, так и для проверки валидности токена
	tokenHeader, Err := hash.GetHeaderFromRepairToken(repairToken)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(tokenHeader.UserId))+logger.NO_COLOR+
		" успешно обратился из письма для восстановления пароля. Переадресую его на страницу заполнения нового пароля")
	http.Redirect(w, r,
		conf.SocketRedirect+conf.PasswdResetRedirect+"?repairToken="+url.QueryEscape(repairToken),
		http.StatusTemporaryRedirect)
	println(tokenHeader.UserId)
}
