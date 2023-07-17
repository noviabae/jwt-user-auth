package controller

//Kode ini mengimpor paket-paket yang diperlukan untuk mengimplementasikan fungsi-fungsi kontroler, seperti backend-api/response untuk menghasilkan respons API, backend-api/token untuk menghasilkan dan memverifikasi token, database/sql untuk berinteraksi dengan database, net/http untuk mengatur status HTTP, dan github.com/gin-gonic/gin untuk menggunakan framework web Gin dan github.com/go-playground/validator/v10 untuk melakukan validasi struktur data.
import (
	"backend-api/response"
	"backend-api/token"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// Struktur AuthController memiliki satu bidang yaitu Db yang merupakan objek *sql.DB untuk berinteraksi dengan database.
type AuthController struct {
	Db *sql.DB
}

// Struktur RegisterRequest digunakan untuk menguraikan data yang diterima saat pendaftaran pengguna baru. Struktur LoginRequest digunakan untuk menguraikan data yang diterima saat proses login. Struktur Auth digunakan untuk menyimpan data pengguna yang ditemukan dalam database.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password"`
}

type Auth struct {
	Id       int
	Email    string
	Password string
}

// Variabel queryCreate berisi pernyataan SQL untuk memasukkan data pengguna baru ke dalam tabel auth. Variabel queryFindByEmail berisi pernyataan SQL untuk mencari pengguna berdasarkan alamat email.
var (
	queryCreate = `	INSERT INTO auth (name, email, password)
		VALUES ($1, $2, $3)`

	queryFindByEmail = `
		SELECT id, email, password
		FROM auth
		WHERE email=$1`
)

// Fungsi ini digunakan untuk menangani permintaan pendaftaran pengguna baru. Data permintaan yang diterima dari ctx (konteks Gin) diuraikan ke dalam variabel req yang merupakan instansi dari struktur RegisterRequest.
func (a *AuthController) Register(ctx *gin.Context) {

	var req = RegisterRequest{}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Validator digunakan untuk memvalidasi struktur req berdasarkan tag validasi yang didefinisikan pada struktur tersebut. Jika validasi gagal, maka akan dikembalikan respons HTTP dengan status BadRequest dan pesan kesalahan validasi.
	val := validator.New()
	err = val.Struct(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//cara agar password tidak terlihat (jadi di enkripsi pakai kriptografi)
	//byte: jumlah berapa kali dia di enkripsi
	//Password yang diterima dalam req dienkripsi menggunakan fungsi bcrypt.GenerateFromPassword() sebelum disimpan ke database.

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.Password = string(hash)

	//Pernyataan SQL disiapkan menggunakan objek *sql.DB (a.Db) dan pernyataan queryCreate untuk memasukkan data pengguna baru ke dalam tabel auth.
	statement, err := a.Db.Prepare(queryCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Pernyataan SQL dieksekusi dengan mengisi parameter dengan nilai yang diperlukan dari req.
	_, err = statement.Exec(
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Jika pendaftaran pengguna baru berhasil, respons JSON dengan status HTTP 201 Created dan pesan sukses dikirimkan ke klien.
	resp := response.ResponseApi{
		StatusCode: http.StatusCreated,
		Message:    "CREATED SUCCESS",
	}
	ctx.JSON(resp.StatusCode, resp)
}

// Fungsi ini digunakan untuk menangani permintaan login pengguna. Data permintaan yang diterima dari ctx (konteks Gin) diuraikan ke dalam variabel req yang merupakan instansi dari struktur LoginRequest.
func (a *AuthController) Login(ctx *gin.Context) {
	var req = LoginRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	statement, err := a.Db.Prepare(queryFindByEmail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//kalau data banyak pakai query, kalau data cuma satu pakai query row
	//Pernyataan SQL dieksekusi menggunakan QueryRow() untuk mendapatkan satu baris hasil data pengguna berdasarkan alamat email yang diberikan dalam req. Hasil tersebut di-scan ke dalam variabel auth yang merupakan instansi dari struktur Auth.
	row := statement.QueryRow(req.Email)

	var auth = Auth{}

	err = row.Scan(
		&auth.Id,
		&auth.Email,
		&auth.Password,
	)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Password yang diterima dalam req dibandingkan dengan password yang dienkripsi yang ditemukan dalam auth. Jika tidak cocok, maka akan dikirimkan respons HTTP dengan status NotFound dan pesan kesalahan.
	err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(req.Password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	tok := token.PayloadToken{
		AuthId: auth.Id,
	}

	//Token akses di-generate menggunakan objek tok yang berisi informasi identitas pengguna, kemudian diubah menjadi string.
	tokString, err := token.GenerateToken(&tok)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Jika proses login berhasil, respons JSON dengan status HTTP 200 OK, pesan sukses, dan token akses dikirimkan ke klien.
	resp := response.ResponseApi{
		StatusCode: http.StatusOK,
		Message:    "LOGIN SUCCESS",
		Payload: gin.H{
			"token": tokString,
		},
	}
	ctx.JSON(resp.StatusCode, resp)
}

//Dengan menggunakan kontroler ini, kita dapat mengimplementasikan endpoint API untuk pendaftaran pengguna baru dan proses login dengan menggunakan Gin dan berinteraksi dengan database PostgreSQL yang dihubungkan melalui objek *sql.DB.
