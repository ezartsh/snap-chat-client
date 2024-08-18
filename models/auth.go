package models

type Auth struct {
	accessToken string
	name        string
	username    string
}

func NewAuth(name string, username string, accessToken string) *Auth {
	return &Auth{
		accessToken: accessToken,
		name:        name,
		username:    username,
	}
}

func (a Auth) Name() string {
	if a.name == "" {
		return "ANON"
	}
	return a.name
}

func (a Auth) Username() string {
	return a.username
}

func (a Auth) Token() string {
	return a.accessToken
}
