package controller

import (
	"auth_backend/controller/hash"
	"auth_backend/errors"
	"auth_backend/logger"
	"net/http"
	"time"
)

/*
**	Все middleware функции для эндпоинтов описаны в этом файле
 */
func panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if ok {
					logger.Error(r, errors.UnknownInternalError.SetArgs("Произошла ПАНИКА", "PANIC happened").SetOrigin(err))
				} else {
					logger.Error(r, errors.UnknownInternalError.SetArgs("Произошла ПАНИКА, отсутствует интерфейс ошибки",
						"PANIC happened, error interface expected"))
				}
				errorResponse(w, errors.UnknownInternalError)
				return
			}
		}()
		next.ServeHTTP(w, r)
		logger.Duration(r, time.Since(t))
	})
}

func authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("accessToken")
		if accessToken == "" {
			logger.Error(r, errors.UserNotLogged.SetArgs("отсутствует токен доступа", "access token expected"))
			errorResponse(w, errors.UserNotLogged)
			return
		}
		if Err := hash.CheckAccessTokenSignature(accessToken); Err != nil {
			logger.Error(r, Err)
			errorResponse(w, errors.UserNotLogged)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsPut(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "PUT,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,accessToken")

		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "PUT" {
			logger.Warning(r, "wrong request method. Should be PUT method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsDelete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "DELETE,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,accessToken")

		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "DELETE" {
			logger.Warning(r, "wrong request method. Should be DELETE method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,accessToken,repairToken")

		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "POST" {
			logger.Warning(r, "wrong request method. Should be POST method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsPatch(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "PATCH,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,accessToken,repairToken")

		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "PATCH" {
			logger.Warning(r, "wrong request method. Should be PATCH method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,Authorization,accessToken,repairToken")
		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "GET" {
			logger.Warning(r, "wrong request method. Should be GET method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsGetPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,Authorization,accessToken")
		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "GET" && r.Method != "POST" {
			logger.Warning(r, "wrong request method. Should be GET or POST method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}
