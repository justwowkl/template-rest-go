package util

import (
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/validator.v9"
)

// KeyJWT jwt key
var KeyJWT = []byte("hihi")
var keyAES = "hihi"
var expireHour = time.Duration(24)
var jwtValidate = validator.New()

// JwtData data
type JwtData struct {
	ID int    `validate:"required"`
	IP string `validate:"required"`
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
	token.Claims = jwt.MapClaims{
		"data": dataEncrypted,
		"exp":  time.Now().Add(time.Hour * expireHour).Unix(),
	}
	// claims := token.Claims.(jwt.MapClaims)
	// claims["data"] = dataEncrypted
	// claims["exp"] = time.Now().Add(time.Hour * expireHour).Unix()

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString(KeyJWT)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// JwtVerify verify data
func JwtVerify(tokenString string, ip string) (JwtData, error) {

	// parse jwt token
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return KeyJWT, nil
	})
	if err != nil {
		return JwtData{}, err
	}
	if !token.Valid {
		return JwtData{}, jwt.NewValidationError("not validated", jwt.ValidationErrorMalformed)
	}

	// get encrypted Data from claim
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return JwtData{}, jwt.NewValidationError("not validated", jwt.ValidationErrorMalformed)
	}
	dataEncrypted, ok := claims["data"].(string)
	if !ok {
		return JwtData{}, jwt.NewValidationError("not validated", jwt.ValidationErrorMalformed)
	}

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
	// check user IP
	if data.IP != ip {
		return JwtData{}, err
	}

	return data, nil

}

// // JwtVerifyData verify data
// func JwtVerifyData(dataEncrypted string) (JwtData, error) {

// 	// decrypt to json
// 	dataJSON, err := Decrypt(dataEncrypted, keyAES)
// 	if err != nil {
// 		return JwtData{}, err
// 	}

// 	// json to struct
// 	var data JwtData
// 	err = json.Unmarshal([]byte(dataJSON), &data)
// 	if err != nil {
// 		return JwtData{}, err
// 	}

// 	// validate struct
// 	if err := jwtValidate.Struct(data); err != nil {
// 		return JwtData{}, err
// 	}

// 	return data, nil

// }
