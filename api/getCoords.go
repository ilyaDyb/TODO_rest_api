package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type NominatimResponses struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func GetCoordinates(place string) (float64, float64, error) {
	encodedPlace := url.QueryEscape(place)
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", encodedPlace)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var nominatimResponses []NominatimResponses
	if err := json.NewDecoder(resp.Body).Decode(&nominatimResponses); err != nil {
		return 0, 0, err
	}
	if len(nominatimResponses) == 0 {
		return 0, 0, fmt.Errorf("no results found for place: %s", place)
	}
	lat, err := strconv.ParseFloat(nominatimResponses[0].Lat, 64)
	if err != nil {
		return 0, 0, err
	}
	lon, err := strconv.ParseFloat(nominatimResponses[0].Lon, 64)
	if err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}
