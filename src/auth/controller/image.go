package controller

import (
	"auth_backend/logger"
	"io/ioutil"
	"net/http"
)

var (
	imageEmailPatch   []byte
	imagePasswdRepair []byte
)

/*
**	/api/image
**	на GET запрос отвечает содержимым картинки
**	Используется только для /api/info
**	-- Проверено
 */
func image(w http.ResponseWriter, r *http.Request) {
	var err error

	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	imageName := r.FormValue("image")

	switch imageName {
	case "emailPatch":
		if imageEmailPatch == nil {
			logger.Log(r, "Изображение еще не считано. Читаю его с диска")
			imageEmailPatch, err = ioutil.ReadFile(conf.ProjectRoot + "/images/emailPatch.jpg")
			if err != nil {
				logger.Warning(r, "Не смог считать изображение с диска - "+err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		successResponse(w, imageEmailPatch)
	case "passwdRepair":
		if imagePasswdRepair == nil {
			logger.Log(r, "Изображение еще не считано. Читаю его с диска")
			imagePasswdRepair, err = ioutil.ReadFile(conf.ProjectRoot + "/images/passwdRepair.jpg")
			if err != nil {
				logger.Warning(r, "Не смог считать изображение с диска - "+err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		successResponse(w, imagePasswdRepair)
	default:
		logger.Warning(r, "Изображение неизвестно: `"+imageName+"`")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	logger.Success(r, "Успешно отгрузил изображение "+imageName)
}
