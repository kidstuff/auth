package auth

type Notificator interface {
	SendMail(ctx *AuthContext, subject, message, from, to string) error
}
