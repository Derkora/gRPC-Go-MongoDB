# UTS_5027221021_Steven Figo
 
mini project Integrasi Sistem menggunakan gRPC dan Protobuf dalam bahasa pemrograman Go, serta menggunakan MongoDB

## link youtube


## Deskripsi Project

Project ini memiliki fungsi sebagai berikut:
- Memiliki fitur Create, Read, Update, Delete data
- Koneksi ke database (MongoDB atau yang lainnya)
- Backend CRUD ke database
- Mengimplementasikan UI

## Cara menjalankan project ini
1. Jalankan perintah `go mod tidy` untuk menambah dan menghapus modul
2. Atur `configs/default-config.yml`
3. Untuk server jalankan perintah `make run-server`
4. Untuk client jalankan perintah `make run-client` (pada terminal yang berbeda)
5. Web dapat diakses lewat `localhost:1000/list`