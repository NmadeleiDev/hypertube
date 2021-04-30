package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/postgres"
	"net/http"
	"net/url"
	"strconv"
)

/*
**	/api/email/confirm
**	Подтверждение почтового адреса
**	-- Проверено
 */
func emailConfirm(w http.ResponseWriter, r *http.Request) {
	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())), //base64.StdEncoding.EncodeToString(Err.ToJson())
			http.StatusTemporaryRedirect)
		return
	}

	emailToken, Err := parseCodeFromRequest(r)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	tokenHeader, Err := hash.GetHeaderFromEmailToken(emailToken)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	user, Err := postgres.UserGetBasicByEmail(tokenHeader.NewEmail)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	if user.IsEmailConfirmed == true {
		logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
			" уже подтвердил свою почту")
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect,
			http.StatusTemporaryRedirect)
		return
	}

	if Err = postgres.UserConfirmEmailBasic(user); Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" подтвердил свой почтовый адрес "+logger.BLUE+*user.Email+logger.NO_COLOR)
	http.Redirect(w, r,
		conf.SocketRedirect+conf.ErrorRedirect, // ErrorRedirect -  это страница авторизации
		http.StatusTemporaryRedirect)
}

func parseCodeFromRequest(r *http.Request) (string, *errors.Error) {
	var confirmCode = r.FormValue("code")
	if confirmCode == "" {
		return "", errors.NoArgument.SetArgs("Отсутствует код подтверждения", "confirm code expected")
	}
	return confirmCode, nil
}
