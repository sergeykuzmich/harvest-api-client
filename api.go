package harvest_api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const CLIENT_VERSION = "1.1.2"
const HARVEST_DOMAIN = "api.harvestapp.com"
const HARVEST_API_VERSION = "v2"

type API struct {
	client       *http.Client
	BaseURL      string
	AccountID    string
	AccessToken  string
	RefreshToken string
}

func HarvestClient(accountID string, accessToken string) *API {
	a := API{}
	a.client = http.DefaultClient
	a.BaseURL = "https://" + HARVEST_DOMAIN + "/" + HARVEST_API_VERSION
	a.AccountID = accountID
	a.AccessToken = accessToken
	return &a
}

func (a *API) GetPaginated(path string, args Arguments, target Pageable, afterFetch func()) error {
	page := 1
	args["page"] = fmt.Sprintf("%d", page)
	err := a.Get(path, args, target)
	if err != nil {
		return err
	}

	afterFetch()

	for target.HasNextPage() {
		page++
		args["page"] = fmt.Sprintf("%d", page)
		err = a.Get(path, args, target)
		if err != nil {
			return err
		}
		afterFetch()
	}
	return nil
}

func (a *API) Get(path string, args Arguments, target interface{}) error {
	url := fmt.Sprintf("%s%s", a.BaseURL, path)
	urlWithParams := fmt.Sprintf("%s?%s", url, args.ToURLValues().Encode())

	req, err := http.NewRequest("GET", urlWithParams, nil)
	if err != nil {
		return errors.Wrapf(err, "Invalid GET request %s", url)
	}
	a.AddHeaders(req)

	res, err := a.client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "HTTP request failure on %s", url)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "HTTP request failure on %s: %s %s", url, string(body), err)
		}
		return errors.Errorf("HTTP request failure on %s: %s", url, string(body))
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(target)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "JSON decode failed on %s: %s", url, string(body))
	}

	return nil
}

func (a *API) Put(path string, args Arguments, postData interface{}, target interface{}) error {
	return a.PPP("PUT", path, args, postData, target)
}

func (a *API) Patch(path string, args Arguments, postData interface{}, target interface{}) error {
	return a.PPP("PATCH", path, args, postData, target)
}

func (a *API) Post(path string, args Arguments, postData interface{}, target interface{}) error {
	return a.PPP("POST", path, args, postData, target)
}

func (a *API) PPP(method string, path string, args Arguments, postData interface{}, target interface{}) error {
	url := fmt.Sprintf("%s%s", a.BaseURL, path)
	urlWithParams := fmt.Sprintf("%s?%s", url, args.ToURLValues().Encode())

	buffer := new(bytes.Buffer)
	if postData != nil {
		json.NewEncoder(buffer).Encode(postData)
	}

	req, err := http.NewRequest(method, urlWithParams, buffer)
	if err != nil {
		return errors.Wrapf(err, "Invalid %s request %s", method, url)
	}
	a.AddHeaders(req)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := a.client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "HTTP request failure on %s", url)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "HTTP request failure on %s: %s %s", url, string(body), err)
		}
		return errors.Errorf("HTTP request failure on %s: %s", url, string(body))
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(target)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "JSON decode failed on POST to %s: %s", url, string(body))
	}

	return nil
}

func (a *API) Delete(path string, args Arguments) error {
	url := fmt.Sprintf("%s%s", a.BaseURL, path)
	urlWithParams := fmt.Sprintf("%s?%s", url, args.ToURLValues().Encode())

	req, err := http.NewRequest("DELETE", urlWithParams, nil)
	if err != nil {
		return errors.Wrapf(err, "Invalid DELETE request %s", url)
	}
	a.AddHeaders(req)

	res, err := a.client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "HTTP request failure on %s", url)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "HTTP request failure on %s: %s %s", url, string(body), err)
		}
		return errors.Errorf("HTTP request failure on %s: %s", url, string(body))
	}

	return nil
}

// Applies relevant User-Agent, Accept & Authorization
func (a *API) AddHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "github.com/sergeykuzmich/harvest-api-client v"+CLIENT_VERSION)
	req.Header.Set("Harvest-Account-Id", a.AccountID)
	req.Header.Set("Authorization", "Bearer "+a.AccessToken)
}
