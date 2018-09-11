package mp

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
		log.Fatal(err)
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

// XMLResponse contains the XML results of the func GetMyPriceForSKU()
type XMLResponse struct {
	XMLName xml.Name    `xml:"GetMyPriceForSKUResponse"`
	Results []XMLResult `xml:"GetMyPriceForSKUResult"`
}

// XMLResult is the xml container for GetMyPriceForSKU() Responses
type XMLResult struct {
	XMLName xml.Name `xml:"GetMyPriceForSKUResult"`
	// ASIN    string   `xml:"SellerSKU,attr"`
	SellerSKU string `xml:"SellerSKU,attr"`
	Status    string `xml:"status,attr"`
	TooSoon   bool
	Product   Product
	// Identifiers Identifier `xml:"Identifiers"`
	// Offers      []Offer    /* `xml:"Offers"`*/
}

// Product describes a Products Identifiers & Offer
type Product struct {
	XMLName     xml.Name   `xml:"Product"`
	Identifiers Identifier `xml:"Identifiers"`
	Offers      []Offer    `xml:">Offer"`
	Summary     Summary    `xml:"Summary"`
}

// Identifier describes ASIN & SellerSKU.
// MarketplaceId & SellerId are not returned
type Identifier struct {
	XMLName xml.Name `xml:"Identifiers"`
	// ASIN    string   `xml:"ASIN"`
	// SKU     string   `xml:"SellerSKU"`
	ASIN string `xml:"MarketplaceASIN>ASIN"`
	SKU  string `xml:"SKUIdentifier>SellerSKU"`
}

// Offer contains the provided Offer data from GetLowestOfferListingsForASINResponse()
type Offer struct {
	XMLName xml.Name `xml:"Offer"`
	// Offer       string       `xml:">Offer"`
	BuyingPrice BuyingPrice `xml:"BuyingPrice"`
	// The current price excluding any promotions that apply to the product. Excludes the shipping cost.
	RegularPrice string `xml:"RegularPrice>Amount"`
	// CurrenyCode can be pulled from every price element within BuyingPrice
	CurrencyCode       string `xml:"RegularPrice>CurrencyCode"`
	FulfillmentChannel string `xml:"FulfillmentChannel"`
	ItemCondition      string `xml:"ItemCondition"`
	ItemSubCondition   string `xml:"ItemSubCondition"`
}

// Summary is returned when status="ActiveButTooSoonForProcessing"
type Summary struct {
	XMLName             xml.Name `xml:"Summary"`
	TotalOfferCount     int      `xml:"TotalOfferCount"`
	OffersAvailableTime string   `xml:"OffersAvailableTime"`
	ParsedTime          time.Time
}

// BuyingPrice contains the pricing fields within BuyingPrice
// ListingPrice – The current price including any promotions that apply to the product.
// Shipping – The shipping cost of the product.
// LandedPrice – ListingPrice plus Shipping.
// Note that if the landed price is not returned, the listing price represents the product with the lowest landed price.
type BuyingPrice struct {
	XMLName       xml.Name `xml:"BuyingPrice"`
	LandedPrice   string   `xml:"LandedPrice>Amount"`
	ListingPrice  string   `xml:"ListingPrice>Amount"`
	ShippingPrice string   `xml:"Shipping>Amount"`
}
