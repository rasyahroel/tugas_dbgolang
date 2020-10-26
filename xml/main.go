package main

import (
	"database/sql"
	"encoding/xml"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type mahasiswa struct {
	Id_Mahasiswa string `json:"id_mhs"`
	Nama         string `json:"nama"`
	Alamat       struct {
		Jalan     string `json:"jalan"`
		Kelurahan string `json:"kelurahan"`
		Kecamatan string `json:"kecamatan"`
		Kota      string `json:"kota"`
		Provinsi  string `json:"provinsi"`
	} `json:"alamat"`
	Jurusan string  `json:"jurusan"`
	Prodi   string  `json:"prodi"`
	Nilai   []nilai `json:"nilai"`
}

type nilai struct {
	Id_Mahasiswa string  `json:"id_mhs"`
	Id_Matkul    string  `json:"id_matkul"`
	Mata_kuliah  string  `json:"matkul"`
	Nilai        float32 `json:"nilai"`
	Semester     int8    `json:"semester"`
}

func getMahasiswa(w http.ResponseWriter, r *http.Request) {

	var mhs mahasiswa
	var ni nilai
	params := mux.Vars(r)

	sql := `SELECT
				id_mhs,
				IFNULL(nama,'') nama,
				IFNULL(jalan,'') jalan,
				IFNULL(kelurahan,'') kelurahan,
				IFNULL(kecamatan,'') kecamatan,
				IFNULL(kota,'') kota,
				IFNULL(provinsi,'') provinsi,
				IFNULL(jurusan,'') jurusan,
				IFNULL(prodi,'') prodi				
				FROM mahasiswa WHERE id_mhs IN (?)`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		err := result.Scan(&mhs.Id_Mahasiswa, &mhs.Nama, &mhs.Alamat.Jalan, &mhs.Alamat.Kelurahan, &mhs.Alamat.Kecamatan, &mhs.Alamat.Kota, &mhs.Alamat.Provinsi, &mhs.Jurusan, &mhs.Prodi)

		if err != nil {
			panic(err.Error())
		}

		sqlNilai := `SELECT
						id_mhs		
						, matakuliah.id_matkul
						, matakuliah.matkul
						, nilai
						, semester
					FROM
						nilai INNER JOIN matakuliah
						ON (nilai.id_matkul = matakuliah.id_matkul)
						WHERE id_mhs = ?`

		id_mahasiswa := &mhs.Id_Mahasiswa

		resultDetail, errDet := db.Query(sqlNilai, *id_mahasiswa)

		defer resultDetail.Close()

		if errDet != nil {
			panic(err.Error())
		}

		for resultDetail.Next() {

			err := resultDetail.Scan(&ni.Id_Mahasiswa, &ni.Id_Matkul, &ni.Mata_kuliah, &ni.Nilai, &ni.Semester)

			if err != nil {
				panic(err.Error())
			}

			mhs.Nilai = append(mhs.Nilai, ni)

		}

	}
	//header for exml
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))
	xml.NewEncoder(w).Encode(mhs)
}
func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/mahasiswa/{id}", getMahasiswa).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
