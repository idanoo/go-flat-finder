package flatfinder

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type GoogleMapsDistanceMatrixResponse struct {
	DestinationAddresses []string `json:"destination_addresses"`
	OriginAddresses      []string `json:"origin_addresses"`
	Rows                 []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

// getDistanceFromAddress - Return distance between 2 points
func (c *LocalConfig) getDistanceFromAddress(address string, toLat float64, toLong float64) string {

	mapsURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/distancematrix/json?units=metric&mode=%s&origins=%f,%f&destinations=%s&key=%s",
		"walking",
		toLat,
		toLong,
		url.QueryEscape(address),
		c.GoogleApiToken,
	)

	client := http.Client{}
	req, err := http.NewRequest("GET", mapsURL, nil)
	if err != nil {
		log.Print(err)
		return "UNKNOWN"
	}

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return "UNKNOWN"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return "UNKNOWN"
		}

		// Decode JSON
		var mapsResult GoogleMapsDistanceMatrixResponse
		err = json.Unmarshal(bodyBytes, &mapsResult)
		if err != nil {
			log.Print(err)
			return "UNKNOWN"
		}

		dist := "N/A"
		time := "N/A"
		for _, rows := range mapsResult.Rows {
			for _, element := range rows.Elements {
				dist = element.Distance.Text
				time = element.Duration.Text
			}
		}

		return fmt.Sprintf("%s (%s)", dist, time)
	} else {
		log.Printf("Maps API error: %s", resp.Status)
	}

	return "UNKNOWN"
}
