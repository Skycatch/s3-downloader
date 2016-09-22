package main

// DownloadHandler hello world
type DownloadHandler interface {
	id() *string
	body() *string
	initialize()
	receive() bool
	success()
}
