package util

import (
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/validator.v9"
)

var keyAES = "hihi"
var keyJWT = "hihi"
var expireHour = time.Duration(24)
var jwtValidate = validator.New()

// JwtData data
type JwtData struct {
	ID int `validate:"required"`
}

// JwtCreate create new jwt token
func JwtCreate(data JwtData) (string, error) {

	// validate struct
	if err := jwtValidate.Struct(data); err != nil {
		return "", err
	}

	// data to json
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	// encrypt json
	dataEncrypted, err := Encrypt(string(dataJSON), keyAES)
	if err != nil {
		return "", err
	}
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["data"] = dataEncrypted
	claims["exp"] = time.Now().Add(time.Hour * expireHour).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(keyJWT))
	if err != nil {
		return "", err
	}

	return t, nil

}

// JwtVerifyData verify data
func JwtVerifyData(dataEncrypted string) (JwtData, error) {

	// decrypt to json
	dataJSON, err := Decrypt(dataEncrypted, keyAES)
	if err != nil {
		return JwtData{}, err
	}

	// json to struct
	var data JwtData
	err = json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		return JwtData{}, err
	}

	// validate struct
	if err := jwtValidate.Struct(data); err != nil {
		return JwtData{}, err
	}

	return data, nil

}
