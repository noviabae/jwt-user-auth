package main

//Baris-baris import mengimpor paket-paket yang diperlukan dalam program, seperti paket config, controller, response, token, dan paket-paket lain yang dibutuhkan.
import (
	"backend-api/config"
	"backend-api/controller"
	"backend-api/response"
	"backend-api/token"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {

	//Menghubungkan aplikasi ke database dengan menggunakan fungsi ConnectDB() yang terdapat dalam paket config. Koneksi database disimpan dalam variabel db.
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	//Membuat instance baru dari router Gin.
	router := gin.New()

	//Menambahkan logger Gin sebagai middleware untuk mencatat setiap permintaan yang masuk.
	router.Use(gin.Logger(), CROS())

	//Membuat instance dari AuthController yang ada di dalam paket controller, dengan menyediakan objek db yang merupakan koneksi database.
	authController := controller.AuthController{
		Db: db,
	}

	//Membuat grup v1 pada router, yang akan digunakan sebagai basis URL untuk versi API.
	v1 := router.Group("v1")

	//Menambahkan endpoint GET "/ping" yang akan mengeksekusi fungsi Ping ketika permintaan GET diterima.
	router.GET("/ping", Ping)

	//Membuat grup "auth" sebagai sub-grup dari grup v1.
	auth := v1.Group("auth")
	{
		//Menambahkan endpoint POST "/auth/register" yang akan mengeksekusi fungsi Register dalam AuthController ketika permintaan POST diterima.
		auth.POST("register", authController.Register)
		//Menambahkan endpoint POST "/auth/login" yang akan mengeksekusi fungsi Login dalam AuthController ketika permintaan POST diterima.
		auth.POST("login", authController.Login)
	}

	router.Run(":8089")
}

func CROS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Request-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		ctx.Header("Access-Control-Request-Headers", "Autgorization, Content-Type")
		ctx.Next()
	}
}

//Fungsi ini adalah handler untuk endpoint "/ping". Ketika permintaan GET diterima, fungsi ini akan mengembalikan respons JSON dengan pesan "OK".

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}

// middleware

//Fungsi CheckAuth adalah sebuah middleware yang akan digunakan untuk memverifikasi token autentikasi pada setiap permintaan API. Fungsi ini menerima konteks Gin dan melakukan langkah-langkah verifikasi token:

// Mengambil nilai header "Authorization" dari permintaan.
// Memisahkan token dari string "bearer <token>".
// Jika token tidak memiliki format yang benar, mengembalikan respons JSON dengan status HTTP 401 Unauthorized.
// Memvalidasi token menggunakan fungsi ValidateToken dari paket token. Jika token tidak valid, mengembalikan respons JSON dengan status HTTP 401 Unauthorized.
// Jika token valid, menyimpan AuthId dari payload token ke konteks Gin menggunakan ctx.Set("authId", payload.AuthId).
// Melanjutkan penanganan permintaan dengan menjalankan ctx.Next().
func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "bearer ")
		fmt.Println(bearerToken)
		if len(bearerToken) != 2 {
			resp := response.ResponseApi{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
			}
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		payload, err := token.ValidateToken(bearerToken[1])
		if err != nil {
			resp := response.ResponseApi{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid Token",
				Payload:    err.Error(),
			}
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		ctx.Set("authId", payload.AuthId)

		ctx.Next()
	}
}

//Dengan konfigurasi ini, server backend-api akan berjalan dan menerima permintaan pada port 8089. Endpoint /ping dapat diakses dengan metode GET, sedangkan endpoint /auth/register dan /auth/login dapat diakses dengan metode POST setelah melalui middleware CheckAuth untuk memverifikasi token autentikasi.
