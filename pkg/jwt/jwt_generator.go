package jwt

import (
	"context"
	"net/http"
	"time"

	pkgErr "golang-rest-api/pkg/error"
	"golang-rest-api/pkg/log"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrFailedProcessJWT            = pkgErr.NewCustomError("Failed Process JWT", "FAILED_PROCESS_JWT", http.StatusInternalServerError)
	ErrFailedJWTMethodNotSupported = pkgErr.NewCustomError("JWT Method Not Supported", "JWT_METHOD_NOT_SUPPORTED", http.StatusInternalServerError)
)

type User struct {
	ID       string
	Username string
}

type JWTResult struct {
	AccessToken           string
	ExpiresAt             time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

//go:generate mockgen -destination=mock/jwt_generator.go -package=mock golang-rest-api/pkg/jwt JWTGenerator
type JWTGenerator interface {
	GenerateJWT(ctx context.Context, user User) (JWTResult, error)
}

type JWTGeneratorOptions func(*jwtGenerator) error

// JWTGeneratorWithSigningMethod assign jwt signing method and jwt key
// IF method is RSA key should fill by private key
func JWTGeneratorWithSigningMethod(methodName JWTSigningMethodName, key string) JWTGeneratorOptions {
	return func(jg *jwtGenerator) error {
		jg.jwtSigningMethodName = methodName
		jg.signingMethod = methodName.GetSigningMethod()

		jg.jwtKeyString = key
		switch jg.jwtSigningMethodName.GetFamily() {
		case JWTFamilyHMACSHA:
			jg.jwtKey = []byte(jg.jwtKeyString)

		case JWTFamilyRSA:
			privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(jg.jwtKeyString))
			if err != nil {
				return err
			}
			jg.jwtKey = privateKey

		default:
			return ErrFailedJWTMethodNotSupported
		}

		return nil
	}
}

func JWTGeneratorWithExpireDuration(duration time.Duration) JWTGeneratorOptions {
	return func(jg *jwtGenerator) error {
		jg.expireDuration = duration
		return nil
	}
}

func JWTGeneratorWithRefreshTokenExpireDurationI(duration time.Duration) JWTGeneratorOptions {
	return func(jg *jwtGenerator) error {
		jg.refreshTokenExpireDuration = duration
		return nil
	}
}

func JWTGeneratorWithIssuer(issuer string) JWTGeneratorOptions {
	return func(jg *jwtGenerator) error {
		jg.issuer = issuer
		return nil
	}
}

func NewJWTGenerator(options ...JWTGeneratorOptions) jwtGenerator {
	gen := &jwtGenerator{
		timeNowFunc:                time.Now,
		expireDuration:             24 * time.Hour,
		refreshTokenExpireDuration: 48 * time.Hour,
	}

	for _, apply := range options {
		err := apply(gen)

		if err != nil {
			panic(err)
		}
	}

	// TODO: validation

	return *gen
}

type jwtGenerator struct {
	jwtKeyString string
	// IF signing method is RSA should be filled by private key
	jwtKey                     any
	jwtSigningMethodName       JWTSigningMethodName
	signingMethod              jwt.SigningMethod
	expireDuration             time.Duration
	refreshTokenExpireDuration time.Duration
	issuer                     string
	timeNowFunc                func() time.Time
}

func (jg jwtGenerator) GenerateJWT(ctx context.Context, u User) (JWTResult, error) {
	now := jg.timeNowFunc()
	claims := JWTClaims{
		ExpireAt: &jwt.NumericDate{Time: now.Add(jg.expireDuration)},
		IssuedAt: &jwt.NumericDate{Time: now},
		Issuer:   jg.issuer,
		Subject:  u.ID,
	}

	token := jwt.NewWithClaims(jg.signingMethod, claims)
	tokenString, err := token.SignedString(jg.jwtKey)
	if err != nil {
		log.Error(ctx, "error get signed string access token", err)
		return JWTResult{}, pkgErr.NewCustomErrWithOriginalErr(ErrFailedProcessJWT, err)
	}

	refreshTokenClaims := claims
	refreshTokenClaims.ExpireAt = &jwt.NumericDate{Time: now.Add(jg.refreshTokenExpireDuration)}
	refreshToken := jwt.NewWithClaims(jg.signingMethod, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jg.jwtKey)
	if err != nil {
		log.Error(ctx, "error get signed string refresh token", err)
		return JWTResult{}, pkgErr.NewCustomErrWithOriginalErr(ErrFailedProcessJWT, err)
	}

	return JWTResult{
		AccessToken:           tokenString,
		ExpiresAt:             claims.ExpireAt.Time,
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresAt: refreshTokenClaims.ExpireAt.Time,
	}, nil
}
