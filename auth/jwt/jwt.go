package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type claims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

type RevokedTokenChecker interface {
	RevokedTokenCheck(token string) (bool, error)
}

// Errors
var ErrTokenExpired error = errors.New("token has expired")
var ErrTokenRevoked error = errors.New("token has been revoked")

func getSigningKey() ([]byte, error) {
	if err := godotenv.Load("../.env"); err != nil {
		return nil, err
	}
	key := os.Getenv("JWT_SIGNING_KEY")
	return []byte(key), nil
}

func GenerateSignedToken(id string, exp time.Time) (string, error) {
	claims := claims{
		Id:               id,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(exp)},
	}
	unsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	mySigningKey, err := getSigningKey()
	if err != nil {
		return "", err
	}
	signed, err := unsigned.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func GenerateTokens(uid string, accExp time.Time, refExp time.Time) (string, string, error) {
	accToken, err := GenerateSignedToken(uid, accExp)
	if err != nil {
		return "", "", err
	}
	refToken, err := GenerateSignedToken(uid, refExp)
	if err != nil {
		return "", "", err
	}
	return accToken, refToken, nil
}

func ParseTokenWithoutValidation(tokenString string) (uid string, exp time.Time, err error) {
	decoded := claims{}

	_, err = jwt.ParseWithClaims(
		tokenString,
		&decoded,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v", 
					token.Header["alg"],
				)
			}
			mySigningKey, err := getSigningKey()
			if err != nil {
				return nil, fmt.Errorf("failed to get jwt signing key from .env")
			}
			return mySigningKey, nil
		},
		jwt.WithoutClaimsValidation(),
	)

	if err != nil {
		return uid, exp, err
	}

	numericDate, err := decoded.GetExpirationTime()
	if err != nil {
		return uid, exp, err
	}

	return decoded.Id, numericDate.Time, nil
}

func VerifyToken(tokenString string) (uid string, exp time.Time, err error) {
	decoded := claims{}
	_, err = jwt.ParseWithClaims(
		tokenString,
		&decoded,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v", 
					token.Header["alg"],
				)
			}
			mySigningKey, err := getSigningKey()
			if err != nil {
				return nil, fmt.Errorf("failed to get jwt signing key from .env")
			}
			return mySigningKey, nil
		},
	)

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return uid, exp, ErrTokenExpired
		}
		return uid, exp, err
	}

	numericDate, err := decoded.GetExpirationTime()
	if err != nil {
		return uid, exp, err
	}

	return decoded.Id, numericDate.Time, nil
}

func Refresh(refToken string, accExp time.Time, refExp time.Time) (string, string, error) {
	uid, _, err := VerifyToken(refToken)
	if err != nil {
		return "", "", err
	}
	newAccToken, err := GenerateSignedToken(uid, accExp)
	if err != nil {
		return "", "", err
	}
	newRefToken, err := GenerateSignedToken(uid, refExp)
	if err != nil {
		return "", "", err
	}
	return newAccToken, newRefToken, nil
}

func CalcTokenExpTime(envPath string) (accExp time.Time, refExp time.Time, err error) {
	if err := godotenv.Load(envPath); err != nil {
		return accExp, refExp, err
	}
	accExpIn, err := strconv.Atoi(os.Getenv("JWT_ACC_EXP_IN"))
	if err != nil {
		return accExp, refExp, err
	}
	refExpIn, err := strconv.Atoi(os.Getenv("JWT_REF_EXP_IN"))
	if err != nil {
		return accExp, refExp, err
	}
	accExp = time.Now().Add(time.Hour * time.Duration(accExpIn))
	refExp = time.Now().Add(time.Hour * time.Duration(refExpIn))
	return accExp, refExp, nil
}
