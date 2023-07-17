package config

//Kode ini mengimpor paket-paket yang diperlukan untuk menghubungkan dan berinteraksi dengan database PostgreSQL. database/sql adalah paket standar Golang untuk mengakses database, sedangkan github.com/lib/pq adalah driver database PostgreSQL untuk Go.
import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

//Fungsi ConnectDB() digunakan untuk melakukan koneksi ke database PostgreSQL. Pada bagian pertama fungsi, terdapat variabel dsn yang berisi konfigurasi koneksi ke database. Konfigurasi ini termasuk host, port, nama pengguna (user), kata sandi (password), nama database, dan opsi sslmode. Dalam kasus ini, sslmode diatur ke "disable" yang berarti koneksi tidak menggunakan SSL.

//Selanjutnya, menggunakan sql.Open("postgres", dsn), fungsi ini membuka koneksi ke database PostgreSQL dengan menggunakan driver "postgres" dan mengembalikan objek *sql.DB yang mewakili koneksi ke database.

//Setelah itu, fungsi melakukan ping ke database menggunakan db.Ping() untuk memastikan bahwa koneksi berhasil. Jika ping gagal, maka akan dikembalikan error.

// Jika semua langkah di atas berhasil, fungsi akan mengembalikan objek *sql.DB yang siap digunakan untuk berinteraksi dengan database PostgreSQL, serta nilai error yang nil (jika tidak ada error).

func ConnectDB() (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"@Syabaniyah99",
		"godb",
	)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
