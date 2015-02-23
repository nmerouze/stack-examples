package main

import (
	"net/http"

	"../../upload"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func main() {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}

	server := s3.New(auth, aws.USEast)
	b := server.Bucket("foobar.com")
	h := upload.Service(b)
	http.ListenAndServe(":8080", h)
}
