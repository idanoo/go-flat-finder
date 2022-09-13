# FlatFinder Bot

[![Build Status](https://ci.tinker.nz/api/badges/idanoo/flat-finder/status.svg)](https://ci.tinker.nz/idanoo/flat-finder) 
    

* Uses the Trade Me API to grab new rental properties that have been recently listed
* Checks if fibre and VDSL are available by querying Chorus
* Includes travel times to various locations

## Requirements
* Linux environment
* Trade Me API Key [(register an application)](https://www.trademe.co.nz/MyTradeMe/Api/RegisterNewApplication.aspx)

## Optionals
* Google Distance Matrix API Key [(get a key)](https://developers.google.com/maps/documentation/distance-matrix/start#get-a-key)
* Discord API key

## Installation
* Download latest build
* Run with below exe with environment variables set

## Configuration
Copy `.env.example`to `.env` and set variables. Leave blank to disable parts.

```
SINCE="2 hours ago"
DISCORD_WEBHOOK="abcd"
GOOGLE_API_KEY="abcd"
GOOGLE_LOCATION_1="42 Wallaby Way, Sydney"
GOOGLE_LOCATION_2="43 Wallaby Way, Sydney"
DISTRICTS="47,52"
BEDROOMS_MIN="2"
BEDROOMS_MAX="4"
PRICE_MAX="700"
PROPERTY_TYPE="House,Townhouse,Apartment"
```

Reference: [http://developer.trademe.co.nz/api-reference/search-methods/rental-search/](http://developer.trademe.co.nz/api-reference/search-methods/rental-search/)
