package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"auth_backend/model"
	"auth_backend/postgres"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type tokenVk struct {
	AccessToken string `json:"access_token"`
	UserId      uint   `json:"user_id"`
	Email       string `json:"email,omitempty"`
	ExpiresIn   uint   `json:"expires_in"`
	ExpiresAt   time.Time
}

type responseVk struct {
	UserVkId     uint   `json:"id"`
	Fname        string `json:"first_name"`
	Lname        string `json:"last_name"`
	Username     string `json:"screen_name"`
	IsImageExist uint   `json:"has_photo"`
	ImageBody    string `json:"photo_200"`
}

type profileVk struct {
	Response []responseVk `json:"response"`
}

/*
**	/api/auth/oauthVk
**	Авторизация oauth ресурса vk.com
**	Особенность - метод api vk.com которым я пользуюсь - не возвращает email пользователя
**	Поэтому его аккаунт всегда создается без указания почтового адреса
**	-- Проверено
 */
func authOauthVk(w http.ResponseWriter, r *http.Request) {
	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	params, Err := parseRequestParamsVk(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}
	/*
	**	Get vk api access token
	 */
	token, Err := getTokenFromVk(params)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	/*
	**	Getting user profile from vk api and fills it into *model.UserVk
	 */
	user, Err := getUserVk(token)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	var userBasic *model.UserBasic
	/*
	**	getting user from db if it exists
	 */
	userFromDb, Err := postgres.UserGetVkById(user.UserVkId)
	if Err != nil {
		if errors.UserNotExist.IsOverlapWithError(Err) {
			// user not exists
			logger.Log(r, "UserVk with userVkId "+strconv.Itoa(int(user.UserVkId))+" not found in database. Creating new one")
			if userBasic, Err = postgres.UserSetVk(user); Err != nil {
				logger.Error(r, Err)
				http.Redirect(w, r,
					conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
					http.StatusTemporaryRedirect)
				return
			}
		} else {
			// database error
			logger.Error(r, Err)
			http.Redirect(w, r,
				conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
				http.StatusTemporaryRedirect)
			return
		}
	} else {
		user.UserId = userFromDb.UserId
		if Err = postgres.UserUpdateVk(user); Err != nil {
			logger.Error(r, Err)
			http.Redirect(w, r,
				conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
				http.StatusTemporaryRedirect)
			return
		}
		userBasic, Err = postgres.UserGetBasicById(user.UserId)
		if Err != nil {
			logger.Error(r, Err)
			http.Redirect(w, r,
				conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
				http.StatusTemporaryRedirect)
			return
		}
		if Err = postgres.UserUpdateVk(user); Err != nil {
			logger.Error(r, Err)
			http.Redirect(w, r,
				conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
				http.StatusTemporaryRedirect)
			return
		}
	}

	accessToken, Err := hash.CreateAccessToken(userBasic)
	if Err != nil {
		logger.Warning(r, "cannot get password hash - "+Err.Error())
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	logger.Success(r, "user #"+strconv.Itoa(int(user.UserId))+" was authenticated")

	http.Redirect(w, r,
		conf.SocketRedirect+conf.OauthRedirect+"?accessToken="+accessToken,
		http.StatusTemporaryRedirect)
}

/*
**	Parsing GET params from request
 */
func parseRequestParamsVk(r *http.Request) (requestParams, *errors.Error) {
	var params requestParams

	params.Code = r.FormValue("code")
	params.State = r.FormValue("state")
	params.Error = r.FormValue("error")
	params.ErrorDescription = r.FormValue("error_description")

	if params.Error != "" || params.ErrorDescription != "" {
		return params, errors.AccessDenied.SetHidden("Сервер авторизации vk.com ответил: " +
			params.Error + " - " + params.ErrorDescription)
	}
	if params.Code == "" || params.State == "" {
		return params, errors.AccessDenied.SetHidden("Сервер авторизации vk.com прислал невалидные данные. code: " +
			params.Code + " state" + params.State)
	}
	return params, nil
}

/*
**	Request to vk API for token
 */
func getTokenFromVk(params requestParams) (tokenVk, *errors.Error) {
	var result tokenVk

	conf, Err := getConfig()
	if Err != nil {
		return result, Err
	}
	portString := strconv.FormatUint(uint64(conf.ServerPort), 10)

	formData := url.Values{
		"client_id":     {conf.VkClientId},
		"client_secret": {conf.VkSecret},
		"code":          {params.Code},
		"redirect_uri":  {"http://" + conf.ServerIp + ":" + portString + "/api/auth/oauthVk"},
	}
	resp, err := http.PostForm("https://oauth.vk.com/access_token", formData)
	if err != nil {
		return result, errors.AccessDenied.SetHidden("Запрос токена из vk.com провален").SetOrigin(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, errors.AccessDenied.SetHidden("Декодирование json дало ошибку").SetOrigin(err)
	}
	duration, err := time.ParseDuration(strconv.FormatUint(uint64(result.ExpiresIn), 10) + "s")
	if err != nil {
		return result, errors.UnknownInternalError.SetArgs("ошибка парсинга времени", "time parse fail").SetOrigin(err)
	}
	result.ExpiresAt = time.Now().Add(duration)
	return result, nil
}

/*
**	Request to vk API for user profile
 */
func getUserProfileVk(token tokenVk) (profileVk, *errors.Error) {
	var profile profileVk
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	url := "https://api.vk.com/method/users.get?access_token=" + token.AccessToken +
		"&user_ids=" + strconv.Itoa(int(token.UserId)) +
		"&fields=has_photo,photo_200,nickname,screen_name" +
		"&v=5.130"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Запрос данных пользователя vk.com провален").SetOrigin(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Запрос данных пользователя vk.com провален").SetOrigin(err)
	}
	defer resp.Body.Close() // важный пункт!
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Чтение данных пользователя vk.com провалено").SetOrigin(err)
	}
	if err = json.Unmarshal(respBody, &profile); err != nil {
		return profile, errors.AccessDenied.SetHidden("Декодирование данных пользователя из json дало ошибку").SetOrigin(err)
	}
	if len(profile.Response) == 0 {
		return profile, errors.EmptyResponse
	}
	return profile, nil
}

/*
**	Forming UserVk structure
 */
func getUserVk(token tokenVk) (*model.UserVk, *errors.Error) {
	profile, Err := getUserProfileVk(token)
	if Err != nil {
		return nil, Err
	}

	var imageBodyPtr *string

	if profile.Response[0].IsImageExist == 1 {
		imageBodyPtr = &profile.Response[0].ImageBody
	} else {
		imageBodyPtr = nil
	}

	return &model.UserVk{
		Fname:     profile.Response[0].Fname,
		Lname:     profile.Response[0].Lname,
		Username:  profile.Response[0].Username,
		ImageBody: imageBodyPtr,
		UserVkModel: model.UserVkModel{
			UserVkId:    profile.Response[0].UserVkId,
			AccessToken: &token.AccessToken,
			ExpiresAt:   &token.ExpiresAt,
		},
	}, nil
}
