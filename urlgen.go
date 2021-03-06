package amazonmws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

// MWSAPI contains required client config information
type MWSAPI struct {
	AccessKey     string
	SecretKey     string
	Host          string
	AuthToken     string
	MarketplaceID string
	SellerID      string
}

func (api MWSAPI) genSignAndFetch(Action string, ActionPath string, Parameters map[string]string) (string, error) {
	genURL, err := GenerateAmazonURL(api, Action, ActionPath, Parameters)
	if err != nil {
		return "", err
	}

	SetTimestamp(genURL)

	signedurl, err := SignAmazonURL(genURL, api)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(signedurl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func (api MWSAPI) genSignAndGet(Action string, ActionPath string, Parameters map[string]string, dst string) error {
	genURL, err := GenerateAmazonURL(api, Action, ActionPath, Parameters)
	if err != nil {
		return err
	}

	SetTimestamp(genURL)

	signedurl, err := SignAmazonURL(genURL, api)
	if err != nil {
		return err
	}

	resp, err := http.Get(signedurl)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		out.Close()
	}()
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	// if err != nil {
	// 	return "", err
	// }

	return nil
}

// GenerateAmazonURL prepares the url in genSignAndFetch
func GenerateAmazonURL(api MWSAPI, Action string, ActionPath string, Parameters map[string]string) (finalURL *url.URL, err error) {
	result, err := url.Parse(api.Host)
	if err != nil {
		return nil, err
	}

	result.Host = api.Host
	result.Scheme = "https"
	result.Path = ActionPath

	values := url.Values{}
	values.Add("Action", Action)

	if api.AuthToken != "" {
		values.Add("MWSAuthToken", api.AuthToken)
	}

	values.Add("AWSAccessKeyId", api.AccessKey)
	values.Add("SellerId", api.SellerID)
	values.Add("SignatureVersion", "2")
	values.Add("SignatureMethod", "HmacSHA256")
	// split allows for accessing other API References ex: Reports
	// values.Add("Version", "2011-10-01")
	split := strings.Split(ActionPath, "/")
	values.Add("Version", split[len(split)-1])

	for k, v := range Parameters {
		values.Set(k, v)
	}

	params := values.Encode()
	result.RawQuery = params

	return result, nil
}

// SetTimestamp adds a RFC3339 timestamp to the URL
func SetTimestamp(origURL *url.URL) (err error) {
	values, err := url.ParseQuery(origURL.RawQuery)
	if err != nil {
		return err
	}
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))
	origURL.RawQuery = values.Encode()

	return nil
}

// SignAmazonURL encodes the SecretKey signing the URL
func SignAmazonURL(origURL *url.URL, api MWSAPI) (signedURL string, err error) {
	escapeURL := strings.Replace(origURL.RawQuery, ",", "%2C", -1)
	escapeURL = strings.Replace(escapeURL, ":", "%3A", -1)

	params := strings.Split(escapeURL, "&")
	sort.Strings(params)
	sortedParams := strings.Join(params, "&")

	toSign := fmt.Sprintf("GET\n%s\n%s\n%s", origURL.Host, origURL.Path, sortedParams)

	hasher := hmac.New(sha256.New, []byte(api.SecretKey))
	_, err = hasher.Write([]byte(toSign))
	if err != nil {
		return "", err
	}

	hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	hash = url.QueryEscape(hash)

	newParams := fmt.Sprintf("%s&Signature=%s", sortedParams, hash)

	origURL.RawQuery = newParams

	return origURL.String(), nil
}
