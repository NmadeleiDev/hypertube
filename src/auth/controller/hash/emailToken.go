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
**	Создание подписи к токену изменения почтового адреса
 */
func createEmailTokenSignature(headerBase64 string) (string, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return "", Err
	}
	headerBase64 += conf.PasswdSalt + "emailPatch"
	crcH := crc32.ChecksumIEEE([]byte(headerBase64))
	return strconv.FormatUint(uint64(crcH), 20), nil
}

/*
**	Проверка подписи токена восстановления пароля
 */
func checkEmailTokenPartsSignature(headerBase64, origSignatureBase64 string) *errors.Error {
	signatureBase64, Err := createEmailTokenSignature(headerBase64)
	if Err != nil {
		return Err
	}
	if signatureBase64 != origSignatureBase64 {
		return errors.InvalidToken.SetHidden("подпись содержит ошибку")
	}
	return nil
}

/*
**	Создание токена восстановления пароля
 */
func CreateEmailToken(user *model.UserBasic) (string, *errors.Error) {
	var header model.EmailTokenHeader
	header.UserId = user.UserId

	if user.NewEmail == nil {
		return "", errors.HashError.SetHidden("Поле NewEmail равно nil")
	}
	header.NewEmail = *user.NewEmail

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", errors.MarshalError.SetOrigin(err)
	}
	headerBase64 := base64.StdEncoding.EncodeToString(headerJson)
	signatureBase64, Err := createEmailTokenSignature(headerBase64)
	if Err != nil {
		return "", Err
	}
	return base64.StdEncoding.EncodeToString([]byte(headerBase64 + "." + signatureBase64)), nil
}

/*
**	Проверка правильности подписи токена восстановления пароля и распаковка его хидера
 */
func GetHeaderFromEmailToken(emailTokenBase64 string) (model.EmailTokenHeader, *errors.Error) {
	var header model.EmailTokenHeader

	decodedEmailToken, err := base64.StdEncoding.DecodeString(emailTokenBase64)
	if err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	tokenParts := strings.Split(string(decodedEmailToken), ".")
	if len(tokenParts) != 2 {
		return header, errors.InvalidToken.SetHidden("Токен должен состоять из 2 частей - но содержит " + strconv.Itoa(len(tokenParts)))
	}

	if Err := checkEmailTokenPartsSignature(tokenParts[0], tokenParts[1]); Err != nil {
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
