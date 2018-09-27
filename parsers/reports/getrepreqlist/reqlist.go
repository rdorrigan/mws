package getrepreqlist

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

// Parser parses the xml response for MWS ReportRequestInfos operations
func (p *XMLParser) Parser(body []byte) *XMLResponse {
	var i XMLResponse
	if err := xml.Unmarshal(body, &i); err != nil {
		log.Fatal(err)
	}
	return &i
}

// XMLResponse contains the XML results of the func GetMyPriceForSKU()
type XMLResponse struct {
	XMLName          xml.Name         `xml:"GetReportRequestListResponse"`
	Result           XMLResult        `xml:"GetReportRequestListResult"`
	ResponseMetadata ResponseMetadata `xml:"ResponseMetadata"`
}

// XMLResult is the xml container for GetMyPriceForSKU() Responses
type XMLResult struct {
	XMLName   xml.Name `xml:"GetReportRequestListResult"`
	NextToken string   `xml:"NextToken"`
	HasNext   bool     `xml:"HasNext"`
	Info      []Info   `xml:"ReportRequestInfo"`
}

// Info describes a ReportRequestInfos Identifiers & Offer
type Info struct {
	XMLName                xml.Name `xml:"ReportRequestInfo"`
	ReportRequestID        string   `xml:"ReportRequestId"`
	ReportType             string   `xml:"ReportType"`
	StartDate              string   `xml:"StartDate"`
	EndDate                string   `xml:"EndDate"`
	Scheduled              bool     `xml:"Scheduled"`
	SubmittedDate          string   `xml:"SubmittedDate"`
	ReportProcessingStatus string   `xml:"ReportProcessingStatus"`
	GeneratedReportID      string   `xml:"GeneratedReportId"`
	StartedProcessingDate  string   `xml:"StartedProcessingDate"`
	CompletedDate          string   `xml:"CompletedDate"`
}

// ResponseMetadata contains the ID
type ResponseMetadata struct {
	XMLName   xml.Name `xml:"ResponseMetadata"`
	RequestID string   `xml:"RequestId"`
}
