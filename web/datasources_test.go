package web;

import (
	"testing"
	"net/url"
	"fmt"
)

func TestURLDatasource(t *testing.T) {
	rawurl := "http://user:password@example.com/path?foo=bar#fragment"
	testUrl, err := url.Parse(rawurl)

	if err != nil {
		t.Error("! Unexpected error.", err)
		return
	}
	ds := new(URLDatasource)
	ds.Init(testUrl)

	// Test the string values.
	arr := map[string]string{
		"scheme": "http",
		"Path": "/path",
		"host": "example.com",
		"fragment": "fragment",
	}
	for key, val := range arr {
		if ds.Value(key) != val {
			t.Error(fmt.Sprintf("! Expected '%s', got '%s'", val, ds.Value(key)))
		}
	}

	// Test the Query Values object.
	qvals := ds.Value("Query").(url.Values)
	if qvals.Get("foo") != "bar" {
		t.Error("! Expected to find foo=bar query param. Found ", qvals["foo"])
	}

	// Test the Userinfo object
	uinfo, ok := ds.Value("User").(url.Userinfo)
	if (!ok) {
		t.Error("Expected a Userinfo object.")
	}
	if uinfo.Username() != "user" {
		t.Error("! Expected user name 'user', got ", uinfo.Username)
	}
}

func TestQueryParameterDatasource (t *testing.T) {
	testUrl, err := url.ParseRequestURI("/foo?a=b&c=foo+bar&d=1234&d=5678")
	if err != nil {
		t.Error("! Unexpected URL parse error.")
	}
	ds := new(QueryParameterDatasource).Init(testUrl.Query())

	// Test the string values.
	arr := map[string]string{
		"a": "b",
		"c": "foo bar",
		"d": "5678",
	}
	for key, val := range arr {
		if ds.Value(key) != val {
			t.Error(fmt.Sprintf("! Expected '%s', got '%s'", val, ds.Value(key)))
		}
	}

}
