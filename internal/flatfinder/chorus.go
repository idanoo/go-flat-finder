package flatfinder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ChorusAddressSearchResponse struct {
	Results []struct {
		Aid   string `json:"aid"`
		Label string `json:"label"`
		Links []struct {
			Rel    string `json:"rel"`
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"links"`
	} `json:"results"`
}

type ChorusUniqueIdResponse struct {
	FormattedAddress struct {
		Line1 string      `json:"line1"`
		Line2 interface{} `json:"line2"`
		Line3 string      `json:"line3"`
		Line4 string      `json:"line4"`
	} `json:"formattedAddress"`
	StructuredAddress struct {
		LevelNumber     interface{} `json:"levelNumber"`
		LevelType       interface{} `json:"levelType"`
		StreetNumber    int         `json:"streetNumber"`
		SituationNumber int         `json:"situationNumber"`
		Suffix          interface{} `json:"suffix"`
		Unit            interface{} `json:"unit"`
		UnitType        interface{} `json:"unitType"`
		StreetName      string      `json:"streetName"`
		RoadType        string      `json:"roadType"`
		RoadAbv         string      `json:"roadAbv"`
		RoadSuffix      interface{} `json:"roadSuffix"`
		Suburb          string      `json:"suburb"`
		RuralDelivery   interface{} `json:"ruralDelivery"`
		Town            string      `json:"town"`
		Postcode        string      `json:"postcode"`
		BoxNumber       interface{} `json:"boxNumber"`
		BoxLobby        interface{} `json:"boxLobby"`
		BoxType         interface{} `json:"boxType"`
		Region          string      `json:"region"`
		Country         string      `json:"country"`
		IsPrimary       string      `json:"isPrimary"`
	} `json:"structuredAddress"`
	Location struct {
		NztmX    float64 `json:"nztmX"`
		NztmY    float64 `json:"nztmY"`
		NzmgX    float64 `json:"nzmgX"`
		NzmgY    float64 `json:"nzmgY"`
		Wgs84Lat float64 `json:"wgs84Lat"`
		Wgs84Lon float64 `json:"wgs84Lon"`
	} `json:"location"`
	References struct {
		Aid   string      `json:"aid"`
		Dpid  interface{} `json:"dpid"`
		Tui   int         `json:"tui"`
		Tlc   int         `json:"tlc"`
		Plsam int         `json:"plsam"`
	} `json:"references"`
	Related []interface{} `json:"related"`
	Links   []struct {
		Rel    string `json:"rel"`
		Href   string `json:"href"`
		Method string `json:"method"`
	} `json:"links"`
}

type ChorusAddressLookupResponse struct {
	RegionRsp                string `json:"region_rsp"`
	SubregionRsp             string `json:"subregion_rsp"`
	AreaHyperfibre           string `json:"area_hyperfibre"`
	AlternativeFibreProvider string `json:"alternative_fibre_provider"`
	AreaFibreSupplier        string `json:"area_fibre_supplier"`
	PointOfInterconnect      string `json:"point_of_interconnect"`
	ProductZoneType          string `json:"product_zone_type"`
	ActiveServices           []struct {
		Service     string `json:"service"`
		SpeedMbps   int    `json:"speed_mbps"`
		SpeedUlMbps int    `json:"speed_ul_mbps"`
	} `json:"active_services"`
	AvailableServices []struct {
		Service              string  `json:"service"`
		ServiceIndicator     string  `json:"service_indicator"`
		Capable              string  `json:"capable"`
		SpeedMbps            float64 `json:"speed_mbps"`
		InstallLeadTimeDays  string  `json:"install_lead_time_days,omitempty"`
		InstallLeadTimeWeeks string  `json:"install_lead_time_weeks,omitempty"`
		SpeedUlMbps          float64 `json:"speed_ul_mbps,omitempty"`
	} `json:"available_services"`
	FutureServices []interface{} `json:"future_services"`
	Fibre          struct {
		BuildRequired      string      `json:"build_required"`
		ConsentRequired    string      `json:"consent_required"`
		ConsentStatus      string      `json:"consent_status"`
		DesignRequired     string      `json:"design_required"`
		DwellingType       string      `json:"dwelling_type"`
		FibreInADayCapable string      `json:"fibre_in_a_day_capable"`
		Greenfields        string      `json:"greenfields"`
		IntactOnt          string      `json:"intact_ont"`
		MduBuildStatus     string      `json:"mdu_build_status"`
		MduClass           string      `json:"mdu_class"`
		MduDesignStatus    interface{} `json:"mdu_design_status"`
		PermitDelayLikely  string      `json:"permit_delay_likely"`
		RightOfWay         string      `json:"right_of_way"`
	} `json:"fibre"`
	Copper struct {
		PremiseWiringRecommended string `json:"premise_wiring_recommended"`
	} `json:"copper"`
}

// getAvailableSpeeds - Checks if VDSL/FIBRE is available
func getAvailableSpeeds(address string) (string, string) {
	aid, err := chorusAddressLookup(address)
	if err != nil {
		log.Print(err)
		return "UNK", "UNK"
	}
	// log.Printf("Using AID: %s", aid)

	tlc, err := chorusGetUnqiueID(aid)
	if err != nil {
		log.Print(err)
		return "UNK", "UNK"
	}
	// log.Printf("Using TLC: %d", tlc)

	chorusURL := fmt.Sprintf("https://www.chorus.co.nz/api/bbc/bcc/%d", tlc)
	client := http.Client{}
	req, err := http.NewRequest("GET", chorusURL, nil)
	if err != nil {
		log.Print(err)
		return "UNK", "UNK"
	}

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return "UNK", "UNK"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return "UNK", "UNK"
		}

		// Decode JSON
		var chorusResult ChorusAddressLookupResponse
		err = json.Unmarshal(bodyBytes, &chorusResult)
		if err != nil {
			log.Print(err)
			return "UNK", "UNK"
		}

		hasFibre := "No"
		if chorusResult.Fibre.BuildRequired == "N" {
			hasFibre = "Yes"
		}

		maxSpeed := 0.0
		for _, available := range chorusResult.AvailableServices {
			if available.SpeedMbps > maxSpeed && available.Capable == "YES" {
				maxSpeed = available.SpeedMbps
			}
		}

		current := "None"
		for _, active := range chorusResult.ActiveServices {
			current = fmt.Sprintf("%s (%d Mbps)", active.Service, active.SpeedMbps)
		}

		return fmt.Sprintf("%s (%.0f Mbps)", hasFibre, maxSpeed), current
	} else {
		log.Print("Invalid response from API: " + resp.Status)
	}

	return "UNK", "UNK"
}

// chorusAddressLookup - Try get the AID for the address
func chorusAddressLookup(address string) (string, error) {
	lookupURL := fmt.Sprintf(
		"https://api.chorus.co.nz/addresslookup/v1/addresses/?fuzzy=true&q=%s",
		url.QueryEscape(address),
	)

	// Build HTTP request
	client := http.Client{}
	req, err := http.NewRequest("GET", lookupURL, nil)
	if err != nil {
		return "", err
	}

	// Magic numbers - May need to dynamically receive these
	req.Header.Set("X-Chorus-Client-Id", "82d4b4a8050c4d5e97c5f06120ef9c04")
	req.Header.Set("X-Chorus-Client-Secret", "8899c64746474Cf18849c6B721b5Db51")
	req.Header.Set("X-Transaction-Id", "ca5ef871-9b3f-4af4-958d-ee9c51094e08")
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// They return a 203
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNonAuthoritativeInfo {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// Decode JSON
		var chorusResult ChorusAddressSearchResponse
		err = json.Unmarshal(bodyBytes, &chorusResult)
		if err != nil {
			return "", err
		}

		// If we have a result, return the first one
		for _, result := range chorusResult.Results {
			return result.Aid, nil
		}

		return "", errors.New("No results found for address: " + address)
	}

	return "", errors.New("Invalid response from API: " + resp.Status)
}

// chorusGetUniqueID - Return ID needed to get avail services
func chorusGetUnqiueID(aid string) (int64, error) {
	lookupURL := fmt.Sprintf(
		"https://api.chorus.co.nz/addresslookup/v1/addresses/aid:%s",
		aid,
	)

	// Build HTTP request
	client := http.Client{}
	req, err := http.NewRequest("GET", lookupURL, nil)
	if err != nil {
		return 0, err
	}

	// Magic numbers - May need to dynamically receive these
	req.Header.Set("X-Chorus-Client-Id", "82d4b4a8050c4d5e97c5f06120ef9c04")
	req.Header.Set("X-Chorus-Client-Secret", "8899c64746474Cf18849c6B721b5Db51")
	req.Header.Set("X-Transaction-Id", "ca5ef871-9b3f-4af4-958d-ee9c51094e08")
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// They return a 203
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNonAuthoritativeInfo {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		// Decode JSON
		var chorusResult ChorusUniqueIdResponse
		err = json.Unmarshal(bodyBytes, &chorusResult)
		if err != nil {
			return 0, err
		}

		return int64(chorusResult.References.Tlc), nil
	}

	return 0, errors.New("Invalid response from API: " + resp.Status)
}
