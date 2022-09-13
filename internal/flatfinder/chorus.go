package flatfinder

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ChorusAddressLookupResponse struct {
}

func chorusAddressLookup(address string) string {
	log.Printf("Querying address: %s", address)
	lookupURL := fmt.Sprintf(
		"https://api.chorus.co.nz/addresslookup/v1/addresses/?fuzzy=true&q=%s",
		url.QueryEscape(address),
	)
	//curl 'https://api.chorus.co.nz/addresslookup/v1/addresses/?fuzzy=true&q=35%20Rosalind' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0' -H 'Accept: application/json, text/plain, */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'X-Chorus-Client-Id: 82d4b4a8050c4d5e97c5f06120ef9c04' -H 'X-Chorus-Client-Secret: 8899c64746474Cf18849c6B721b5Db51' -H 'X-Transaction-Id: ca5ef871-9b3f-4af4-958d-ee9c51094e08' -H 'Origin: https://www.chorus.co.nz'
	// Build HTTP request
	client := http.Client{}
	req, err := http.NewRequest("GET", lookupURL, nil)
	if err != nil {
		log.Print(err)
		return "UNK"
	}

	// Magic numbers - May need to dynamically receive these
	req.Header.Set("X-Chorus-Client-Id", "82d4b4a8050c4d5e97c5f06120ef9c04")
	req.Header.Set("X-Chorus-Client-Secret", "8899c64746474Cf18849c6B721b5Db51")
	req.Header.Set("X-Transaction-Id", "ca5ef871-9b3f-4af4-958d-ee9c51094e08")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "https://tinker.nz/idanoo/flat-finder")

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return "UNK"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNonAuthoritativeInfo {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return "UNK"
		}

		log.Println(string(bodyBytes))
		return "N/A"
	} else {
		log.Print("Invalid response from API: " + resp.Status)
	}

	return "UNK"
}

// getAvailableSpeeds - Checks if VDSL/FIBRE is available
func getAvailableSpeeds(address string) string {
	return chorusAddressLookup(address)

	// // Build HTTP request
	// client := http.Client{}
	// req, err := http.NewRequest("GET", chorusURL, nil)
	// if err != nil {
	// 	log.Print(err)
	// 	return "UNK"
	// }

	// req.Header.Set("Content-TypeContent-Type", "application/json")
	// req.Header.Set("User-Agent", "https://tinker.nz/idanoo/flat-finder")

	// // Do the request
	// resp, err := client.Do(req)
	// if err != nil {
	// 	log.Print(err)
	// 	return "UNK"
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode == http.StatusOK {
	// 	bodyBytes, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Print(err)
	// 		return "UNK"
	// 	}

	// 	log.Println(string(bodyBytes))

	// 	return "N/A"
	// } else {
	// 	log.Print("Invalid response from API: " + resp.Status)
	// }

	return "UNK"
}
