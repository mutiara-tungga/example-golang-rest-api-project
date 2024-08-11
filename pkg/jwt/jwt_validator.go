package jwt

import (
	"context"
	"fmt"
	pkgErr "golang-rest-api/pkg/error"
	"golang-rest-api/pkg/log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -destination=mock/jwt_validator.go -package=mock golang-rest-api/pkg/jwt JWTValidator
type JWTValidator interface {
	Validate(ctx context.Context, tokenString string) (JWTClaims, error)
}

var (
	ErrorJWTInvalid = pkgErr.NewCustomError("jwt not valid", "INVALID_JWT", http.StatusUnauthorized)
)

type JWTValidatorOptions func(*jwtValidator) error

// JWTValidatorWithSigningMethod assign jwt signing method and jwt key
// IF method is RSA key should fill by public key
func JWTValidatorWithSigningMethod(methodName JWTSigningMethodName, key string) JWTValidatorOptions {
	return func(jv *jwtValidator) error {
		jv.jwtSigningMethodName = methodName
		jv.signingMethod = methodName.GetSigningMethod()

		jv.jwtKeyString = key
		switch jv.jwtSigningMethodName.GetFamily() {
		case JWTFamilyHMACSHA:
			jv.jwtKey = []byte(jv.jwtKeyString)

		case JWTFamilyRSA:
			privateKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(jv.jwtKeyString))
			if err != nil {
				return err
			}
			jv.jwtKey = privateKey

		default:
			return ErrFailedJWTMethodNotSupported
		}

		return nil
	}
}

func JWTValidatorWithValidIssuer(issuer string) JWTValidatorOptions {
	return func(jv *jwtValidator) error {
		jv.validIssuer = issuer
		return nil
	}
}

func NewJWTValidator(options ...JWTValidatorOptions) jwtValidator {
	v := &jwtValidator{}

	for _, apply := range options {
		err := apply(v)

		if err != nil {
			panic(err)
		}
	}

	// TODO: validation

	return *v
}

type jwtValidator struct {
	jwtKeyString string
	// IF signing method is RSA should be filled by public key
	jwtKey               any
	jwtSigningMethodName JWTSigningMethodName
	signingMethod        jwt.SigningMethod
	validIssuer          string
}

func (jv jwtValidator) Validate(ctx context.Context, tokenString string) (JWTClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (any, error) { return jv.jwtKey, nil },
		jwt.WithValidMethods([]string{jv.signingMethod.Alg()}),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(jv.validIssuer),
	)
	if err != nil {
		log.Error(ctx, "error parse claim when validate jwt", err)
		return JWTClaims{}, pkgErr.NewCustomErrWithOriginalErr(ErrorJWTInvalid, err)
	}

	jwtClaims, ok := jwtToken.Claims.(*JWTClaims)
	if !ok {
		err := fmt.Errorf("invalid jwt claims")
		log.Error(ctx, "error casting jwt claims when validate jwt", err)
		return JWTClaims{}, pkgErr.NewCustomErrWithOriginalErr(ErrorJWTInvalid, err)
	}

	return *jwtClaims, nil
}
