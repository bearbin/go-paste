package main

type service interface {
	Put(text, title string) (id string, err error)
	Get(id string) (text, title string, err error)
}
