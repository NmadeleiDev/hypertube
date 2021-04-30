package validator

import (
	"auth_backend/errors"
	"unicode/utf8"
)

func isLetter(c rune) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= 'A' && c <= 'Z' {
		return true
	}
	if c >= 'а' && c <= 'я' {
		return true
	}
	if c >= 'А' && c <= 'Я' {
		return true
	}
	return false
}

func isDigit(c rune) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func isNameRunePermitted(c rune) bool {
	if isLetter(c) || isDigit(c) {
		return true
	}
	if c == '_' || c == '-' || c == ' ' {
		return true
	}
	return false
}

func isMailRunePermitted(c rune) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= 'A' && c <= 'Z' {
		return true
	}
	if isDigit(c) {
		return true
	}
	if c == '_' || c == '-' || c == '.' || c == '@' {
		return true
	}
	return false
}

func ValidatePassword(pass string) *errors.Error {
	var (
		wasLetter      bool
		wasDigit       bool
		wasSpacialChar bool
		buf            = []rune(pass)
	)
	if utf8.RuneCountInString(pass) < 6 {
		return errors.InvalidArgument.SetArgs("слишком короткий пароль", "too short password")
	}

	if utf8.RuneCountInString(pass) > 25 {
		return errors.InvalidArgument.SetArgs("слишком короткий пароль", "too short password")
	}

	for i := 0; i < len(buf); i++ {
		if isLetter(buf[i]) {
			wasLetter = true
		}
		if buf[i] >= '0' && buf[i] <= '9' {
			wasDigit = true
		}
		if buf[i] == '!' || buf[i] == '@' || buf[i] == '#' || buf[i] == '$' ||
			buf[i] == '%' || buf[i] == '^' || buf[i] == '&' || buf[i] == '*' {
			wasSpacialChar = true
		}
	}
	if !wasLetter {
		return errors.InvalidArgument.SetArgs("пароль должен содержать буквы", "Password should contain letters")
	}
	if !wasDigit {
		return errors.InvalidArgument.SetArgs("пароль должен содержать цифры", "Password should contain digits")
	}
	if !wasSpacialChar {
		return errors.InvalidArgument.SetArgs("пароль должен содержать специальные символы",
			"Password should contain special chars")
	}
	return nil
}

func ValidateEmail(email string) *errors.Error {
	var (
		buf        = []rune(email)
		length     = len(buf)
		doggyCount int
		dots       int
	)

	if utf8.RuneCountInString(email) < 6 {
		return errors.InvalidArgument.SetArgs("слишком короткий почтовый адрес", "too short mail address")
	}
	if utf8.RuneCountInString(email) > 40 {
		return errors.InvalidArgument.SetArgs("слишком длинный почтовый адрес", "too long mail address")
	}

	if buf[0] == '_' || buf[0] == '-' || buf[0] == '@' ||
		buf[0] == '.' || (buf[0] >= '0' && buf[0] <= '9') {
		return errors.InvalidArgument.SetArgs("первый символ почтового адреса невалиден",
			"invalid first mail address symbol")
	}

	if buf[length-1] == '_' || buf[length-1] == '-' || buf[length-1] == '@' ||
		buf[length-1] == '.' || (buf[length-1] >= '0' && buf[length-1] <= '9') {
		return errors.InvalidArgument.SetArgs("последний символ почтового адреса невалиден",
			"invalid last mail address symbol")
	}

	for i := 0; i < length; i++ {
		if !isMailRunePermitted(buf[i]) {
			return errors.InvalidArgument.SetArgs("найден запрещенный символ почтового адреса",
				"forbidden symbol in mail")
		}
		if buf[i] == '@' {
			doggyCount++
			if i > 0 && buf[i-1] == '.' {
				return errors.InvalidArgument.SetArgs("почтовый адрес невалиден", "invalid mail address")
			}
		}
		if buf[i] == '.' && doggyCount > 0 {
			dots++
			if buf[i-1] == '.' || buf[i-1] == '@' {
				return errors.InvalidArgument.SetArgs("почтовый адрес невалиден", "invalid mail address")
			}
		}
	}
	if doggyCount != 1 {
		return errors.InvalidArgument.SetArgs("невалидное количество символов '@' в почтовом адресе",
			"invalid amount of '@' symbols in mail address")
	}
	if dots != 1 && dots != 2 {
		return errors.InvalidArgument.SetArgs("невалидное количество символов '.' в почтовом адресе",
			"invalid amount of '.' symbols in mail address")
	}
	return nil
}

func ValidateName(name string) *errors.Error {
	if utf8.RuneCountInString(name) < 3 {
		return errors.InvalidArgument.SetArgs("слишком короткое имя", "too short name")
	}
	if utf8.RuneCountInString(name) > 20 {
		return errors.InvalidArgument.SetArgs("слишком длинное имя", "too long name")
	}
	return nil
}
