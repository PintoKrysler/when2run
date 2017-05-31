package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pintokrysler/when2run/models"
)

// MakeWeatherAPIcall ...
func MakeWeatherAPIcall(s models.Settings) models.Responsetype {
	apiKey := "4793867f02934a10b3033be4d68f385c"
	baseURL := "http://api.openweathermap.org/data/2.5/forecast?q=lakewood,co&units=imperial"
	query := baseURL + "&appid=" + apiKey + "&id=5427946"

	res, err := http.Get(query)
	if err != nil {
		log.Fatal(err)
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var r = models.Responsetype{}
	json.Unmarshal(response, &r)
	fmt.Println("compare values with settings", s)
	r = parseData(r, s)

	return r
}

// parseData
// This function parses the Weather API data
// Transforms ts into readable data for the view
func parseData(data models.Responsetype, s models.Settings) models.Responsetype {
	for i := 0; i < len(data.List); i++ {
		elem := data.List[i]
		tsString := strconv.Itoa(elem.Ts)
		tsFormatted, err := strconv.ParseInt(tsString, 10, 64)
		if err != nil {
			panic(err)
		}
		data.List[i].TsFormatted = time.Unix(tsFormatted, 0)

		if s.MinTemp <= elem.TempValues.TempMin && s.MaxTemp >= elem.TempValues.TempMax {
			data.List[i].GoRun = true
		}

	}
	return data
}
