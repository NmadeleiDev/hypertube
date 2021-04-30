package errors

import (
	"encoding/json"
	"net/http"
)

type errorType struct {
	errType    string
	httpStatus uint
}

var (
	ArgumentType     = errorType{errType: "Ошибка аргумента", httpStatus: http.StatusUnprocessableEntity}          // 422
	RequestType      = errorType{errType: "Невалидный запрос", httpStatus: http.StatusBadRequest}                  // 400
	InternalType     = errorType{errType: "Внутренняя ошибка сервера", httpStatus: http.StatusInternalServerError} // 500
	ExternalType     = errorType{errType: "Ошибка внешнего сервиса", httpStatus: http.StatusInternalServerError}   // 500
	BusinessType     = errorType{errType: "Ошибка бизнес логики", httpStatus: http.StatusNotAcceptable}            // 406
	UnauthorizedType = errorType{errType: "Пользователь не авторизован", httpStatus: http.StatusUnauthorized}      // 401
)

type Error struct {
	code           uint
	errorType      errorType
	ruDescription  string
	engDescription string
	ruArgument     string
	engArgument    string
	hidden         string
	original       error
}

type jsonExportedError struct {
	Code uint   `json:"code"`
	Ru   string `json:"ru"`
	Eng  string `json:"en"`
}

func new(code uint, errorType errorType, ruDescription string, engDescription string) *Error {
	var Err Error

	Err.code = code
	Err.errorType = errorType
	Err.ruDescription = ruDescription
	Err.engDescription = engDescription

	return &Err
}

func (Err Error) GetHttpStatus() int {
	return int(Err.errorType.httpStatus)
}

func (Err *Error) SetArgs(ruArgument, engArgument string) *Error {
	Err.ruArgument = ruArgument
	Err.engArgument = engArgument
	return Err
}

func (Err *Error) SetHidden(hiddenArgument string) *Error {
	Err.hidden = hiddenArgument
	return Err
}

func (Err *Error) SetOrigin(err error) *Error {
	Err.original = err
	return Err
}

func (Err *Error) ToJson() []byte {
	var exported = jsonExportedError{Code: Err.code}
	if Err.errorType == InternalType || Err.errorType == ExternalType {
		exported.Eng = "internal server error"
		exported.Ru = "Внутренняя ошибка сервера"
	} else {
		exported.Eng = Err.engDescription + " - " + Err.engArgument
		exported.Ru = Err.ruDescription + " - " + Err.ruArgument
	}
	dst, _ := json.Marshal(exported)
	return dst
}

func (Err Error) Error() string {
	if Err.original == nil {
		return Err.ruDescription + " " + Err.ruArgument + " " + Err.hidden
	}
	return Err.ruDescription + " " + Err.ruArgument + " " + Err.hidden + ": " + Err.original.Error()
}

func (E Error) IsOverlapWithError(err error) bool {
	if err == nil || err == (*Error)(nil) {
		return false
	}
	if Err, ok := err.(*Error); ok {
		if Err.code == E.code {
			return true
		}
	}
	return false
}

func (Err Error) IsOverlapWithJsonError(jsonError []byte) bool {
	if jsonError == nil {
		return false
	}
	var exported jsonExportedError

	if err := json.Unmarshal(jsonError, &exported); err != nil {
		return false
	}
	if exported.Code == Err.code {
		return true
	}
	return false
}

func (Err Error) IsInternalOrExternal() bool {
	if Err.errorType == InternalType || Err.errorType == ExternalType {
		return true
	}
	return false
}

var (
	// User errors (100 - 109)
	UserNotLogged = new(100, UnauthorizedType,
		"Пользователь не авторизован",
		"User not authorized")
	UserNotExist = new(101, BusinessType,
		"Пользователь не существует",
		"User not exists")
	AuthFail = new(102, BusinessType,
		"Не могу авторизовать пользователя. Неверная почта или пароль",
		"Cannot authorize user. Wrong mail or password")
	InvalidToken = new(103, UnauthorizedType,
		"Авторизационный токен невалиден",
		"Authorization token is invalid")
	NotConfirmedMail = new(104, BusinessType,
		"Пожалуйста, подтвердите вашу почту. Письмо выслано на ваш почтовый ящик",
		"Please confirm your mail. Mail was sent to your email address")
	RegFailUserExists = new(105, BusinessType,
		"Такой пользователь уже существует",
		"Same user already exists")
	AccessDenied = new(106, UnauthorizedType,
		"Не могу авторизоваться, доступ запрещен",
		"Cannot authorize, access denied")

	// Request errors (110 - 119)
	InvalidRequestBody = new(110, RequestType,
		"Тело запроса содержит ошибку",
		"Request body is invalid")
	NoArgument = new(111, ArgumentType,
		"Отстутствует одно из обязательных полей",
		"One of the required fields is missing")
	InvalidArgument = new(112, ArgumentType,
		"Ошибка в аргументе",
		"Argument error")

	// Common errors (120 - 129)
	RecordNotFound = new(120, BusinessType,
		"Такой записи не существует в базе данных",
		"Record not found in database")
	ImpossibleToExecute = new(121, BusinessType,
		"Невозможно выполнить команду",
		"Imposible to execute command")

	// Internal errors (130 - 139)
	MarshalError = new(130, InternalType,
		"Произошла ошибка при упаковке данных",
		"An error occurred while packing data")
	UnmarshalError = new(131, InternalType,
		"Произошла ошибка при распаковке данных",
		"An error occurred while unpacking data")
	UnknownInternalError = new(132, InternalType,
		"Произошла неизвестная ошибка",
		"An unknown error occurred")
	HashError = new(133, InternalType,
		"Ошибка при создании хэша",
		"Error during creating hash")
	NotConfiguredPackage = new(134, InternalType,
		"Не сконфигурирован пакет",
		"Not configured package")
	NotInitializedPackage = new(135, InternalType,
		"Не инициализирован пакет",
		"Not initialized package")
	ConfigurationFail = new(136, InternalType,
		"Ошибка конфигурации пакета",
		"Failed configuration of package")
	InitializationFail = new(137, InternalType,
		"Ошибка инициализации пакета",
		"Failed initialization of package")

	// External errors (140 - 149)
	DatabaseError = new(150, ExternalType,
		"База данных вернула ошибку",
		"Database returned error")
	DatabasePreparingError = new(151, ExternalType,
		"База данных вернула ошибку во время подготовки запроса",
		"Database returned error in time of preparing query")
	DatabaseExecutingError = new(152, ExternalType,
		"База данных вернула ошибку во время выполнения запроса",
		"Database returned error in time of executing query")
	DatabaseScanError = new(153, ExternalType,
		"База данных вернула ошибку во время возвращения результатов запроса",
		"Database returned error in time of scaning results of query")
	DatabaseTransactionError = new(154, ExternalType,
		"База данных вернула ошибку во время создания транзакции",
		"Database returned error in time of transaction preparing")
	RedisError = new(154, ExternalType,
		"Redis вернул ошибку",
		"Redis returned error")
	MailerError = new(155, ExternalType,
		"Mailer вернул ошибку",
		"Mailer returned error")
	ApiError = new(156, ExternalType,
		"Внешний сервис вернул ошибку",
		"External service returned error")
	EmptyResponse = new(157, ExternalType,
		"Внешний сервис не вернул данные",
		"External service returned no data")
)
