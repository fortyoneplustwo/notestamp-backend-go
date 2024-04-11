package auth

import (
	"errors"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey []byte = []byte("AllYourBase") // TODO: put in env file

type claims struct {
  Id int `json:"id"`
  jwt.RegisteredClaims
}


// Errors
var ErrTokenExpired error = errors.New("Token has expired")
var ErrTokenRevoked error = errors.New("Token has been revoked")


func GenerateSignedToken(id int, exp time.Time) (string, error) {
  claims := claims{
    Id: id,
    RegisteredClaims: jwt.RegisteredClaims{ ExpiresAt: jwt.NewNumericDate(exp) },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  signed, err := token.SignedString(mySigningKey)
  if err != nil {
    return "", err
  }

  return signed, nil
}


func VerifyToken(tokenString string) (int, time.Time, error) {
  decoded, err := jwt.ParseWithClaims(
    tokenString, 
    &claims{}, 
    func(token *jwt.Token) (interface{}, error) {
      if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
      }
      return mySigningKey, nil
  })

  if err != nil {
    if err == jwt.ErrTokenExpired {
      return 0, time.Time{}, ErrTokenExpired
    }
    return 0, time.Time{}, err
  }

  claims, ok := decoded.Claims.(*claims)
  if !ok {
    return 0, time.Time{}, errors.New("Could not decode claims")
  }

  exp, _ := claims.GetExpirationTime()
  var t time.Time = exp.Time
  return claims.Id, t, nil
}


func Refresh(token string, exp time.Time, s RevokedStore) (string, error) {
  id, exp, err := VerifyToken(token)
  if err != nil {
    return "", err
  }
  
  revoked, err := s.IsRevoked(token, id)
  if err != nil {
    return "", err
  }
  if revoked {
    return "", ErrTokenRevoked
  }

  newToken, err := GenerateSignedToken(id, exp)
  if err != nil {
    return "", err
  }

  return newToken, nil
}
