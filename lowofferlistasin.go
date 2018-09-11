package amazonmws

import (
	"encoding/xml"
	"log"
	"regexp"
	"strconv"
)

// Credit To github.com/ezkl his repos go-amazon-mws-parser and go-amazon-mws-api were the foundation for extending mine and this section is an exact copy of parser.go

// Document contains the XML results of the func GetLowestOfferListingsForASINResponse()
type Document struct {
	XMLName xml.Name `xml:"GetLowestOfferListingsForASINResponse"`
	Results []Result `xml:"GetLowestOfferListingsForASINResult"`
}

// Result is the xml container for GetLowestOfferListingsForASINResult
type Result struct {
	XMLName xml.Name `xml:"GetLowestOfferListingsForASINResult"`
	ASIN    string   `xml:"ASIN,attr"`
	Status  string   `xml:"status,attr"`
	Product *Product
}

// Product describes a Products Lowest Offer
type Product struct {
	XMLName xml.Name `xml:"Product"`
	Offers  []Offer  `xml:"LowestOfferListings>LowestOfferListing"`
}

// Offer contains the provided Offer data from GetLowestOfferListingsForASINResponse()
type Offer struct {
	XMLName              xml.Name `xml:"LowestOfferListing"`
	ConditionString      string   `xml:"Qualifiers>ItemSubcondition"`
	Condition            int
	DomesticString       string `xml:"Qualifiers>ShipsDomestically"`
	Domestic             bool
	ShippingTimeString   string `xml:"Qualifiers>ShippingTime>Max"`
	ShippingTime         int
	FeedbackRatingString string `xml:"Qualifiers>SellerPositiveFeedbackRating"`
	FeedbackRating       int
	FeedbackCount        int    `xml:"SellerFeedbackCount"`
	ListingPriceString   string `xml:"Price>ListingPrice>Amount"`
	ShippingPriceString  string `xml:"Price>Shipping>Amount"`
	ListingPrice         int
	ShippingPrice        int
}

func parseMoney(priceStr string) (priceInt int) {
	if priceStr == "0.00" {
		priceInt = 0
		return
	}

	p, err := strconv.ParseFloat(priceStr, 64)

	if err != nil {
		log.Fatal(err)
	}

	priceInt = int(p * 100.0)

	return
}

func parseCondition(condStr string) int {
	switch condStr {
	case "New":
		return 1
	case "Mint":
		return 2
	case "VeryGood":
		return 3
	case "Good":
		return 4
	}
	return 5
}

func parseDomestic(domStr string) bool {
	switch domStr {
	case "True":
		return true
	case "False":
		return false
	}
	return false
}

func parseFeedbackRating(fbStr string) int {
	ratingRegex, err := regexp.Compile(`([\d]+)%$`)

	if err != nil {
		log.Fatal(err)
	}

	if ratingRegex.Match([]byte(fbStr)) {
		m := ratingRegex.FindStringSubmatch(fbStr)
		i, err := strconv.Atoi(m[1])

		if err != nil {
			log.Fatal("Couldn't parse Feedback Rating: ", err)
		}

		return i
	}

	return -1
}

func parseMaxShipping(shipStr string) int {
	maxRegex, err := regexp.Compile(`([0-9]+) .*days`)

	if err != nil {
		log.Fatal(err)
	}

	if maxRegex.Match([]byte(shipStr)) {
		m := maxRegex.FindStringSubmatch(shipStr)
		i, err := strconv.Atoi(m[1])

		if err != nil {
			log.Fatal("Couldn't parse Max Shipping: ", err)
		}
		return i
	}

	return -1
}

// Parse parses the xml response for GetLowestOfferListingsForASINResponse()
func Parse(body []byte) (mws Document) {
	mws = Document{}

	if err := xml.Unmarshal(body, &mws); err != nil {
		log.Fatal(err)
	}

	for _, result := range mws.Results {
		for k, o := range result.Product.Offers {
			result.Product.Offers[k].ListingPrice = parseMoney(o.ListingPriceString)
			result.Product.Offers[k].ShippingPrice = parseMoney(o.ShippingPriceString)

			result.Product.Offers[k].Condition = parseCondition(o.ConditionString)
			result.Product.Offers[k].Domestic = parseDomestic(o.DomesticString)
			result.Product.Offers[k].ShippingTime = parseMaxShipping(o.ShippingTimeString)

			result.Product.Offers[k].FeedbackRating = parseFeedbackRating(o.FeedbackRatingString)
		}
	}

	return mws
}
