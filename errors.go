package amazonmws

import (
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
)

// ErrorResponse contains the returned xml errors
type ErrorResponse struct {
	Error    error
	Response XMLErrorResponse
	Throttle Throttle
}

// NOT NEEDED
// ErrorCode contains common MWS errors
// type ErrorCode struct {
// 	Name    string
// 	Code    int
// 	Message string /*Description*/
// }

// XMLErrorResponse is the first element returned indicating an ErrorrResponse
type XMLErrorResponse struct {
	XMLName   xml.Name            `xml:"ErrorResponse"`
	Error     []XMLResponseErrors `xml:"Error"`
	RequestID XMLRequestID        `xml:"RequestID"`
}

// XMLResponseErrors contains the type of Errror
type XMLResponseErrors struct {
	XMLName xml.Name `xml:"Error"`
	Type    string   `xml:"Type"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
	Detail  string   `xml:"Detail"`
	SKU     string   `xml:"SKU,attr"`
}

// XMLRequestID holds the ID
type XMLRequestID struct {
	XMLName xml.Name `xml:"RequestID"`
	ID      string   `xml:"RequestID"`
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

const (
	// Disconnect will terminate
	Disconnect = "InputStreamDisconnected"
	// Parameter will terminate
	Parameter = "InvalidParameterValue"
	// AccessDenied will terminate
	AccessDenied = "AccessDenied"
	// InvalidAccessKey will terminate
	InvalidAccessKey = "InvalidAccessKeyId"
	// SignDoesNotMatch will terminate
	SignDoesNotMatch = "SignatureDoesNotMatch"
	// InvalidAddress will terminate
	InvalidAddress = "InvalidAddress"
	// InternalError will be throttled
	InternalError = "InternalError"
	// QuotaExceeded will be throttled
	QuotaExceeded = "QuotaExceeded"
	// ReqThrottled will be throttled
	ReqThrottled = "RequestThrottled"
)

// NewError implements ErrorResponse.Error
func (er *ErrorResponse) NewError(i int) {
	er.Error = fmt.Errorf("error Code: %s\nresponse Message: %s", er.Response.Error[i].Code, er.Response.Error[i].Message)

}

// ParseError parses an mws xml error response
func (p *XMLParser) ParseError(e error) *ErrorResponse {
	ER := ErrorResponse{}
	b := []byte(e.Error())
	if err := xml.Unmarshal(b, &ER.Response); err != nil {
		ER.Error = err
	}

	return &ER
}

// ResolveError finds a solution to errors that are due to throttling
// and exits when the error is due to an invalid request
func (er *ErrorResponse) ResolveError() {
	for k, e := range er.Response.Error {
		er.CheckCode(k)
		fmt.Println("Error: ", e.Code)
		if er.Throttle.Throttled == true {
			er.Throttle.NewTicker()
		} else {
			er.NewError(k)
		}
	}

}

// CheckCode iterates through the const valid codes
func (er *ErrorResponse) CheckCode(i int) {
	validcodes := []string{
		InternalError,
		QuotaExceeded,
		ReqThrottled,
	}
	exitcodes := []string{
		Disconnect,
		Parameter,
		AccessDenied,
		InvalidAccessKey,
		SignDoesNotMatch,
		InvalidAddress,
	}
	for _, s := range validcodes {
		if strings.Compare(s, er.Response.Error[i].Code) == 0 {
			er.Throttle.Throttled = true
			return
		}
	}
	for _, e := range exitcodes {
		if strings.Compare(e, er.Response.Error[i].Code) == 0 {
			er.Throttle.Throttled = false
			return
		}
	}
}

// Implement all
/*Error code	HTTP status code	Description
InputStreamDisconnected	400	There was an error reading the input stream.
InvalidParameterValue	400	An invalid parameter value was used, or the request size exceeded the maximum accepted size, or the request expired.
AccessDenied	401	Access was denied.
InvalidAccessKeyId	403	An invalid AWSAccessKeyId value was used.
SignatureDoesNotMatch	403	The signature used does not match the server's calculated signature value.
InvalidAddress	404	An invalid API section or operation value was used, or an invalid path was used.
InternalError	500	There was an internal service failure.
QuotaExceeded	503	The total number of requests in an hour was exceeded.
RequestThrottled	503	The frequency of requests was greater than allowed.*/

/*Service errors
The common response to a 500 or 503 service error is to try the request again.
Such service errors are usually only temporary and will resolve themselves.
If you want to retry an operation call after receiving a 500 or 503 error,
you can immediately retry after the first error response.
However, if you want to retry multiple times, Amazon recommends that you implement an
"exponential backoff" approach (i.e. pausing between retrys), with up to four retries.
Then, log the error and proceed with a manual follow-up and investigation.
For example, you can time your retries with the following time spacing: 1s, 4s, 10s, 30s.
The actual backoff times and limit will depend upon your business processes.*/

// NOT NEEDED
// ErrorHandler returns an ErrorResponse. ONLY PASS IN PARSED XML
// func ErrorHandler(name string, code int, msg string) *ErrorResponse {
// 	if name != "" && code != 0 && msg != "" {
// 		ec := ErrorCode{
// 			Name:    name,
// 			Code:    code,
// 			Message: msg,
// 		}
// 		return &ErrorResponse{Code: ec}
// 	}
// 	return &ErrorResponse{Error: errors.New("ErrorHandler is missing required data")}
// }
