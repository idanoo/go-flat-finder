package flatfinder

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

// Load discord client
func (c *LocalConfig) initDiscord() {
	// Webhook URL splitting
	webhookString := strings.ReplaceAll(c.DiscordWebhook, "https://discord.com/api/webhooks/", "")
	webhookParts := strings.Split(webhookString, "/")
	if len(webhookParts) != 2 {
		log.Fatal("Invalid DISCORD_WEBHOOK")
	}

	// Convert snowflakeID to uint64
	i, err := strconv.ParseInt(webhookParts[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Start client!
	client := webhook.New(snowflake.ID(i), webhookParts[1])
	c.DiscordClient = client

	log.Print("Discord client loaded succesfully")
}

// sendEmbeddedMessage - Build an embedded message from listing data
func (c *LocalConfig) sendEmbeddedMessage(listing TradeMeListing) {
	log.Printf("New listing: %s", listing.Title)

	hasFibre, currentConn := getAvailableSpeeds(
		fmt.Sprintf(
			"%s, %s, %s",
			strings.TrimSpace(listing.Address),
			strings.TrimSpace(listing.Suburb),
			strings.TrimSpace(listing.Region),
		),
	)

	embed := discord.NewEmbedBuilder().
		SetTitle(listing.Title).
		SetURL(fmt.Sprintf("https://trademe.co.nz/%d", listing.ListingID)).
		SetColor(1127128).
		SetImage(listing.PictureHref).
		AddField(
			"Location",
			fmt.Sprintf(
				"[%s](https://maps.google.com/maps?z=12&t=m&q=loc:%f+%f)",
				listing.Address,
				listing.GeographicLocation.Latitude,
				listing.GeographicLocation.Longitude,
			),
			true,
		).
		AddField("Bedrooms", fmt.Sprintf("%d", listing.Bedrooms), true).
		AddField("Fibre Avail", hasFibre, false).
		AddField("Current Connection", currentConn, false)

	// Only add address if token set
	if c.GoogleApiToken != "" && c.GoogleLocation1 != "" {
		distance := getDistanceFromAddress(listing.Address, c.GoogleLocation1)
		embed = embed.AddField(fmt.Sprintf("Distance to %s", c.GoogleLocation1), distance, false)
	}

	embeds := []discord.Embed{}
	embeds = append(embeds, embed.Build())
	_, err := c.DiscordClient.CreateEmbeds(embeds)
	if err != nil {
		log.Print(err)
	}
}
