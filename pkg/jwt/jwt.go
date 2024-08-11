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
	ErrFailedProcessJWT = pkgErr.NewCustomError("Failed Process JWT", "FAILED_PROCESS_JWT", http.StatusInternalServerError)
)

type User struct {
	ID       string
	Username string
}

type JWTResult struct {
	AccessToken  string
	RefreshToken string
}

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

//go:generate mockgen -destination=mock/jwt.go -package=mock transport-service/pkg/jwt JWT
type JWTGenerator interface {
	GenerateJWT(ctx context.Context, user User) (JWTResult, error)
}

type JWTGeneratorOptions func(*jwtGenerator)

func JWTGeneratorWithKey(key string) JWTGeneratorOptions {
	return func(jg *jwtGenerator) {
		jg.jwtKey = key
	}
}

func JWTGeneratorWithSigningMethod(name JWTSigningMethodName) JWTGeneratorOptions {
	return func(jg *jwtGenerator) {
		jg.jwtSigningMethodName = name
		jg.signingMethod = name.GetSigningMethod()
	}
}

func JWTGeneratorWithExpireDurationInSecond(durationInSecond int64) JWTGeneratorOptions {
	return func(jg *jwtGenerator) {
		jg.expireDurationInSecond = durationInSecond
	}
}

func JWTGeneratorWithRefreshTokenExpireDurationInSecond(durationInSecond int64) JWTGeneratorOptions {
	return func(jg *jwtGenerator) {
		jg.refreshTokenExpireDurationInSecond = durationInSecond
	}
}

func JWTGeneratorWithIssuer(issuer string) JWTGeneratorOptions {
	return func(jg *jwtGenerator) {
		jg.issuer = issuer
	}
}

func New(options []JWTGeneratorOptions) jwtGenerator {
	defaultExpireDuration := int64((24 * time.Hour).Seconds())
	gen := &jwtGenerator{
		timeNowFunc:                        time.Now,
		expireDurationInSecond:             defaultExpireDuration,
		refreshTokenExpireDurationInSecond: defaultExpireDuration,
	}

	for _, apply := range options {
		apply(gen)
	}

	// TODO: validation

	return *gen
}

type jwtGenerator struct {
	jwtKey                             string
	jwtSigningMethodName               JWTSigningMethodName
	signingMethod                      jwt.SigningMethod
	expireDurationInSecond             int64
	refreshTokenExpireDurationInSecond int64
	issuer                             string
	timeNowFunc                        func() time.Time
}

var (
	MapClaimsKeyID       string = "id"
	MapClaimsKeyUsername string = "username"
	MapClaimsKeyExpireAt string = "exp"
	MapClaimsKeyIssuedAt string = "iat"
	MapClaimsKeyIssuer   string = "issuer"
)

func (jg jwtGenerator) GenerateJWT(ctx context.Context, u User) (JWTResult, error) {
	nowUnix := jg.timeNowFunc().Unix()
	claims := jwt.MapClaims{
		MapClaimsKeyID:       u.ID,
		MapClaimsKeyUsername: u.Username,
		MapClaimsKeyExpireAt: nowUnix + jg.expireDurationInSecond,
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
	refreshTokenClaims[MapClaimsKeyExpireAt] = nowUnix + jg.refreshTokenExpireDurationInSecond

	refreshToken := jwt.NewWithClaims(jg.signingMethod, claims)
	refreshTokenString, err := refreshToken.SignedString(jg.jwtKey)
	if err != nil {
		log.Error(ctx, "error get signed string refresh token", err)
		return JWTResult{}, pkgErr.NewCustomErrWithOriginalErr(ErrFailedProcessJWT, err)
	}

	return JWTResult{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
