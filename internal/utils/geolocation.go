package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type Location struct {
	DisplayName string `json:"display_name"`
}

func ReverseGeocode(lat, lon string) (string, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%s&lon=%s", lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var location Location
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return "", err
	}

	return location.DisplayName, nil
}

func ReverseGeocodeWithCity(lat, lon string) (string, string, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%s&lon=%s", lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var location struct {
		Name string `json:"name"`
		Address struct {
			City  string `json:"city"`
		} `json:"address"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return "", "", err
	}

	city := location.Address.City

	return location.Name, city, nil
}

func TestReverseGeocode(t *testing.T) {
	lat := "21.038938"
	lon := "-89.661438"
	expected := "Playa del Carmen, Quintana Roo, México" // Cambia esto según la respuesta esperada

	address, err := ReverseGeocode(lat, lon)
	if err != nil {
		t.Fatalf("Error al obtener la ubicación: %v", err)
	}

	if address != expected {
		t.Errorf("Se esperaba %s, pero se obtuvo %s", expected, address)
	}
}
