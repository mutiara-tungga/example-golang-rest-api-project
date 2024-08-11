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
type JWTParser interface {
	ParseAndValidate(ctx context.Context, tokenString string) (JWTClaims, error)
}

var (
	ErrorJWTInvalid = pkgErr.NewCustomError("jwt not valid", "INVALID_JWT", http.StatusUnauthorized)
)

type JWTParserOptions func(*jwtParser) error

// JWTParserWithSigningMethod assign jwt signing method and jwt key
// IF method is RSA key should fill by public key
func JWTParserWithSigningMethod(methodName JWTSigningMethodName, key string) JWTParserOptions {
	return func(jp *jwtParser) error {
		jp.jwtSigningMethodName = methodName
		jp.signingMethod = methodName.GetSigningMethod()

		jp.jwtKeyString = key
		switch jp.jwtSigningMethodName.GetFamily() {
		case JWTFamilyHMACSHA:
			jp.jwtKey = []byte(jp.jwtKeyString)

		case JWTFamilyRSA:
			privateKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(jp.jwtKeyString))
			if err != nil {
				return err
			}
			jp.jwtKey = privateKey

		default:
			return ErrFailedJWTMethodNotSupported
		}

		return nil
	}
}

func JWTParserWithValidIssuer(issuer string) JWTParserOptions {
	return func(jp *jwtParser) error {
		jp.validIssuer = issuer
		return nil
	}
}

func NewJWTParser(options ...JWTParserOptions) jwtParser {
	v := &jwtParser{}

	for _, apply := range options {
		err := apply(v)

		if err != nil {
			panic(err)
		}
	}

	// TODO: validation

	return *v
}

type jwtParser struct {
	jwtKeyString string
	// IF signing method is RSA should be filled by public key
	jwtKey               any
	jwtSigningMethodName JWTSigningMethodName
	signingMethod        jwt.SigningMethod
	validIssuer          string
}

func (jp jwtParser) ParseAndValidate(ctx context.Context, tokenString string) (JWTClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (any, error) { return jp.jwtKey, nil },
		jwt.WithValidMethods([]string{jp.signingMethod.Alg()}),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(jp.validIssuer),
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
