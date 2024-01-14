package net

type Client interface {
	Send(target string, path string) error
}
