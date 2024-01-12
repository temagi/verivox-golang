package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestAddressCheckingService(t *testing.T) {
	baseURL := "https://service.verivox.de/geo/latestv2/cities"

	t.Run("Valid German postcodes", func(t *testing.T) {
		validPostcodes := []struct {
			Postcode int
			Cities   []string
		}{
			{Postcode: 10409, Cities: []string{"Berlin"}},
			{Postcode: 77716, Cities: []string{"Fischerbach", "Haslach", "Hofstetten"}},
		}

		for _, code := range validPostcodes {
			t.Run(fmt.Sprintf("should return cities for postcode %d", code.Postcode), func(t *testing.T) {
				uri := fmt.Sprintf("%d", code.Postcode)
				response := makeAPIRequest(baseURL, uri)

				require.Equal(t, fasthttp.StatusOK, response.StatusCode())
				type CitiesResponse struct {
					Cities []string `json:"Cities"`
				}

				var result CitiesResponse
    			json.Unmarshal([]byte(response.Body()), &result)
				assert.Equal(t, code.Cities, result.Cities)
			})
		}
	})

	t.Run("Invalid German postcode", func(t *testing.T) {
		invalidPostcode := 22333
		t.Run(fmt.Sprintf("should return HTTP 404 for invalid postcode %d", invalidPostcode), func(t *testing.T) {
			uri := fmt.Sprintf("%d", invalidPostcode)
			response := makeAPIRequest(baseURL, uri)

			require.Equal(t, fasthttp.StatusNotFound, response.StatusCode())
			assert.Empty(t, response.Body())
		})
	})

	t.Run("Find the streets for a given postcode", func(t *testing.T) {
		streets := []struct {
			Postcode int
			City     string
			Streets  []string
		}{
			{Postcode: 10409, City: "Berlin", Streets: readDataFile("Berlin")},
			{Postcode: 77716, City: "Fischerbach", Streets: readDataFile("Fischerbach")},
			{Postcode: 77716, City: "Haslach", Streets: readDataFile("Haslach")},
			{Postcode: 77716, City: "Hofstetten", Streets: readDataFile("Hofstetten")},
		}

		for _, street := range streets {
			t.Run(fmt.Sprintf("should return streets for postcode %d and city %s", street.Postcode, street.City), func(t *testing.T) {
				uri := fmt.Sprintf("%d/%s/streets", street.Postcode, street.City)
				response := makeAPIRequest(baseURL, uri)

				type StreetsResponse struct {
					Streets []string `json:"Streets"`
				}
				var result StreetsResponse
    			json.Unmarshal([]byte(response.Body()), &result)

				require.Equal(t, fasthttp.StatusOK, response.StatusCode())
				assert.Equal(t, street.Streets, result.Streets)
			})
		}
	})
}

func makeAPIRequest(baseURL, uri string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/%s", baseURL, uri))

	resp := fasthttp.AcquireResponse()
	fasthttp.Do(req, resp)

	return resp
}

func readDataFile(filename string) []string{
	data, err := os.ReadFile("testdata/" + filename + ".txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File does not exist")  
		} else if errors.Is(err, os.ErrPermission) {
			fmt.Println("Permission denied")
		} else {
			fmt.Printf("Unhandled error %v occurred\n", err)
			panic(err)
		}
	}
	return strings.Split(string(data), "\n")
}
