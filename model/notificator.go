package model

var DEFAULT_NOTIFICATOR Notificator

type Notificator interface {
	SendMail(subject, message, from, to string) error
}
