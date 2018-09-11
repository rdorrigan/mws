package lowp

import (
	"encoding/xml"
	"log"
	"sync"
	"time"
)

// Status values per request
const (
	// Success is returned on a successful request
	Success = "Success"
	// NoBuyableOffers is returned when there are no buyable offers
	NoBuyableOffers = "NoBuyableOffers"
	// NoOfferDueToMissingShippingCharge is returned when there are no shipping charges
	NoOfferDueToMissingShippingCharge = "NoOfferDueToMissingShippingCharge"
	// ActiveButTooSoonForProcessing is due to Throttling
	ActiveButTooSoonForProcessing = "ActiveButTooSoonForProcessing"
)

// XMLParse extends Parser
type XMLParse interface {
	Parser(body []byte)
}

// XMLParser represents an XML parser.
type XMLParser struct {
	decoder  *xml.Decoder
	decMutex *sync.Mutex
	mapMutex *sync.Mutex
}

// NewXMLParser creates a new XML parser.
func NewXMLParser() *XMLParser {
	return &XMLParser{nil, &sync.Mutex{}, &sync.Mutex{}}
}

// Parser parses the xml response for MWS Products operations
func (p *XMLParser) Parser(body []byte) *XMLResponse {
	var i XMLResponse
	if err := xml.Unmarshal(body, &i); err != nil {
		log.Println(err)
	}
	i.tooSoon()
	for _, r := range i.Results {
		if r.TooSoon {
			i.parseTime()
		}
	}
	return &i
}
func (r *XMLResponse) tooSoon() {
	for _, p := range r.Results {
		if p.Status == "ActiveButTooSoonForProcessing" {
			p.TooSoon = true
		}
	}
}
func (r *XMLResponse) parseTime() {
	for _, p := range r.Results {

		c, err := time.Parse(p.Product.Identifiers.TimeOfOfferChange, time.RFC3339)
		if err != nil {
			p.Product.Identifiers.ParsedTime = c
		}
	}
}
func (r *XMLResponse) parseStatus() {
	for _, i := range r.Results {
		switch i.Status {
		case Success:
			i.Statusbool = true
		case NoBuyableOffers:
			i.Statusbool = false
		case NoOfferDueToMissingShippingCharge:
			i.Statusbool = false
		case ActiveButTooSoonForProcessing:
			i.Statusbool = false
		default:
			i.Statusbool = false
		}
	}
}

// XMLResponse contains the XML results of the func GetLowestOfferListingsForSKU
type XMLResponse struct {
	XMLName          xml.Name         `xml:"GetLowestPricedOffersForSKUResponse"`
	Results          []XMLResult      `xml:"GetLowestPricedOffersForSKUResult"`
	ResponseMetadata ResponseMetadata `xml:"ResponseMetadata"`
}

// XMLResult is the xml container for GetMyPriceForSKU() Responses
type XMLResult struct {
	XMLName       xml.Name `xml:"GetLowestPricedOffersForSKUResult"`
	MarketplaceID string   `xml:"MarketplaceID,attr"`
	SKU           string   `xml:"SKU,attr"`
	ItemCondition string   `xml:"ItemCondition,attr"`
	Status        string   `xml:"status,attr"`
	Statusbool    bool
	TooSoon       bool
	Product       Product
}

// Product describes a Products Identifiers & Offer
type Product struct {
	XMLName     xml.Name   `xml:"Product"`
	Identifiers Identifier `xml:"Identifier"`
	Offers      []Offer    `xml:"Offers"`
	Summary     Summary    `xml:"Summary"`
}

// Identifier describes ASIN & SellerSKU.
type Identifier struct {
	XMLName           xml.Name `xml:"Identifier"`
	MarketplaceID     string   `xml:"MarketplaceID"`
	SellerSKU         string   `xml:"SellerSKU"`
	ItemCondition     string   `xml:"ItemCondition"`
	TimeOfOfferChange string   `xml:"TimeOfOfferChange"`
	ParsedTime        time.Time
}

// Offer contains the provided Offer data from GetLowestPricedOffersForSKU
type Offer struct {
	XMLName              xml.Name             `xml:"Offer"`
	MyOffer              string               `xml:"MyOffer"`
	SubCondition         string               `xml:"SubCondition"`
	SellerFeedbackRating SellerFeedbackRating `xml:"SellerFeedbackRating"`
	ShippingTime         ShippingTime         `xml:"ShippingTime"`
	ListingPrice         string               `xml:"ListingPrice>Amount"`
	ShippingPrice        string               `xml:"ShippingPrice>Amount"`
	IsFulfilledByAmazon  string               `xml:"IsFulfilledByAmazon"`
	IsBuyBoxWinner       string               `xml:"IsBuyBoxWinner"`
	IsFeaturedMerchant   string               `xml:"IsFeaturedMerchant"`
}

// Summary contains price information about the product
type Summary struct {
	XMLName                         xml.Name                        `xml:"Summary"`
	TotalOfferCount                 int                             `xml:"TotalOfferCount"`
	NumberOfOffers                  NumberOfOffers                  `xml:"NumberOfOffers"`
	LowestPrices                    []LowestPrice                   `xml:"LowestPrices"`
	BuyBoxPrices                    []BuyBoxPrice                   `xml:"BuyBoxPrices"`
	ListPrice                       ListPrice                       `xml:"ListPrice"`
	SuggestedLowerPricePlusShipping SuggestedLowerPricePlusShipping `xml:"SuggestedLowerPricePlusShipping"`
	BuyBoxEligibleOffers            BuyBoxEligibleOffers            `xml:"BuyBoxEligibleOffers"`
	// OffersAvailableTime string   `xml:"OffersAvailableTime"`
	// ParsedTime          time.Time

}

// NumberOfOffers contains some offer attributes
type NumberOfOffers struct {
	XMLName            xml.Name `xml:"NumberOfOffers"`
	OfferCount         int      `xml:"OfferCount"`
	Condition          string   `xml:"condition,attr"`
	FulfillmentChannel string   `xml:"fulfillmentChannel,attr"`
}

// LowestPrice contains the pricing fields within BuyingPrice
// ListingPrice – The current price including any promotions that apply to the product.
// Shipping – The shipping cost of the product.
// LandedPrice – ListingPrice plus Shipping.
// Note that if the landed price is not returned, the listing price represents the product with the lowest landed price.
type LowestPrice struct {
	XMLName            xml.Name `xml:"LowestPrice"`
	Condition          string   `xml:"condition,attr"`
	FulfillmentChannel string   `xml:"fulfillmentChannel,attr"`
	LandedPrice        string   `xml:"LandedPrice>Amount"`
	CurrencyCode       string   `xml:"LandedPrice>CurrencyCode"`
	ListingPrice       string   `xml:"ListingPrice>Amount"`
	ShippingPrice      string   `xml:"Shipping>Amount"`
}

// BuyBoxPrice describes buy box pricing
type BuyBoxPrice struct {
	XMLName       xml.Name `xml:"BuyBoxPrice"`
	Condition     string   `xml:"condition,attr"`
	LandedPrice   string   `xml:"LandedPrice>Amount"`
	CurrencyCode  string   `xml:"LandedPrice>CurrencyCode"`
	ListingPrice  string   `xml:"ListingPrice>Amount"`
	ShippingPrice string   `xml:"Shipping>Amount"`
}

// ListPrice has currency and and an amount
type ListPrice struct {
	XMLName      xml.Name `xml:"ListPrice"`
	CurrencyCode string   `xml:"CurrencyCode"`
	Amount       string   `xml:"Amount"`
}

// SuggestedLowerPricePlusShipping has currency and and an amount
type SuggestedLowerPricePlusShipping struct {
	XMLName      xml.Name `xml:"SuggestedLowerPricePlusShipping"`
	CurrencyCode string   `xml:"CurrencyCode"`
	Amount       string   `xml:"Amount"`
}

// BuyBoxEligibleOffers is the last Summary element
type BuyBoxEligibleOffers struct {
	XMLName            xml.Name `xml:"BuyBoxEligibleOffers"`
	OfferCount         int      `xml:"OfferCount"`
	Condition          string   `xml:"condition,attr"`
	FulfillmentChannel string   `xml:"fulfillmentChannel,attr"`
}

// SellerFeedbackRating contains feedback info
type SellerFeedbackRating struct {
	XMLName                      xml.Name `xml:"SellerFeedbackRating"`
	SellerPositiveFeedbackRating string   `xml:"SellerPositiveFeedbackRating"`
	FeedbackCount                string   `xml:"FeedbackCount"`
}

// ShippingTime contains min man and availability
type ShippingTime struct {
	XMLName          xml.Name `xml:"ShippingTime"`
	MinHours         string   `xml:"minimumHours,attr"`
	MaxHours         string   `xml:"maximumHours,attr"`
	AvailabilityType string   `xml:"availabilityType,attr"`
}

// ResponseMetadata contains the RequestID
type ResponseMetadata struct {
	XMLName   xml.Name `xml:"ResponseMetadata"`
	RequestID string   `xml:"RequestId"`
}
