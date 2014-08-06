package model

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/securecookie"
	"time"
)

var HashPwd = DefaultHashPwd

func DefaultHashPwd(pwd string) (Password, error) {
	p := Password{}
	p.InitAt = time.Now()
	p.Salt = securecookie.GenerateRandomKey(32)

	pwdBytes := []byte(pwd)
	tmp := make([]byte, len(pwdBytes)+len(p.Salt))
	copy(tmp, pwdBytes)
	tmp = append(tmp, p.Salt...)
	b, err := bcrypt.GenerateFromPassword(tmp, bcrypt.DefaultCost)
	p.Hashed = b

	return p, err
}
