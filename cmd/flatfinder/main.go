package main

import (
	"flatfinder/internal/flatfinder"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load .env file in current directory")
	}
}

func main() {
	// Load env vars and validate
	flatfinder.Conf = flatfinder.LocalConfig{}

	// Load webhook
	flatfinder.Conf.DiscordWebhook = os.Getenv("DISCORD_WEBHOOK")
	if flatfinder.Conf.DiscordWebhook == "" {
		log.Fatal("DISCORD_WEBHOOK not set")
	}
	flatfinder.Conf.DiscordTag = os.Getenv("DISCORD_TAG")

	// Load Google stuff
	flatfinder.Conf.GoogleApiToken = os.Getenv("GOOGLE_API_KEY")
	if flatfinder.Conf.GoogleApiToken == "" {
		log.Print("GOOGLE_API_KEY not set. Not using map logicc")
	}
	flatfinder.Conf.GoogleLocation1 = os.Getenv("GOOGLE_LOCATION_1")
	flatfinder.Conf.GoogleLocation2 = os.Getenv("GOOGLE_LOCATION_2")

	// Load trademe config
	flatfinder.Conf.TradeMeKey = os.Getenv("TRADEME_API_KEY")
	flatfinder.Conf.TradeMeSecret = os.Getenv("TRADEME_API_SECRET")
	if flatfinder.Conf.TradeMeKey == "" || flatfinder.Conf.TradeMeSecret == "" {
		log.Fatal("TRADEME_API_KEY or TRADEME_API_SECRET not set")
	}

	// Load filterse
	flatfinder.Conf.Suburbs = os.Getenv("SUBURBS")
	if flatfinder.Conf.Suburbs == "" {
		log.Fatal("SUBURBS not set")
	}

	flatfinder.Conf.BedroomsMin = os.Getenv("BEDROOMS_MIN")
	if flatfinder.Conf.BedroomsMin == "" {
		log.Fatal("BEDROOMS_MIN not set")
	}

	flatfinder.Conf.BedroomsMax = os.Getenv("BEDROOMS_MAX")
	if flatfinder.Conf.BedroomsMax == "" {
		log.Fatal("BEDROOMS_MAX not set")
	}

	flatfinder.Conf.PriceMax = os.Getenv("PRICE_MAX")
	if flatfinder.Conf.PriceMax == "" {
		log.Fatal("PRICE_MAX not set")
	}

	flatfinder.Conf.PropertyTypes = os.Getenv("PROPERTY_TYPE")
	if flatfinder.Conf.PropertyTypes == "" {
		log.Fatal("PROPERTY_TYPE not set")
	}

	// Start the stuff
	flatfinder.Launch()
}
