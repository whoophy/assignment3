package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"time"
)

const (
	fileName = "data.json"
)

type Variable struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Status struct {
	Status Variable `json:"status"`
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteJson() {
	var status Status
	water, _ := rand.Int(rand.Reader, big.NewInt(100))
	wind, _ := rand.Int(rand.Reader, big.NewInt(100))
	/*
		karena random dengan crypto/rand merupakan *big.Int
		maka perlu di rubah ke int
	*/
	status.Status.Water = int(water.Int64())
	status.Status.Wind = int(wind.Int64())
	statusByte, err := json.Marshal(status)
	panicErr(err)
	_ = ioutil.WriteFile(fileName, statusByte, 0644)
}

func OpenJson() *os.File {
	data, err := os.Open(fileName)
	panicErr(err)
	return data
}

func Reload(w http.ResponseWriter, r *http.Request) {
	var status Status
	data := OpenJson()
	defer data.Close()

	go func() {
		for {
			time.Sleep(time.Second * 15)
			WriteJson()
		}
	}()

	JsonSting, err := ioutil.ReadAll(data)
	panicErr(err)

	json.Unmarshal([]byte(JsonSting), &status)

	var res = make(map[string]interface{})
	res["Water"] = status.Status.Water
	res["Wind"] = status.Status.Wind

	switch {
	case status.Status.Water <= 5:
		res["WaterStatus"] = "Aman"
	case status.Status.Water > 5 && status.Status.Water < 9:
		res["WaterStatus"] = "Siaga"
	default:
		res["WaterStatus"] = "Bahaya"
	}

	switch {
	case status.Status.Wind <= 6:
		res["WindStatus"] = "Aman"
	case status.Status.Wind > 6 && status.Status.Wind < 16:
		res["WindStatus"] = "Siaga"
	default:
		res["WindStatus"] = "Bahaya"
	}

	template, err := template.ParseFiles("index.html")
	panicErr(err)
	template.Execute(w, res)
}

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", Reload)
	fmt.Println("server running on port 8080")
	http.ListenAndServe(":8080", mux)
}
