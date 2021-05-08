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

type tokenFb struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
	ExpiresAt   time.Time
}

type profileFbId struct {
	IdString string `json:"id"`
	IdUint   uint   `json:"-"`
}

type avatarFb struct {
	Height       uint   `json:"height"`
	Width        uint   `json:"width"`
	Url          string `json:"url"`
	IsSilhouette bool   `json:"is_silhouette"`
}

type avatarDataFb struct {
	Data avatarFb `json:"data"`
}

type profileFb struct {
	Email     *string       `json:"email"`
	Fname     string        `json:"first_name"`
	Lname     string        `json:"last_name"`
	Username  string        `json:"short_name"`
	ImageBody *avatarDataFb `json:"picture"`
}

/*
**	/api/auth/oauthFb
**	Авторизация oauth facebook.com
**	Особенности: несколько запросов. Ко всем запросам надо параметром задавать id пользователя
**	поэтому первым запросом выясняю его id. Далее получаю профиль пользователя. Но фотография в
**	профиле имеет расширение 50*50 поэтому делаю еще один запрос для получения фотографии в
**	более высоком расрешении. Фотография при этом имеет бинарный вид. Это надо протестить.
**	-- Нужно протестировать фотографию пользователя
 */
func authOauthFb(w http.ResponseWriter, r *http.Request) {
	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	params, Err := parseRequestParamsFb(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	/*
	**	Get facebook api access token
	 */
	token, Err := getTokenFromFb(params)
	if Err != nil {
		logger.Error(r, Err)
		http.Redirect(w, r,
			conf.SocketRedirect+conf.ErrorRedirect+"?error="+url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	/*
	**	Get user profile from facebook api and fills it into *model.UserFb
	 */
	user, Err := getUserFb(token)
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
	userFromDb, Err := postgres.UserGetFbById(user.UserFbId)
	if Err != nil {
		if errors.UserNotExist.IsOverlapWithError(Err) {
			// user not exists
			logger.Log(r, "UserFb with userFbId "+strconv.Itoa(int(user.UserFbId))+" not found in database. Creating new one")
			if userBasic, Err = postgres.UserSetFb(user); Err != nil {
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
		// user exists
		user.UserId = userFromDb.UserId
		if Err = postgres.UserUpdateFb(user); Err != nil {
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
		if Err = postgres.UserUpdateFb(user); Err != nil {
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
func parseRequestParamsFb(r *http.Request) (requestParams, *errors.Error) {
	var params requestParams

	params.Code = r.FormValue("code")
	params.Error = r.FormValue("error")
	params.ErrorDescription = r.FormValue("error_description")

	if params.Error != "" || params.ErrorDescription != "" {
		return params, errors.AccessDenied.SetHidden("Сервер авторизации facebook ответил: " +
			params.Error + " - " + params.ErrorDescription)
	}
	if params.Code == "" {
		return params, errors.AccessDenied.SetHidden("Сервер авторизации facebook прислал невалидные данные. code: " + params.Code)
	}
	return params, nil
}

/*
**	Request to facebook server API for token
 */
func getTokenFromFb(params requestParams) (tokenFb, *errors.Error) {
	var result tokenFb

	conf, Err := getConfig()
	if Err != nil {
		return result, Err
	}
	portString := strconv.FormatUint(uint64(conf.ServerPort), 10)

	url := "https://graph.facebook.com/v10.0/oauth/access_token" +
		"?client_id=" + conf.FacebookClientId +
		"&redirect_uri=http://" + conf.ServerIp + ":" + portString + "/api/auth/oauthFb" +
		"&client_secret=" + conf.FacebookSecret +
		"&code=" + params.Code

	resp, err := http.Get(url)
	if err != nil {
		return result, errors.AccessDenied.SetHidden("Запрос токена из facebook провален").SetOrigin(err)
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
**	Request to facebook server API for user ID
 */
func getUserIdFb(accessToken string) (profileFbId, *errors.Error) {
	var profileId profileFbId
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
	url := "https://graph.facebook.com/me?access_token=" + accessToken
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return profileId, errors.AccessDenied.SetHidden("Запрос данных пользователя facebook провален").SetOrigin(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return profileId, errors.AccessDenied.SetHidden("Запрос данных пользователя facebook провален").SetOrigin(err)
	}
	defer resp.Body.Close() // важный пункт!
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return profileId, errors.AccessDenied.SetHidden("Чтение данных пользователя facebook провалено").SetOrigin(err)
	}
	if err = json.Unmarshal(respBody, &profileId); err != nil {
		return profileId, errors.AccessDenied.SetHidden("Декодирование данных пользователя из json дало ошибку").SetOrigin(err)
	}
	idInt, err := strconv.Atoi(profileId.IdString)
	if err != nil {
		return profileId, errors.AccessDenied.SetHidden("Декодирование id пользователя из строки дало ошибку").SetOrigin(err)
	}
	profileId.IdUint = uint(idInt)
	return profileId, nil
}

/*
**	Request to facebook server API for user profile
 */
func getUserProfileFb(accessToken, userId string) (profileFb, *errors.Error) {
	var profile profileFb
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
	url := "https://graph.facebook.com/" + userId +
		"?access_token=" + accessToken +
		"&fields=email,first_name,last_name,picture,short_name"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Запрос данных пользователя facebook провален").SetOrigin(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Запрос данных пользователя facebook провален").SetOrigin(err)
	}
	defer resp.Body.Close() // важный пункт!
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Чтение данных пользователя facebook провалено").SetOrigin(err)
	}
	if err = json.Unmarshal(respBody, &profile); err != nil {
		return profile, errors.AccessDenied.SetHidden("Декодирование данных пользователя из json дало ошибку").SetOrigin(err)
	}
	return profile, nil
}

/*
**	Request to facebook server API for user photo in good resolution
 */
func getUserPhotoFb(accessToken, userId string) (*string, *errors.Error) {
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
	url := "https://graph.facebook.com/" + userId + "/picture" +
		"?access_token=" + accessToken +
		// small, normal, album, large, square
		"&type=normal"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.AccessDenied.SetHidden("Запрос фотографии пользователя facebook провален").SetOrigin(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.AccessDenied.SetHidden("Запрос фотографии пользователя facebook провален").SetOrigin(err)
	}
	defer resp.Body.Close() // важный пункт!
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.AccessDenied.SetHidden("Чтение фотографии пользователя facebook провалено").SetOrigin(err)
	}
	imageBody := string(respBody)
	print("image body ")
	println(imageBody)
	return &imageBody, nil
}

/*
**	Forming UserFb structure
 */
func getUserFb(token tokenFb) (*model.UserFb, *errors.Error) {
	profileId, Err := getUserIdFb(token.AccessToken)
	if Err != nil {
		return nil, Err
	}

	profile, Err := getUserProfileFb(token.AccessToken, profileId.IdString)
	if Err != nil {
		return nil, Err
	}

	// imageBodyPtr, Err := getUserPhotoFb(token.AccessToken, profileId.IdString)
	// if Err != nil {
	// 	return nil, Err
	// }

	return &model.UserFb{
		Email:     profile.Email,
		Fname:     profile.Fname,
		Lname:     profile.Lname,
		Username:  profile.Username,
		// ImageBody: imageBodyPtr,
		ImageBody: &profile.ImageBody.Data.Url,
		UserFbModel: model.UserFbModel{
			UserFbId:    profileId.IdUint,
			AccessToken: &token.AccessToken,
			ExpiresAt:   &token.ExpiresAt,
		},
	}, nil
}
