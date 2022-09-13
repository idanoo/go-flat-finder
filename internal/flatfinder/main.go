package flatfinder

import (
	"log"
	"time"

	"github.com/disgoorg/disgo/webhook"
)

// Our local struct we will store data during runtime
type LocalConfig struct {
	DiscordWebhook string         `json:"-"`
	DiscordClient  webhook.Client `json:"-"`

	GoogleApiToken  string `json:"-"`
	GoogleLocation1 string `json:"-"`
	GoogleLocation2 string `json:"-"`

	TradeMeKey    string `json:"-"`
	TradeMeSecret string `json:"-"`

	Suburbs       string `json:"-"`
	BedroomsMin   string `json:"-"`
	BedroomsMax   string `json:"-"`
	PriceMax      string `json:"-"`
	PropertyTypes string `json:"-"`

	PostedProperties map[int64]bool `json:"properties"`
}

var Conf LocalConfig

// Launch!
func Launch() {
	// Load discord
	Conf.initDiscord()

	// Load previously posted properties
	Conf.loadConfig()

	// Intial run
	Conf.pollUpdates()

	// Run every minute!
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			Conf.pollUpdates()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// pollUpdates - check for new listings!
func (c *LocalConfig) pollUpdates() {
	err := Conf.searchTrademe()
	if err != nil {
		log.Println(err)
		return
	}

	// Update config
	c.storeConfig()
}
