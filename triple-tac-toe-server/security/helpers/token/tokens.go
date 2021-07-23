package tokens

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"gregvader/triple-tac-toe/domain"
	"gregvader/triple-tac-toe/globals"
	"gregvader/triple-tac-toe/security"
)

const (
	userIdKey     = "user_id"
	usernameKey   = "username"
	authorizedKey = "authorized"
	expiredKey    = "exp"
)

func ExtractAuthDetailsFromToken(r *http.Request) (security.AuthDetails, error) {
	tokenString := ExtractToken(r)
	return ExtractAuthDetailsFromTokenString(tokenString)
}

func ExtractAuthDetailsFromTokenString(tokenString string) (security.AuthDetails, error) {
	authDetails := security.AuthDetails{}
	token, verErr := VerifyToken(tokenString)
	if verErr != nil {
		return authDetails, verErr
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		userId, extErr := strconv.ParseUint(fmt.Sprintf("%.f", claims[userIdKey]), 10, 64)
		if extErr != nil {
			return authDetails, extErr
		}
		authDetails.UserId = userId
		authDetails.Username = fmt.Sprintf("%s", claims[usernameKey])

		return authDetails, nil
	}

	return authDetails, errors.New("token is not valid")
}

func CreateAuthInfo(user domain.User) (security.AuthInfo, error) {
	claims := jwt.MapClaims{}
	claims[authorizedKey] = true
	claims[userIdKey] = user.Id
	claims[usernameKey] = user.Username
	claims[expiredKey] = getExpiryDate()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := getSignedTokenString(token)

	if err != nil {
		return security.AuthInfo{}, err
	}

	return security.AuthInfo{
		Token:    signedToken,
		UserId:   user.Id,
		Username: user.Username,
	}, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

	//Token example: "Bearer TOKEN_YYY"
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(globals.JwtTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(ExtractToken(r))
	if err != nil {
		return err
	}
	if !token.Valid {
		return err
	}
	return nil
}

func getExpiryDate() int64 {
	lifeLength, err := strconv.ParseInt(globals.JwtTokenLifeLength, 10, 32)

	if err != nil {
		log.Println("JWT life length could not be loaded from globals.go.. Setting it to 24 hours..")
		lifeLength = 24
	}
	return time.Now().Add(time.Hour * time.Duration(lifeLength)).Unix()
}

func getSignedTokenString(token *jwt.Token) (string, error) {
	tkn, signErr := token.SignedString([]byte(globals.JwtTokenSecret))
	if signErr != nil {
		return "", signErr
	}

	return tkn, nil
}
