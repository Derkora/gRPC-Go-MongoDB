# UTS_5027221021_Steven Figo
Mini project Integrasi Sistem menggunakan gRPC dan Protobuf dalam bahasa pemrograman Go, serta menggunakan MongoDB

![Screenshot 2024-05-06 232846](https://github.com/Derkora/UTS_5027221021_Steven-Figo/assets/110652010/31104f0f-5397-438b-8848-9a950bc73f7a)

## link youtube
https://youtu.be/PXX0ncV-zVo 

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
