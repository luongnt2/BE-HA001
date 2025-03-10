package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	"math"
	"net/http"
)

type IPGeoResponse struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371

	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func GetLocationFromIP(ctx context.Context, ip string) (float64, float64, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var data IPGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	return data.Lat, data.Lon, nil
}

func GetCityCoordinates(city string) (float64, float64, error) {
	geocoder := openstreetmap.Geocoder()
	location, err := geocoder.Geocode(city)
	if err != nil {
		return 0, 0, err
	}
	if location == nil {
		return 0, 0, errors.New("city not found")
	}
	return location.Lat, location.Lng, nil
}
