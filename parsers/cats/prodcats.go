package low

import (
	"encoding/xml"
	"log"
	"sync"
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
	// i.tooSoon()
	// for _, r := range i.Results {
	// 	if r.TooSoon {
	// 		i.parseTime()
	// 	}
	// }
	return &i
}

// func (r *XMLResponse) tooSoon() {
// 	for _, p := range r.Results {
// 		if p.Status == "ActiveButTooSoonForProcessing" {
// 			p.TooSoon = true
// 		}
// 	}
// }

// XMLResponse contains the XML results of the func GetLowestOfferListingsForSKU
type XMLResponse struct {
	XMLName          xml.Name         `xml:"GetProductCategoriesForSKUResponse"`
	Results          XMLResult        `xml:"GetProductCategoriesForSKUResult"`
	ResponseMetadata ResponseMetadata `xml:"ResponseMetadata"`
}

// XMLResult is the xml container for GetProductCategoriesForSKU() Responses
type XMLResult struct {
	XMLName xml.Name `xml:"GetProductCategoriesForSKUResult"`
	// ASIN    string   `xml:"SellerSKU,attr"`
	// SellerSKU string `xml:"SellerSKU,attr"`
	// Status    string `xml:"status,attr"`
	Self    []Self `xml:"Self"`
	TooSoon bool
}

// ResponseMetadata returns a RequestID
type ResponseMetadata struct {
	XMLName   xml.Name `xml:"ResponseMetadata"`
	RequestID string   `xml:"RequestId"`
}

// Self Contains the ProductCategoryId for the product that you submitted.
// Also contains a ProductCategoryId for each of the parent categories of the product,
//  up to the root for the Marketplace.
type Self struct {
	XMLName             xml.Name  `xml:"Self"`
	ProductCategoryID   string    `xml:"ProductCategoryId"`
	ProductCategoryName string    `xml:"ProductCategoryName"`
	Parents             []Parents `xml:"Parent"`
}

// Parents contains the potentially many layers of "Parent" categories
// See Example responses for clarification
type Parents struct {
	Parents []Parent
}

// Parent describes a category of the product
type Parent struct {
	XMLName             xml.Name `xml:"Parent"`
	ProductCategoryID   string   `xml:"ProductCategoryId"`
	ProductCategoryName string   `xml:"ProductCategoryName"`
}
