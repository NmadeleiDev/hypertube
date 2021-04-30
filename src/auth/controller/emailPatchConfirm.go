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
**	api/email/patch/confirm
**	Запускается письмом с почты, которое приходит после успешной
**	активации эндпоинта /api/email/patch
**	В конце запроса должен быть редирект на успешный или не успешный эндпоинт
**	-- Проверено
 */
func emailPatchConfirm(w http.ResponseWriter, r *http.Request) {
	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	emailPatchToken, Err := parseCodeFromRequest(r)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	tokenHeader, Err := hash.GetHeaderFromEmailToken(emailPatchToken)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	user, Err := postgres.UserGetBasicById(tokenHeader.UserId)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	if user.NewEmail == nil {
		logger.Warning(r, "Пользователь уже обновил статус свого почтового адреса")
		Err := errors.ImpossibleToExecute.SetArgs("статус почтового адреса уже обновлен", "email address status already changed")
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}
	if *user.NewEmail != tokenHeader.NewEmail {
		logger.Warning(r, "Пользователь уже обновил статус свого почтового адреса. Это старый код подтверждения")
		Err := errors.ImpossibleToExecute.SetArgs("код подтверждения просрочен", "confirm code is expired")
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	user.Email = user.NewEmail

	if Err = postgres.UserConfirmEmailBasic(user); Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	logger.Success(r, "Пользователь #"+logger.BLUE+strconv.Itoa(int(user.UserId))+logger.NO_COLOR+
		" подтвердил изменения своего почтового адреса "+logger.BLUE+*user.Email+logger.NO_COLOR)
	http.Redirect(w, r,
		conf.SocketRedirect+conf.ErrorRedirect,
		http.StatusTemporaryRedirect)
}
