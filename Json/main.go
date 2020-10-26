package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

func main() {

	url := "http://localhost:8080/mahasiswa/1811082009"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	mahasiswa := mahasiswa{}
	jsonErr := json.Unmarshal(body, &mahasiswa)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println("ID Mahasiswa :", mahasiswa.Id_Mahasiswa)
	fmt.Println("Nama :", mahasiswa.Nama)
	fmt.Println("Jurusan :", mahasiswa.Jurusan)
	fmt.Println("Program Studi :", mahasiswa.Prodi)
	fmt.Println("Alamat :")

	fmt.Println("Jalan : ", mahasiswa.Alamat.Jalan)
	fmt.Println("Kelurahan : ", mahasiswa.Alamat.Kelurahan)
	fmt.Println("Kecamatan : ", mahasiswa.Alamat.Kecamatan)
	fmt.Println("Kota : ", mahasiswa.Alamat.Kota)
	fmt.Println("Provinsi : ", mahasiswa.Alamat.Provinsi)

	for _, nilai := range mahasiswa.Nilai {
		fmt.Println("Nama Matkul : ", nilai.Mata_kuliah)
		fmt.Println("Nilai : ", nilai.Nilai)
		fmt.Println("Semester : ", nilai.Semester)

	}

}
