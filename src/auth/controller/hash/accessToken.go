package hash

import (
	"auth_backend/errors"
	"auth_backend/model"
	"encoding/base64"
	"encoding/json"
	"hash/crc32"
	"strconv"
	"strings"
)

/*
**	Создание подписи к токену доступа
 */
func createAccessTokenSignature(headerBase64 string) (string, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return "", Err
	}
	headerBase64 += conf.PasswdSalt + "accessToken"
	crcH := crc32.ChecksumIEEE([]byte(headerBase64))
	return strconv.FormatUint(uint64(crcH), 20), nil
}

/*
**	Проверка подписи токена доступа (когда он уже частично раскодирован)
 */
func checkAccessTokenPartsSignature(headerBase64, origSignatureBase64 string) *errors.Error {
	signatureBase64, Err := createAccessTokenSignature(headerBase64)
	if Err != nil {
		return Err
	}
	if signatureBase64 != origSignatureBase64 {
		return errors.InvalidToken.SetHidden("подпись содержит ошибку")
	}
	return nil
}

/*
**	Создание токена доступа
 */
func CreateAccessToken(user *model.UserBasic) (string, *errors.Error) {
	var header model.AccessTokenHeader
	header.UserId = user.UserId

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", errors.MarshalError.SetOrigin(err)
	}
	headerBase64 := base64.StdEncoding.EncodeToString(headerJson)
	signature, Err := createAccessTokenSignature(headerBase64)
	if Err != nil {
		return "", Err
	}
	return base64.StdEncoding.EncodeToString([]byte(headerBase64 + "." + signature)), nil
}

/*
**	Проверка подписи токена доступа (использутся при проверке авторизации пользователя)
 */
func CheckAccessTokenSignature(accessTokenBase64 string) *errors.Error {
	decodedAccessToken, err := base64.StdEncoding.DecodeString(accessTokenBase64)
	if err != nil {
		return errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	tokenParts := strings.Split(string(decodedAccessToken), ".")
	if len(tokenParts) != 2 {
		return errors.InvalidToken.SetHidden("Токен должен состоять из 2 частей - но содержит " + strconv.Itoa(len(tokenParts)))
	}
	signature, Err := createAccessTokenSignature(tokenParts[0])
	if Err != nil {
		return Err
	}
	if signature != tokenParts[1] {
		return errors.InvalidToken.SetHidden("подпись содержит ошибку")
	}
	return nil
}

/*
**	Проверка правильности подписи токена доступа и распаковка его хидера
 */
func GetHeaderFromAccessToken(accessTokenBase64 string) (model.AccessTokenHeader, *errors.Error) {
	var header model.AccessTokenHeader

	decodedAccessToken, err := base64.StdEncoding.DecodeString(accessTokenBase64)
	if err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	tokenParts := strings.Split(string(decodedAccessToken), ".")
	if len(tokenParts) != 2 {
		return header, errors.InvalidToken.SetHidden("Токен должен состоять из 2 частей - но содержит " + strconv.Itoa(len(tokenParts)))
	}

	if Err := checkAccessTokenPartsSignature(tokenParts[0], tokenParts[1]); Err != nil {
		return header, Err
	}

	decodedHeader, err := base64.StdEncoding.DecodeString(tokenParts[0])
	if err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	if err = json.Unmarshal(decodedHeader, &header); err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования json").SetOrigin(err)
	}
	return header, nil
}
