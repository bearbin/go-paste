// Package fpaste wraps the basic functions of the Pastebin API and exposes a
// Go API.
package fpaste

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrPutFailed = errors.New("fpaste Put Failed!")
	ErrGetFailed = errors.New("fpaste Get Failed!")
)

type Fpaste struct{}
type fpasteResponse struct {
	Result struct {
		ID string `json:"id"`
	} `json:"result"`
}

// Function Put uploads text to fpaste.org. It returns the ID of the created
// paste or an error. The title is not used, as the service does not support
// titles.
func (f Fpaste) Put(text, title string) (id string, err error) {
	data := url.Values{}
	// Required values.
	data.Set("paste_data", text)
	data.Set("paste_lang", "text")
	data.Set("api_submit", "true")
	data.Set("mode", "json")        // Get the results back in JSON.
	data.Set("paste_private", "no") // Public paste.
	data.Set("paste_expire", "0")   // Never expire.

	resp, err := http.PostForm("http://fpaste.org", data)
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
	decresp := &fpasteResponse{}
	err = json.Unmarshal(respBody, decresp)
	if err != nil {
		return "", err
	}
	return decresp.Result.ID, nil
}

// Function Get returns the text inside the paste identified by ID.
func (f Fpaste) Get(id string) (text string, err error) {
	resp, err := http.Get("http://fpaste.org/" + id + "/raw/")
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

// Function StripURL returns the paste ID from a fpaste URL.
func (f Fpaste) StripURL(url string) string {
	return strings.Replace(url, "http://fpaste.org/", "", -1)
}

// Function WrapID returns the fpaste URL from a paste ID.
func (f Fpaste) WrapID(id string) string {
	return "http://fpaste.org/" + id
}
