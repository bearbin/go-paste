// Package pastebin wraps the basic functions of the Pastebin API and exposes a
// Go API.
package pastebin

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	pastebinDevKey = "d06a9df64b29123b8eeda23f53d6535d"
)

var (
	ErrPutFailed = errors.New("Pastebin Put Failed!")
	ErrGetFailed = errors.New("Pastebin Get Failed!")
)

type Pastebin struct{}

// Function Put uploads text to Pastebin with optional title returning the ID or
// an error.
func (p Pastebin) Put(text, title string) (id string, err error) {
	data := url.Values{}
	// Required values.
	data.Set("api_dev_key", pastebinDevKey)
	data.Set("api_option", "paste") // Create a paste.
	data.Set("api_paste_code", text)
	// Optional values.
	data.Set("api_paste_name", title)      // The paste should have title "title".
	data.Set("api_paste_private", "0")     // Create a public paste.
	data.Set("api_paste_expire_date", "N") // The paste should never expire.

	resp, err := http.PostForm("http://pastebin.com/api/api_post.php", data)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", ErrPutFailed
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return p.StripURL(string(respBody)), nil
}

// Function Get returns the text inside the paste identified by ID.
func (p Pastebin) Get(id string) (text string, err error) {
	resp, err := http.Get("http://pastebin.com/raw.php?i=" + id)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", ErrGetFailed
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

// Function StripURL returns the paste ID from a pastebin URL.
func (p Pastebin) StripURL(url string) string {
	return strings.Replace(url, "http://pastebin.com/", "", -1)
}

// Function WrapID returns the pastebin URL from a paste ID.
func (p Pastebin) WrapID(id string) string {
	return "http://pastebin.com/" + id
}
