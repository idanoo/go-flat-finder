package flatfinder

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var TradeMeBaseURL = "https://api.trademe.co.nz/v1/Search/Property/Rental.json"

// https://developer.trademe.co.nz/api-reference/search-methods/rental-search
type TrademeResultSet struct {
	TotalCount      int              `json:"TotalCount"`
	Page            int              `json:"Page"`
	PageSize        int              `json:"PageSize"`
	List            []TradeMeListing `json:"List"`
	FoundCategories []interface{}    `json:"FoundCategories"`
	SearchQueryID   string           `json:"SearchQueryId"`
}

type TradeMeListing struct {
	ListingID          int64         `json:"ListingId"`
	Title              string        `json:"Title"`
	Category           string        `json:"Category"`
	StartPrice         int           `json:"StartPrice"`
	StartDate          string        `json:"StartDate"`
	EndDate            string        `json:"EndDate"`
	ListingLength      interface{}   `json:"ListingLength"`
	IsFeatured         bool          `json:"IsFeatured,omitempty"`
	HasGallery         bool          `json:"HasGallery"`
	IsBold             bool          `json:"IsBold,omitempty"`
	IsHighlighted      bool          `json:"IsHighlighted,omitempty"`
	AsAt               string        `json:"AsAt"`
	CategoryPath       string        `json:"CategoryPath"`
	PictureHref        string        `json:"PictureHref"`
	RegionID           int           `json:"RegionId"`
	Region             string        `json:"Region"`
	SuburbID           int           `json:"SuburbId"`
	Suburb             string        `json:"Suburb"`
	NoteDate           string        `json:"NoteDate"`
	ReserveState       int           `json:"ReserveState"`
	IsClassified       bool          `json:"IsClassified"`
	OpenHomes          []interface{} `json:"OpenHomes"`
	GeographicLocation struct {
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
		Northing  int     `json:"Northing"`
		Easting   int     `json:"Easting"`
		Accuracy  int     `json:"Accuracy"`
	} `json:"GeographicLocation"`
	PriceDisplay   string   `json:"PriceDisplay"`
	PhotoUrls      []string `json:"PhotoUrls"`
	AdditionalData struct {
		BulletPoints []interface{} `json:"BulletPoints"`
		Tags         []interface{} `json:"Tags"`
	} `json:"AdditionalData"`
	ListingExtras       []string `json:"ListingExtras"`
	MemberID            int      `json:"MemberId"`
	Address             string   `json:"Address"`
	District            string   `json:"District"`
	AvailableFrom       string   `json:"AvailableFrom"`
	Bathrooms           int      `json:"Bathrooms"`
	Bedrooms            int      `json:"Bedrooms"`
	ListingGroup        string   `json:"ListingGroup"`
	Parking             string   `json:"Parking"`
	PetsOkay            int      `json:"PetsOkay"`
	PropertyType        string   `json:"PropertyType"`
	RentPerWeek         int      `json:"RentPerWeek"`
	SmokersOkay         int      `json:"SmokersOkay"`
	Whiteware           string   `json:"Whiteware"`
	AdjacentSuburbNames []string `json:"AdjacentSuburbNames"`
	AdjacentSuburbIds   []int    `json:"AdjacentSuburbIds"`
	DistrictID          int      `json:"DistrictId"`
	Agency              struct {
		ID       int    `json:"Id"`
		Name     string `json:"Name"`
		Website  string `json:"Website"`
		Logo     string `json:"Logo"`
		Branding struct {
			BackgroundColor string `json:"BackgroundColor"`
			TextColor       string `json:"TextColor"`
			StrokeColor     string `json:"StrokeColor"`
			OfficeLocation  string `json:"OfficeLocation"`
			LargeBannerURL  string `json:"LargeBannerURL"`
		} `json:"Branding"`
		Logo2  string `json:"Logo2"`
		Agents []struct {
			FullName string `json:"FullName"`
		} `json:"Agents"`
		IsRealEstateAgency bool `json:"IsRealEstateAgency"`
	} `json:"Agency,omitempty"`
	TotalParking    int    `json:"TotalParking"`
	IsSuperFeatured bool   `json:"IsSuperFeatured"`
	AgencyReference string `json:"AgencyReference"`
	BestContactTime string `json:"BestContactTime"`
	IdealTenant     string `json:"IdealTenant"`
	MaxTenants      int    `json:"MaxTenants"`
	PropertyID      string `json:"PropertyId"`
	Amenities       string `json:"Amenities"`
	Lounges         int    `json:"Lounges"`
}

func (c *LocalConfig) searchTrademe() error {
	// Only pull last 2 hours by default
	dateFrom := time.Now().Add(-time.Hour * 8)

	// Set filters
	queryParams := url.Values{}
	queryParams.Add("photo_size", "FullSize")   // 670x502
	queryParams.Add("sort_order", "Default")    // Standard order
	queryParams.Add("return_metadata", "false") // Include search data
	queryParams.Add("rows", "500")              // Total results

	queryParams.Add("date_from", dateFrom.Format("2006-01-02T15:00"))
	queryParams.Add("suburb", c.Suburbs)
	queryParams.Add("property_type", c.PropertyTypes)
	queryParams.Add("price_max", c.PriceMax)
	queryParams.Add("bedrooms_min", c.BedroomsMin)
	queryParams.Add("bedrooms_max", c.BedroomsMax)

	// Build HTTP request
	client := http.Client{}
	req, err := http.NewRequest("GET", TradeMeBaseURL, nil)
	if err != nil {
		return err
	}

	// Append our filters
	req.URL.RawQuery = queryParams.Encode()

	// Auth
	req.Header.Set("Authorization", "OAuth oauth_consumer_key=\""+c.TradeMeKey+"\", oauth_signature_method=\"PLAINTEXT\", oauth_signature=\""+c.TradeMeSecret+"&\"")

	req.Header.Set("Content-TypeContent-Type", "application/json")
	req.Header.Set("User-Agent", "https://tinker.nz/idanoo/flat-finder")

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return c.handleTrademeResponse(bodyBytes)
	} else {
		return errors.New("Invalid response from API: " + resp.Status)
	}
}

func (c *LocalConfig) handleTrademeResponse(responseJson []byte) error {
	var resultSet TrademeResultSet
	err := json.Unmarshal(responseJson, &resultSet)
	if err != nil {
		return err
	}

	log.Printf("Query complete. Listings: %d", resultSet.TotalCount)
	for _, result := range resultSet.List {
		c.parseTrademeListing(result)
	}

	// Update config if succcess
	c.storeConfig()
	return nil
}

func (c *LocalConfig) parseTrademeListing(listing TradeMeListing) {
	// Only send if we haven't before!
	if _, ok := c.PostedProperties[listing.ListingID]; !ok {
		// Send the message!
		c.sendEmbeddedMessage(listing)

		// Make sure we add the key in to the map so we don't send it again!
		c.PostedProperties[listing.ListingID] = true
	}
}
