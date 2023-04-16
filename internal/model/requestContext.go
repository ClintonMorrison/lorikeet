package model

type RequestContext struct {
	Username     string
	Ip           string
	Password     string // password from basic auth field
	SessionToken string // session token from cookie
}

func (rc RequestContext) ToAuth(decryptToken string) Auth {
	return Auth{
		Username: rc.Username,
		Password: decryptToken,
		Ip:       rc.Ip,
	}
}
