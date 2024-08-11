package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ExpireAt  *jwt.NumericDate `json:"exp"`
	NotBefore *jwt.NumericDate `json:"nbf"`
	IssuedAt  *jwt.NumericDate `json:"iat"`
	Audience  jwt.ClaimStrings `json:"aud"`
	Issuer    string           `json:"iss"`
	Subject   string           `json:"sub"`
}

func (c JWTClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.ExpireAt, nil
}
func (c JWTClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.IssuedAt, nil
}
func (c JWTClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return c.NotBefore, nil
}
func (c JWTClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}
func (c JWTClaims) GetSubject() (string, error) {
	return c.Subject, nil
}
func (c JWTClaims) GetAudience() (jwt.ClaimStrings, error) {
	return c.Audience, nil
}
