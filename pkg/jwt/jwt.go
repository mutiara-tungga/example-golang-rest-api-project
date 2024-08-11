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

//go:generate mockgen -destination=mock/jwt.go -package=mock transport-service/pkg/jwt JWT
type JWTGenerator interface {
	GenerateJWT(ctx context.Context, user User) (JWTResult, error)
}

type JWTGeneratorOptions func(*jwtGenerator) error

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

func JWTGeneratorWithExpireDurationInSecond(durationInSecond int64) JWTGeneratorOptions {
	return func(jg *jwtGenerator) error {
		jg.expireDurationInSecond = durationInSecond
		return nil
	}
}

func JWTGeneratorWithRefreshTokenExpireDurationInSecond(durationInSecond int64) JWTGeneratorOptions {
	return func(jg *jwtGenerator) error {
		jg.refreshTokenExpireDurationInSecond = durationInSecond
		return nil
	}
}

func JWTGeneratorWithIssuer(issuer string) JWTGeneratorOptions {
	return func(jg *jwtGenerator) error {
		jg.issuer = issuer
		return nil
	}
}

func New(options ...JWTGeneratorOptions) jwtGenerator {
	defaultExpireDuration := int64((24 * time.Hour).Seconds())
	defaultRefreshTokenExpireDuration := int64((48 * time.Hour).Seconds())
	gen := &jwtGenerator{
		timeNowFunc:                        time.Now,
		expireDurationInSecond:             defaultExpireDuration,
		refreshTokenExpireDurationInSecond: defaultRefreshTokenExpireDuration,
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
	jwtKeyString                       string
	jwtKey                             any
	jwtSigningMethodName               JWTSigningMethodName
	signingMethod                      jwt.SigningMethod
	expireDurationInSecond             int64
	refreshTokenExpireDurationInSecond int64
	issuer                             string
	timeNowFunc                        func() time.Time
}

var (
	MapClaimsKeyUserID   string = "user_id"
	MapClaimsKeyUsername string = "username"
	MapClaimsKeyExpireAt string = "exp"
	MapClaimsKeyIssuedAt string = "iat"
	MapClaimsKeyIssuer   string = "issuer"
)

func (jg jwtGenerator) GenerateJWT(ctx context.Context, u User) (JWTResult, error) {
	nowUnix := jg.timeNowFunc().Unix()
	tokenExpiresUnix := nowUnix + jg.expireDurationInSecond
	claims := jwt.MapClaims{
		MapClaimsKeyUserID:   u.ID,
		MapClaimsKeyUsername: u.Username,
		MapClaimsKeyExpireAt: tokenExpiresUnix,
		MapClaimsKeyIssuedAt: nowUnix,
		MapClaimsKeyIssuer:   jg.issuer,
	}

	token := jwt.NewWithClaims(jg.signingMethod, claims)
	tokenString, err := token.SignedString(jg.jwtKey)
	if err != nil {
		log.Error(ctx, "error get signed string access token", err)
		return JWTResult{}, pkgErr.NewCustomErrWithOriginalErr(ErrFailedProcessJWT, err)
	}

	refreshTokenClaims := make(jwt.MapClaims)
	for k, v := range claims {
		refreshTokenClaims[k] = v
	}

	refreshTokenExpireUnix := nowUnix + jg.refreshTokenExpireDurationInSecond
	refreshTokenClaims[MapClaimsKeyExpireAt] = refreshTokenExpireUnix

	refreshToken := jwt.NewWithClaims(jg.signingMethod, claims)
	refreshTokenString, err := refreshToken.SignedString(jg.jwtKey)
	if err != nil {
		log.Error(ctx, "error get signed string refresh token", err)
		return JWTResult{}, pkgErr.NewCustomErrWithOriginalErr(ErrFailedProcessJWT, err)
	}

	return JWTResult{
		AccessToken:           tokenString,
		ExpiresAt:             time.Unix(refreshTokenExpireUnix, 0),
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresAt: time.Unix(refreshTokenExpireUnix, 0),
	}, nil
}
