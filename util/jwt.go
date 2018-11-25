package util

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/md5"
// 	"crypto/rand"
// 	"encoding/hex"
// 	"encoding/json"
// 	"io"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/labstack/echo"
// 	"gitlab.com/hartsfield/gencrypt"
// )

// var keyAES = []byte("safh97hf9fhja98ewfj94fjh98djdfjh")
// var gcm *gencrypt.Galois

// // JwtCreate create new jwt token
// func JwtCreate(data map[string]interface{}) (string, error) {

// 	gcm, _ := gencrypt.NewGCM(keyAES)

// 	// Create token
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	// convert to json string
// 	dataString, err := json.Marshal(data)
// 	if err != nil {
// 		return "", err
// 	}

// 	// crypt json string
// 	dataCrypted := string(_cryptEncrypt([]byte(dataString)))
// 	println("dataCrypted : ", dataCrypted)

// 	// Set claims
// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
// 	claims["data"] = dataCrypted

// 	// Generate encoded token and send it as response.
// 	t, err := token.SignedString([]byte("secret"))
// 	if err != nil {
// 		return "", err
// 	}
// 	println("data is : ", string(_cryptDecrypt([]byte(dataCrypted))))
// 	return t, nil
// }

// // MWjwtVerify middleware adds a `Server` header to the response.
// func MWjwtVerify(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		user := c.Get("user").(*jwt.Token)
// 		claims := user.Claims.(jwt.MapClaims)
// 		dataCrypted := claims["data"].(string)
// 		println("dataCrypted : ", dataCrypted)
// 		dataString := string(_cryptDecrypt([]byte(dataCrypted)))
// 		println("data is : ", dataString)
// 		// convert to json string
// 		// dataString, err := json.Unmarshal()
// 		// if err != nil {
// 		// 	return "", err
// 		// }
// 		// c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
// 		return next(c)
// 	}
// }

// // https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/

// func _cryptCreateHash(key string) string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(key))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// func _cryptEncrypt(data []byte) []byte {
// 	block, _ := aes.NewCipher([]byte(_cryptCreateHash(keyAES)))
// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
// 		panic(err.Error())
// 	}
// 	ciphertext := gcm.Seal(nonce, nonce, data, nil)
// 	return ciphertext
// }

// func _cryptDecrypt(data []byte) []byte {
// 	println("0")
// 	key := []byte(_cryptCreateHash(keyAES))
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		println("1")
// 		panic(err.Error())
// 	}
// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		println("2")
// 		panic(err.Error())
// 	}
// 	nonceSize := gcm.NonceSize()
// 	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
// 	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		println("3")
// 		panic(err.Error())
// 	}
// 	return plaintext
// }
