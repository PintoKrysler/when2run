package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//ResponseMain ...
type responseMain struct {
	Temp    float64 `json:"temp"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
}

//ResponseElem ...
type responseElem struct {
	Ts         int          `json:"dt"`
	TempValues responseMain `json:"main"`
}

// Response ...
type responsetype struct {
	List []responseElem `json:"list"`
}

func main() {
	fmt.Println("This is a test")
	apiKey := "4793867f02934a10b3033be4d68f385c"
	baseURL := "http://api.openweathermap.org/data/2.5/forecast?q=lakewood,co"
	query := baseURL + "&appid=" + apiKey + "&id=5427946"
	fmt.Println(query)
	res, err := http.Get(query)
	if err != nil {
		log.Fatal(err)
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var r = responsetype{}
	json.Unmarshal(response, &r)

	// Parse JSON
	fmt.Println(r.List[0])
}
