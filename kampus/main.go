package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

var db *sql.DB
var err error

type yamlconfig struct {
	Connection struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		User     string `yaml:"user"`
		Database string `yaml:"database"`
	}
}
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
	w.Header().Set("Content-Type", "application/json")
	var mahasiswa mahasiswa
	var nilaimhs nilai
	params := mux.Vars(r)

	sql := `select
				mahasiswa.id_mhs,
				mahasiswa.nama,
				jurusan.nama as jurusan,
				prodi.nama as prodi,
				mahasiswa.jalan,
				mahasiswa.kelurahan,
				mahasiswa.kecamatan,
				mahasiswa.kota,
				mahasiswa.provinsi 
				FROM
				mahasiswa.mahasiswa
				INNER JOIN mahasiswa.jurusan
				ON (mahasiswa.Id_Jurusan = jakultas.id_jurusan)
				INNER JOIN mahasiswa.prodi
				ON (mahasiswa.Id_prodi = prodi.id_prodi) where mahasiswa.id_mhs=?`
	result, err := db.Query(sql, params["id"])

	defer result.Close()
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err := result.Scan(&mahasiswa.Id_Mahasiswa, &mahasiswa.Nama, &mahasiswa.Jurusan, &mahasiswa.Prodi,
			&mahasiswa.Alamat.Jalan, &mahasiswa.Alamat.Kelurahan, &mahasiswa.Alamat.Kecamatan, &mahasiswa.Alamat.Kota, &mahasiswa.Alamat.Provinsi)

		Idmahasiswa := &mahasiswa.Id_Mahasiswa

		if err != nil {
			panic(err.Error())
		}

		sqlnilai := `SELECT
						matkul.nama,nilai.nilai,nilai.semester
						FROM
							mahasiswa.nilai
							INNER JOIN mahasiswa.matkul
								ON (nilai.Id_matkul = matkul.id_matkul) where nilai.id_mhs=?;`

		resultnilai, errnilai := db.Query(sqlnilai, *Idmahasiswa)

		defer resultnilai.Close()

		if errnilai != nil {
			panic(err.Error())
		}

		for resultnilai.Next() {
			err := resultnilai.Scan(&nilaimhs.Mata_kuliah, &nilaimhs.Nilai, &nilaimhs.Semester)

			if err != nil {
				panic(err.Error())
			}

			mahasiswa.Nilai = append(mahasiswa.Nilai, nilaimhs)
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mahasiswa)
}
func getNilai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mhsP []mahasiswa

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
		var mhs mahasiswa
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
					WHERE nilai.id_mhs = ?`

		resultNilai, errNilai := db.Query(sqlNilai, mhs.Id_Mahasiswa)

		defer resultNilai.Close()

		if errNilai != nil {
			panic(err.Error())
		}

		for resultNilai.Next() {
			var nilaiP nilai
			err := resultNilai.Scan(&nilaiP.Id_Mahasiswa, &nilaiP.Id_Matkul, &nilaiP.Mata_kuliah, &nilaiP.Nilai, &nilaiP.Semester)
			if err != nil {
				panic(err.Error())
			}
			mhs.Nilai = append(mhs.Nilai, nilaiP)
		}
		mhsP = append(mhsP, mhs)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mhsP)
}

func getNilaiAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mhsG []mahasiswa

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
			FROM mahasiswa`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var mhs2 mahasiswa
		err := result.Scan(&mhs2.Id_Mahasiswa, &mhs2.Nama, &mhs2.Alamat.Jalan, &mhs2.Alamat.Kelurahan, &mhs2.Alamat.Kecamatan, &mhs2.Alamat.Kota, &mhs2.Alamat.Provinsi, &mhs2.Jurusan, &mhs2.Prodi)

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
					WHERE nilai.id_mhs = ?`

		resultNilai, errNilai := db.Query(sqlNilai, mhs2.Id_Mahasiswa)

		defer resultNilai.Close()

		if errNilai != nil {
			panic(err.Error())
		}

		for resultNilai.Next() {
			var nilaiG nilai
			err := resultNilai.Scan(&nilaiG.Id_Mahasiswa, &nilaiG.Id_Matkul, &nilaiG.Mata_kuliah, &nilaiG.Nilai, &nilaiG.Semester)
			if err != nil {
				panic(err.Error())
			}
			mhs2.Nilai = append(mhs2.Nilai, nilaiG)
		}
		mhsG = append(mhsG, mhs2)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mhsG)
}
func updateMahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newNama := r.FormValue("nama")
		newJalan := r.FormValue("jalan")
		newKelurahan := r.FormValue("kelurahan")
		newKecamatan := r.FormValue("kecamatan")
		newKota := r.FormValue("kota")
		newProvinsi := r.FormValue("provinsi")
		newJurusan := r.FormValue("jurusan")
		newProdi := r.FormValue("prodi")

		stmt, err := db.Prepare("UPDATE mahasiswa SET nama = ?, jalan = ?, kelurahan = ?, kecamatan = ?, kota = ?, provinsi = ?, jurusan = ?, prodi = ? WHERE id_mhs = ?")

		_, err = stmt.Exec(newNama, newJalan, newKelurahan, newKecamatan, newKota, newProvinsi, newJurusan, newProdi, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Mahasiswa with id_mhs = %s was updated", params["id"])
	}
}
func createMahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		Idmahasiswa := r.FormValue("id_mhs")
		Nama := r.FormValue("nama")
		Jalan := r.FormValue("jalan")
		Kelurahan := r.FormValue("kelurahan")
		Kecamatan := r.FormValue("kecamatan")
		Kota := r.FormValue("kota")
		Provinsi := r.FormValue("provinsi")
		Jurusan := r.FormValue("jurusan")
		Prodi := r.FormValue("prodi")

		stmt, err := db.Prepare("INSERT INTO mahasiswa (id_mhs, nama, jalan, kelurahan, kecamatan, kota, provinsi, jurusan, prodi) VALUES (?,?,?,?,?,?,?,?,?)")

		_, err = stmt.Exec(Idmahasiswa, Nama, Jalan, Kelurahan, Kecamatan, Kota, Provinsi, Jurusan, Prodi)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}

// func delMahasiswa(w http.ResponseWriter, r *http.Request) {

// 	Idmahasiswa := r.FormValue("id_mahasiswa")
// 	Nama := r.FormValue("nama")

// 	stmt, err := db.Prepare("DELETE FROM mahasiswa WHERE id_mahasiswa = ? AND nama = ?")

// 	_, err = stmt.Exec(Idmahasiswa, Nama)

// 	if err != nil {
// 		fmt.Fprintf(w, "delete failed")
// 	}

// 	fmt.Fprintf(w, "Mahasiswa with ID = %s was deleted", Idmahasiswa)
// }

// Main function

func main() {
	yamlFile, err := ioutil.ReadFile("../Yaml/config.yml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}
	var yamlConfig yamlconfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	host := yamlConfig.Connection.Host
	port := yamlConfig.Connection.Port
	user := yamlConfig.Connection.User
	pass := yamlConfig.Connection.Password
	data := yamlConfig.Connection.Database

	var (
		mySQL = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, pass, host, port, data)
	)

	db, err = sql.Open("mysql", mySQL)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/mahasiswa/{id}", getNilai).Methods("GET")
	r.HandleFunc("/mahasiswaG", getNilaiAll).Methods("GET")
	r.HandleFunc("/mahasiswa/{id}", updateMahasiswa).Methods("PUT")
	r.HandleFunc("/mahasiswaT", createMahasiswa).Methods("POST")
	r.HandleFunc("/mahasiswa", getMahasiswa).Methods("GET")
	//r.HandleFunc("/delmahasiswa", delMahasiswa).Methods("POST")

	fmt.Println("Server on :8181")
	// Start server
	log.Fatal(http.ListenAndServe(":8181", r))
}
