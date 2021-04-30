package logger

import (
	"auth_backend/errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	RED       = "\033[31m"
	GREEN     = "\033[32m"
	YELLOW    = "\033[33m"
	BLUE      = "\033[34m"
	RED_BG    = "\033[41;30m"
	GREEN_BG  = "\033[42;30m"
	YELLOW_BG = "\033[43;30m"
	BLUE_BG   = "\033[44;30m"
	NO_COLOR  = "\033[m"
)

var logger *Logger

type Logger struct {
	IsEnabled bool `conf:"isEnabled"`
}

func GetConfig() *Logger {
	if logger == nil {
		logger = &Logger{}
	}
	return logger
}

func Log(r *http.Request, message string) {
	if logger == nil {
		log.Printf("%sLogger is not configured. Run GetConfig() function first%s\n", RED_BG, NO_COLOR)
		return
	}
	if logger.IsEnabled == false {
		return
	}
	log.Printf("%s %7s %20s %s\n", r.RemoteAddr, r.Method, r.URL.Path, message)
}

func PrintQuery(query string) {
	if logger == nil {
		log.Printf("%sLogger is not configured. Run GetConfig() function first%s\n", RED_BG, NO_COLOR)
		return
	}
	if logger.IsEnabled == false {
		return
	}
	log.Printf("Query: %s\n", query)
}

func Success(r *http.Request, message string) {
	if logger == nil {
		log.Printf("%sLogger is not configured. Run GetConfig() function first%s\n", RED_BG, NO_COLOR)
		return
	}
	if logger.IsEnabled == false {
		return
	}
	log.Printf("%s %7s %20s %s\n", r.RemoteAddr, r.Method, r.URL.Path, GREEN_BG+"SUCCESS: "+NO_COLOR+message)
}

func Warning(r *http.Request, message string) {
	if logger == nil {
		log.Printf("%sLogger is not configured. Run GetConfig() function first%s\n", RED_BG, NO_COLOR)
		return
	}
	if logger.IsEnabled == false {
		return
	}
	log.Printf("%s %7s %20s %s\n", r.RemoteAddr, r.Method, r.URL.Path,
		YELLOW_BG+"WARNING: "+NO_COLOR+message)
}

func error(r *http.Request, message string) {
	if logger == nil {
		log.Printf("%sLogger is not configured. Run GetConfig() function first%s\n", RED_BG, NO_COLOR)
		return
	}
	if logger.IsEnabled == false {
		return
	}
	log.Printf("%s %7s %20s %s\n", r.RemoteAddr, r.Method, r.URL.Path, RED_BG+"ERROR: "+NO_COLOR+message)
}

func Error(r *http.Request, Err *errors.Error) {
	if Err.IsInternalOrExternal() {
		error(r, Err.Error())
	} else {
		Warning(r, Err.Error())
	}
}

func Duration(r *http.Request, dur time.Duration) {
	if logger == nil {
		log.Printf("%sLogger is not configured. Run GetConfig() function first%s\n", RED_BG, NO_COLOR)
		return
	}
	if logger.IsEnabled == false {
		return
	}
	milliseconds := (int)(dur.Milliseconds())
	color := GREEN_BG
	if milliseconds > 2 {
		color = YELLOW_BG
	}
	if milliseconds > 10 {
		color = RED_BG
	}
	log.Printf("%s %7s %20s %s\n", r.RemoteAddr, r.Method, r.URL.Path,
		"time : "+color+strconv.Itoa(milliseconds)+NO_COLOR)
}
