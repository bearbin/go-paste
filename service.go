package main

type service interface {
	Put(text, title string) (id string, err error)
	Get(id string) (text string, err error)
	// StripURL strips the URL to produce the ID of the paste.
	StripURL(url string) string
	// WrapURL wraps an ID to produce a usable URL.
	WrapID(id string) string
}
