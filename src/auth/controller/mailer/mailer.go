package mailer

import (
	"auth_backend/errors"
	"auth_backend/model"
	"net/smtp"
	"strconv"
)

type Config struct {
	Host   string `conf:"host"`
	Email  string `conf:"email"`
	Passwd string `conf:"passwd"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}

func getConfig() (*Config, *errors.Error) {
	if cfg == nil {
		return nil, errors.NotConfiguredPackage.SetArgs("controller/mailer", "controller/mailer")
	}
	return cfg, nil
}

func SendEmailConfirmMessage(user *model.UserBasic, token, serverIp string, serverPort uint) *errors.Error {
	conf, Err := getConfig()
	if Err != nil {
		return Err
	}

	portString := strconv.FormatUint(uint64(serverPort), 10)

	auth := smtp.PlainAuth("", conf.Email, conf.Passwd, conf.Host)
	message := `To: <` + *user.Email + `>
From: "Hypertube administration" <` + conf.Email + `>
Subject: Confirm email in project Hypertube
MIME-Version: 1.0
Content-type: text/html; charset=utf8

<html><head></head><body>
<span style="font-size: 1.3em; color: green;">Hello, ` + user.Username + `, click below to confirm your email
<form method="GET" action="http://`+serverIp+`:` + portString + `/api/email/confirm">
	<input type="hidden" name="code" value="` + token + `">
	<input type="submit" value="Click to confirm mail">
</form>
<a target="_blank" href="http://`+serverIp+`:` + portString + `/api/email/confirm?code=` + token + `">click to confirm mail</a></br>
</br>
if this letter came by mistake - delete it 
</span></body></html>
`

	if err := smtp.SendMail(conf.Host+":587", auth, conf.Email, []string{*user.Email}, []byte(message)); err != nil {
		return errors.MailerError.SetOrigin(err)
	}
	return nil
}

func SendEmailPatchMailAddress(user *model.UserBasic, token, serverIp string, serverPort uint) *errors.Error {
	conf, Err := getConfig()
	if Err != nil {
		return Err
	}

	if user.NewEmail == nil {
		return errors.MailerError.SetArgs("Поле NewEmail равно nil", "NewEmail field is nil")
	}

	portString := strconv.FormatUint(uint64(serverPort), 10)

	auth := smtp.PlainAuth("", conf.Email, conf.Passwd, conf.Host)
	message := `To: <` + *user.NewEmail + `>
From: "Hypertube administration" <` + conf.Email + `>
Subject: Confirm new email in project Hypertube
MIME-Version: 1.0
Content-type: text/html; charset=utf8

<html><head></head><body>
<span style="font-size: 1.3em; color: green;">Hello, ` + user.Username + `, click below to confirm your new email
<form method="GET" action="http://`+serverIp+`:` + portString + `/api/email/patch/confirm">
	<input type="hidden" name="code" value="` + token + `">
	<input type="submit" value="Click to confirm mail">
</form>
<a target="_blank" href="http://`+serverIp+`:` + portString + `/api/email/patch/confirm?code=` + token + `">click to confirm mail</a></br>
</br>
if this letter came by mistake - delete it 
</span></body></html>
`

	if err := smtp.SendMail(conf.Host+":587", auth, conf.Email, []string{*user.NewEmail}, []byte(message)); err != nil {
		return errors.MailerError.SetOrigin(err)
	}
	return nil
}

func SendEmailPasswdRepair(user *model.UserBasic, repairToken, serverIp string, serverPort uint) *errors.Error {
	conf, Err := getConfig()
	if Err != nil {
		return Err
	}

	portString := strconv.FormatUint(uint64(serverPort), 10)

	auth := smtp.PlainAuth("", conf.Email, conf.Passwd, conf.Host)
	message := `To: <` + *user.Email + `>
From: "Hypertube administration" <` + conf.Email + `>
Subject: Password repair in project Hypertube
MIME-Version: 1.0
Content-type: text/html; charset=utf8

<html><head></head><body>
<span style="font-size: 1.3em; color: green;">Hello, ` + user.Username + `, click below to confirm password repair of your account
<form method="GET" action="http://`+serverIp+`:` + portString + `/api/passwd/repair/confirm">
	<input type="hidden" name="code" value="` + repairToken + `">
	<input type="submit" value="Click to confirm mail">
</form>
<a target="_blank" href="http://`+serverIp+`:` + portString + `/api/passwd/repair/confirm?code=` + repairToken +
	`">click to repair password</a></br></br>
</br>
if this letter came by mistake - delete it 
</span></body></html>
`

	if err := smtp.SendMail(conf.Host+":587", auth, conf.Email, []string{*user.Email}, []byte(message)); err != nil {
		return errors.MailerError.SetOrigin(err)
	}
	return nil
}
