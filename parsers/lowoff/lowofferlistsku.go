package low

import (
	"encoding/xml"
	"log"
	"sync"
	"time"
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
		c, err := time.Parse(p.Product.Summary.OffersAvailableTime, time.RFC3339)
		if err != nil {
			p.Product.Summary.ParsedTime = c
		}
	}
}

// XMLResponse contains the XML results of the func GetLowestOfferListingsForSKU
type XMLResponse struct {
	XMLName xml.Name    `xml:"GetLowestOfferListingsForSKUResponse"`
	Results []XMLResult `xml:"GetLowestOfferListingsForSKUResult"`
}

// XMLResult is the xml container for GetMyPriceForSKU() Responses
type XMLResult struct {
	XMLName xml.Name `xml:"GetLowestOfferListingsForSKUResult"`
	// ASIN    string   `xml:"SellerSKU,attr"`
	SellerSKU string `xml:"SellerSKU,attr"`
	Status    string `xml:"status,attr"`
	TooSoon   bool
	Product   Product
}

// Product describes a Products Identifiers & Offer
type Product struct {
	XMLName     xml.Name            `xml:"Product"`
	Identifiers Identifier          `xml:"Identifiers"`
	Offers      LowestOfferListings `xml:"LowestOfferListings"`
	Summary     Summary             `xml:"Summary"`
}

// LowestOfferListings contains each LowestOfferListing
type LowestOfferListings struct {
	LowestOfferListing []LowestOfferListing `xml:"LowestOfferListing"`
}

// Identifier describes ASIN & SellerSKU.
type Identifier struct {
	XMLName xml.Name `xml:"Identifiers"`
	ASIN    string   `xml:"MarketplaceASIN>ASIN"`
	SKU     string   `xml:"SKUIdentifier>SellerSKU"`
}

// LowestOfferListing contains the provided Offer data from GetLowestOfferListingsForASINResponse()
type LowestOfferListing struct {
	XMLName                         xml.Name   `xml:"LowestOfferListing"`
	Qualifier                       Qualifiers `xml:"Qualifiers"`
	NumberOfOfferListingsConsidered int        `xml:"NumberOfOfferListingsConsidered"`
	SellerFeedbackCount             int        `xml:"SellerFeedbackCount"`
	Price                           Price      `xml:"Price"`
	MultipleOffersAtLowestPrice     string     `xml:"MultipleOffersAtLowestPrice"`
}

// Summary is returned when status="ActiveButTooSoonForProcessing"
type Summary struct {
	XMLName             xml.Name `xml:"Summary"`
	TotalOfferCount     int      `xml:"TotalOfferCount"`
	OffersAvailableTime string   `xml:"OffersAvailableTime"`
	ParsedTime          time.Time
}

// Qualifiers contains all low offer Qualifiers
type Qualifiers struct {
	XMLName                      xml.Name `xml:"Qualifiers"`
	ItemCondition                string   `xml:"ItemCondition"`
	ItemSubCondition             string   `xml:"ItemSubcondition"`
	FulfillmentChannel           string   `xml:"FulfillmentChannel"`
	ShipsDomestically            string   `xml:"ShipsDomestically"`
	ShippingTime                 string   `xml:"ShippingTime>Max"`
	SellerPositiveFeedbackRating string   `xml:"SellerPositiveFeedbackRating"`
}

// Price contains the pricing fields within BuyingPrice
// ListingPrice – The current price including any promotions that apply to the product.
// Shipping – The shipping cost of the product.
// LandedPrice – ListingPrice plus Shipping.
// Note that if the landed price is not returned, the listing price represents the product with the lowest landed price.
type Price struct {
	XMLName       xml.Name `xml:"Price"`
	LandedPrice   string   `xml:"LandedPrice>Amount"`
	CurrencyCode  string   `xml:"LandedPrice>CurrencyCode"`
	ListingPrice  string   `xml:"ListingPrice>Amount"`
	ShippingPrice string   `xml:"Shipping>Amount"`
}
