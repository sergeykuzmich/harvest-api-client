package harvest_api_client

import (
	"testing"
)

func TestHarvestClient(t *testing.T) {
	a1 := HarvestClient("ACCOUNTID", "TOKEN")
	if a1.BaseURL != "https://api.harvestapp.com/v2" {
		t.Errorf("Incorrect domain name '%s'.", a1.BaseURL)
	}
	if a1.client == nil {
		t.Error("No http client")
	}
	if a1.AccountID != "ACCOUNTID" {
		t.Error("AccountID not assigned correctly")
	}
	if a1.AccessToken != "TOKEN" {
		t.Error("AccessToken not assigned correctly")
	}
}

func testAPI() *API {
	a := HarvestClient("ACCOUNTID", "TOKEN")
	return a
}
