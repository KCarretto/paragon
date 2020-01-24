package auth

const (
	OAuthStateLength   = 64
	SessionTokenLength = 256
	SessionCookieName  = "pg-session"
	UserCookieName     = "pg-userid"
)

type Secret string

func (secret Secret) String() string {
	return "[REDACTED]"
}
