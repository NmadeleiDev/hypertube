package hash

import (
	"auth_backend/errors"
	// "crypto/aes"
	// "crypto/cipher"
	// "crypto/rand"
	// "encoding/base64"
	"hash/crc32"
	// "io"
	"strconv"
)

/*
**	Хэш пароля. Указателем так как в аккаунте он должен быть указателем
**	(меньше обрабатывать в http хэндлерах)
 */
func PasswdHash(pass string) (*string, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return nil, Err
	}
	pass += conf.PasswdSalt
	crcH := crc32.ChecksumIEEE([]byte(pass))
	passHash := strconv.FormatUint(uint64(crcH), 20)
	return &passHash, nil
}

/*
**	Кодирует email в строку. Используется для валидации почты юзера при регистрации
 */
// func EmailHashEncode(mail string) (string, *errors.Error) {
// 	conf, Err := getConfig()
// 	if Err != nil {
// 		return "", Err
// 	}
// 	c, err := aes.NewCipher([]byte(conf.MasterKey))
// 	if err != nil {
// 		return "", errors.HashError.SetOrigin(err)
// 	}
// 	gcm, err := cipher.NewGCM(c)
// 	if err != nil {
// 		return "", errors.HashError.SetOrigin(err)
// 	}
// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
// 		return "", errors.HashError.SetOrigin(err)
// 	}
// 	token := gcm.Seal(nonce, nonce, []byte(mail), nil)
// 	return base64.URLEncoding.EncodeToString(token), nil
// }

/*
**	Раскодирует строку в email пользователя. Используется для валидации
**	почты юзера при регистрации
 */
// func EmailHashDecode(code string) (string, *errors.Error) {
// 	conf, Err := getConfig()
// 	if Err != nil {
// 		return "", Err
// 	}
// 	encodedToken, err := base64.URLEncoding.DecodeString(code)
// 	if err != nil {
// 		return "", errors.ImpossibleToExecute.SetArgs("код подтверждения невалиден",
// 			"confirm code invalid").SetHidden("На этапе распаковки url кодирования").SetOrigin(err)
// 	}
// 	c, err := aes.NewCipher([]byte(conf.MasterKey))
// 	if err != nil {
// 		return "", errors.ImpossibleToExecute.SetArgs("код подтверждения невалиден",
// 			"confirm code invalid").SetHidden("На этапе расшифровки мастер ключем").SetOrigin(err)
// 	}
// 	gcm, err := cipher.NewGCM(c)
// 	if err != nil {
// 		return "", errors.ImpossibleToExecute.SetArgs("код подтверждения невалиден",
// 			"confirm code invalid").SetHidden("На этапе расшифровки").SetOrigin(err)
// 	}
// 	nonceSize := gcm.NonceSize()
// 	if len(encodedToken) < nonceSize {
// 		return "", errors.ImpossibleToExecute.SetArgs("код подтверждения невалиден",
// 			"confirm code invalid").SetHidden("ошибка размера при декодировании токена")
// 	}
// 	nonce, encodedToken := encodedToken[:nonceSize], encodedToken[nonceSize:]
// 	mail, err := gcm.Open(nil, nonce, encodedToken, nil)
// 	if err != nil {
// 		return "", errors.ImpossibleToExecute.SetArgs("код подтверждения невалиден",
// 			"confirm code invalid").SetHidden("При декодировании токена").SetOrigin(err)
// 	}
// 	return string(mail), nil
// }
