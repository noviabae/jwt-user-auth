package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// type PayloadToken struct:
// Ini adalah deklarasi struktur dengan nama PayloadToken yang digunakan untuk menyimpan payload token. Struktur ini memiliki dua bidang, yaitu AuthId yang merupakan ID autentikasi dan Expired yang merupakan waktu kedaluwarsa token.
type PayloadToken struct {
	AuthId  int
	Expired time.Time
}

// Secretkey di dapat dari http://www.unit-conversion.info/texttools/random-string-generator/
// dengan memasukkan length 32
// const SecretKey = "a9ag7SRaRxIOaswTli0hdJjT7F4dWMVR":
// Konstanta SecretKey adalah kunci rahasia yang digunakan untuk menandatangani token JWT.
const SecretKey = "a9ag7SRaRxIOaswTli0hdJjT7F4dWMVR"

// Fungsi ini digunakan untuk menghasilkan token JWT. Parameter tok adalah pointer ke struktur PayloadToken yang akan digunakan sebagai payload token. Pada fungsi ini, token diatur waktu kedaluwarsanya dan dihasilkan dengan menggunakan kunci rahasia SecretKey yang telah ditentukan sebelumnya.
func GenerateToken(tok *PayloadToken) (string, error) {

	tok.Expired = time.Now().Add(10 * 60 * time.Second)

	claims := jwt.MapClaims{
		"payload": tok,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//Fungsi ini digunakan untuk memvalidasi token JWT. Parameter tokString adalah string token JWT yang akan divalidasi. Pada fungsi ini, token JWT diparsing, kemudian dilakukan validasi tanda tangan dan isi token. Jika token valid, payload token diekstraksi dari klaim token dan dikembalikan dalam bentuk pointer ke struktur PayloadToken.

func ValidateToken(tokString string) (*PayloadToken, error) {
	tok, err := jwt.Parse(tokString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("unauthorized")
	}

	payload := claims["payload"]

	var payloadToken = PayloadToken{}
	payloadByte, _ := json.Marshal(payload)
	err = json.Unmarshal(payloadByte, &payloadToken)
	if err != nil {
		return nil, err
	}

	return &payloadToken, nil
}

//Dengan menggunakan package token ini, kita dapat menghasilkan dan memvalidasi token JWT yang menggunakan kunci rahasia tertentu.
