// Package pastebin wraps the basic functions of the Pastebin API and exposes a
// Go API.
package pastebin

import (
	"bytes"
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

// Function Put uploads text to Pastebin with optional title returning the ID or
// an error.
func Put(text, title string) (id string, err error) {
	data := url.Values{}
	// Required values.
	data.Set("api_dev_key", pastebinDevKey)
	data.Set("api_option", "paste") // Create a paste.
	data.Set("api_paste_code", text)
	// Optional values.
	data.Set("api_paste_name", title)      // The paste should have title "title".
	data.Set("api_paste_private", "0")     // Create a public paste.
	data.Set("api_paste_expire_date", "N") // The paste should never expire.

	// Parse and URLEncode the values ready to pass to Pastebin.
	body := bytes.NewBufferString(data.Encode())

	resp, err := http.Post("http://pastebin.com/api/api_post.php", "application/x-www-form-urlencoded", body)
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
	return strings.Replace(string(respBody), "http://pastebin.com/", "", -1), nil
}

// Function Get returns the text inside the paste identified by ID.
func Get(id string) (text string, err error) {
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
func StripURL(url string) string {
	return strings.TrimPrefix(url, "http://pastebin.com/")
}

// Function WrapID returns the pastebin URL from a paste ID.
func WrapID(id string) string {
	return "http://pastebin.com/" + id
}
