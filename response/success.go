package response

type ResponseApi struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload,omitempty"`
}

//Kode yang diberikan adalah sebuah struktur (struct) dalam bahasa pemrograman Go (Golang) yang digunakan untuk membentuk respons API. Struktur ini memiliki tiga bidang (fields), yaitu StatusCode, Message, dan Payload.

//Berikut adalah penjelasan mengenai kode tersebut:

//type ResponseApi struct:
//Ini adalah deklarasi struktur dengan nama ResponseApi.

//StatusCode int:
//Bidang StatusCode bertipe data int dan digunakan untuk menyimpan kode status HTTP dari respons API. Biasanya, kode status ini akan diberikan sesuai dengan standar HTTP, seperti 200 untuk sukses, 400 untuk kesalahan permintaan klien, 500 untuk kesalahan server, dan sebagainya.

//Message string:
//Bidang Message bertipe data string dan digunakan untuk menyimpan pesan yang akan dikirimkan sebagai bagian dari respons API. Pesan ini dapat berisi informasi tentang status operasi atau kesalahan yang terjadi.

//Payload interface{}:
//Bidang Payload bertipe data interface{} dan digunakan untuk menyimpan data tambahan yang ingin disertakan dalam respons API. Bidang ini menggunakan tipe interface{} agar dapat menampung berbagai jenis data. Penggunaan omitempty pada tag json:payload,omitempty" mengindikasikan bahwa bidang ini akan diabaikan (tidak disertakan dalam respons) jika nil atau kosong.

//Dengan menggunakan struktur ResponseApi ini, kita dapat membuat respons API yang konsisten dengan menyertakan kode status, pesan, dan payload (jika diperlukan) dalam setiap respons yang dikirimkan.
