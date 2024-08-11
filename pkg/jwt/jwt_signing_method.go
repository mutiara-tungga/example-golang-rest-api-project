package jwt

import "github.com/golang-jwt/jwt/v5"

type JWTFamily string

const (
	JWTFamilyRSA     JWTFamily = "RSA"
	JWTFamilyHMACSHA JWTFamily = "HMAC-SHA"
)

type JWTSigningMethodName string

const (
	JWTSigningMethodNameRS256 = "RS256"
	JWTSigningMethodNameRS384 = "RS384"
	JWTSigningMethodNameRS512 = "RS512"
	JWTSigningMethodNameHS256 = "HS256"
	JWTSigningMethodNameHS384 = "HS384"
	JWTSigningMethodNameHS512 = "HS512"
)

func (name JWTSigningMethodName) GetSigningMethod() jwt.SigningMethod {
	switch name {
	case JWTSigningMethodNameRS256:
		return jwt.SigningMethodRS256
	case JWTSigningMethodNameRS384:
		return jwt.SigningMethodRS384
	case JWTSigningMethodNameRS512:
		return jwt.SigningMethodRS512
	case JWTSigningMethodNameHS256:
		return jwt.SigningMethodHS256
	case JWTSigningMethodNameHS384:
		return jwt.SigningMethodHS384
	case JWTSigningMethodNameHS512:
		return jwt.SigningMethodHS512
	default:
		return nil
	}
}

func (name JWTSigningMethodName) GetFamily() JWTFamily {
	switch name {
	case JWTSigningMethodNameRS256, JWTSigningMethodNameRS384, JWTSigningMethodNameRS512:
		return JWTFamilyRSA
	case JWTSigningMethodNameHS256, JWTSigningMethodNameHS384, JWTSigningMethodNameHS512:
		return JWTFamilyHMACSHA
	default:
		return ""
	}
}
