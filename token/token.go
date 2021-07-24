package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	SessionId string `json:"session id"`
	jwt.StandardClaims
}

func CreateToken(sessionId string) (string, error) {
	//h := hmac.New(sha256.New, []byte("this is the key"))
	//h.Write([]byte(sessionId))
	//sig := base64.StdEncoding.EncodeToString(h.Sum(nil))
	//
	//return sig + "|" + sessionId

	claims := CustomClaims{
		sessionId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sig, err := token.SignedString([]byte("this is the key"))
	if err != nil {
		return "", fmt.Errorf("error in SignedString: %w", err)
	}
	return sig, nil
}

func ParseToken(signature string) (string, error) {
	//xs := strings.Split(signature, "|")
	//if len(xs) != 2 {
	//	return "", fmt.Errorf("session id was rewritten by sbd")
	//}
	//messageMAC, err := base64.StdEncoding.DecodeString(xs[0])
	//if err != nil {
	//	return "", fmt.Errorf("error in decoding from base64: %w", err)
	//}
	//sessionId := xs[1]
	//
	//h := hmac.New(sha256.New, []byte("this is the key"))
	//h.Write([]byte(sessionId))
	//if !hmac.Equal(messageMAC, h.Sum(nil)) {
	//	return "", fmt.Errorf("session id was rewritten by sbd")
	//}
	//return sessionId, nil

	token, err := jwt.ParseWithClaims(signature, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("encoding method modified")
		}
		return []byte("this is the key"), nil
	})
	if err != nil {
		return "", fmt.Errorf("error in parsewithclaims: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("JWT modified")
	}
	return token.Claims.(*CustomClaims).SessionId, nil
}
