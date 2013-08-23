package web;

import (
	"testing"
	"net/url"
	"fmt"
	"net/http"
	"strings"
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

func TestFormValuesDatasource(t *testing.T) {
	method := "POST"
	urlString := "http://example.com/form/test"
	body := strings.NewReader("name=Inigo+Montoya&fingers=6")

	request, err := http.NewRequest(method, urlString, body)

	// Canary
	if err != nil {
		t.Error("! Error constructing a request.", err)
	}

	ds := new(FormValuesDatasource).Init(request)

	if ds.Value("name").(string) != "Inigo Montoya" {
		t.Error("! Prepare to die.")
	}

	if ds.Value("fingers") != 6 {
		t.Error("! Expected six fingers, but got less.")
	}
}

func TestPathDatasource(t *testing.T) {
	ds := new(PathDatasource).Init("/foo/bar")
	if ds.Value("1") != "foo" {
		t.Error("! Expected value 1 to be 'foo'. Got ", ds.Value("1"))
	}

}
