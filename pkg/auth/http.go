package auth

const (
	SessionCookieName = "pg-session"
	UserCookieName    = "pg-userid"

	HeaderService   = "X-Pg-Service"
	HeaderSignature = "X-Pg-Signature"
	HeaderIdentity  = "X-Pg-Identity"
	HeaderEpoch     = "X-Pg-Epoch"
)
