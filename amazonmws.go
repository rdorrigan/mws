// Package amazonmws provides methods for interacting with the Amazon Marketplace Services API.
package amazonmws

import (
	"fmt"
	"strings"
)

const (
	bulklimit = 18
	prodAPI   = "/Products/2011-10-01"
	reportAPI = "/Reports/2009-01-01"
)

/*
GetLowestOfferListingsForASIN takes a list of ASINs and returns the result.
*/
func (api MWSAPI) GetLowestOfferListingsForASIN(items []string) (string, error) {
	params := make(map[string]string)

	for k, v := range items {
		key := fmt.Sprintf("ASINList.ASIN.%d", (k + 1))
		params[key] = string(v)
	}

	params["MarketplaceId"] = string(api.MarketplaceID)

	return api.genSignAndFetch("GetLowestOfferListingsForASIN", prodAPI, params)
}

/*
GetCompetitivePricingForASIN takes a list of ASINs and returns the result.
*/
func (api MWSAPI) GetCompetitivePricingForASIN(items []string) (string, error) {
	params := make(map[string]string)

	for k, v := range items {
		key := fmt.Sprintf("ASINList.ASIN.%d", (k + 1))
		params[key] = string(v)
	}

	params["MarketplaceId"] = string(api.MarketplaceID)

	return api.genSignAndFetch("GetCompetitivePricingForASIN", prodAPI, params)
}

// GetMatchingProductForID returns a list of products and their attributes,
// based on a list of product identifier values that you specify.
// Possible product identifiers are ASIN, GCID, SellerSKU, UPC, EAN, ISBN, and JAN.
func (api MWSAPI) GetMatchingProductForID(idType string, idList []string) (string, error) {
	params := make(map[string]string)

	for k, v := range idList {
		key := fmt.Sprintf("IdList.Id.%d", (k + 1))
		params[key] = string(v)
	}

	params["IdType"] = idType
	params["MarketplaceId"] = string(api.MarketplaceID)

	return api.genSignAndFetch("GetMatchingProductForId", prodAPI, params)
}

// GetMyPriceForSKU returns pricing information for your own offer listings,
// based on the ASIN mapped to the SellerSKU and MarketplaceId that you specify.
// Note that if you submit a SellerSKU for a product for which you don’t have an offer listing,
// the operation returns an empty Offers element.
// This operation returns pricing information for a maximum of 20 offer listings.
func (api MWSAPI) GetMyPriceForSKU(items []string) (string, error) {
	params := make(map[string]string)

	for k, v := range items {
		key := fmt.Sprintf("SellerSKUList.SellerSKU.%d", (k + 1))
		params[key] = string(v)
	}
	params["MarketplaceId"] = string(api.MarketplaceID)

	return api.genSignAndFetch("GetMyPriceForSKU", prodAPI, params)
}

// GetLowestOfferListingsForSKU takes a list of SKUs and returns the result.
func (api MWSAPI) GetLowestOfferListingsForSKU(items []string) (string, error) {
	params := make(map[string]string)

	for k, v := range items {
		key := fmt.Sprintf("SellerSKUList.SellerSKU.%d", (k + 1))
		params[key] = string(v)
	}

	params["MarketplaceId"] = string(api.MarketplaceID)

	return api.genSignAndFetch("GetLowestOfferListingsForSKU", prodAPI, params)
}

// GetLowestPricedOffersForSKU takes a single SKU and returns the result.
func (api MWSAPI) GetLowestPricedOffersForSKU(item string) (string, error) {
	params := make(map[string]string)
	// ItemCondition is a required field
	// ItemCondition values: New, Used, Collectible, Refurbished, Club.
	params["ItemCondition"] = "New"
	sku := fmt.Sprintf("SellerSKU")
	params[sku] = item
	params["MarketplaceId"] = string(api.MarketplaceID)

	return api.genSignAndFetch("GetLowestPricedOffersForSKU", prodAPI, params)
}

// GetProductCategoriesForSKU takes a single SKU and returns the result.
func (api MWSAPI) GetProductCategoriesForSKU(item string) (string, error) {
	params := make(map[string]string)
	sku := fmt.Sprintf("SellerSKU")
	params[sku] = item
	params["MarketplaceId"] = string(api.MarketplaceID)

	return api.genSignAndFetch("GetProductCategoriesForSKU", prodAPI, params)
}

// RequestReport allows for requesting a Report from reportAPI
func (api MWSAPI) RequestReport(report string, dateparams []string) (string, error) {
	params := make(map[string]string)
	l := len(dateparams)
	if l > 2 {
		return "", fmt.Errorf("Too many arguments. dateparams cannot exceed 2, a start and an end date")
	}
	report = strings.ToUpper(report)
	params["ReportType"] = report
	// Expected format time.Now().Format("2006-01-02T15:04:05-07")
	params["EndDate"] = dateparams[1]
	params["StartDate"] = dateparams[0]

	params["MarketplaceId"] = string(api.MarketplaceID)

	return api.genSignAndFetch("RequestReport", reportAPI, params)
}

// GetReportRequestList Returns a list of report requests that you can use to get the ReportRequestId for a report.
// ReportRequestIdList A structured list of ReportRequestId values. If you pass in ReportRequestId values, other query conditions are ignored.
func (api MWSAPI) GetReportRequestList(params map[string]string) (string, error) {
	// params := make(map[string]string)
	params["MarketplaceId"] = string(api.MarketplaceID)
	return api.genSignAndFetch("GetReportRequestList", reportAPI, params)
}

// GetReport Returns a list of report requests that you can use to get the ReportRequestId for a report.
func (api MWSAPI) GetReport(id string) error {
	params := make(map[string]string)
	params["RepordId"] = id
	params["MarketplaceId"] = string(api.MarketplaceID)
	id = id + ".txt"
	return api.genSignAndGet("GetReport", reportAPI, params, id)
}
