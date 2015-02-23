package upload

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/goamz/s3"
	"github.com/nmerouze/stack/mux"
)

type appContext struct {
	bucket *s3.Bucket
}

func (c *appContext) upsertFile(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Body must be set", http.StatusBadRequest)
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := mux.Params(r).ByName("path")
	err = c.bucket.Put(path, content, r.Header.Get("Content-Type"), s3.PublicRead)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := c.bucket.URL(path)
	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"url":"%s"}`, url)
}

func (c *appContext) getFile(w http.ResponseWriter, r *http.Request) {
	url := c.bucket.URL(mux.Params(r).ByName("path"))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"url":"%s"}`, url)
}

func (c *appContext) deleteFile(w http.ResponseWriter, r *http.Request) {
	c.bucket.Del(mux.Params(r).ByName("path"))
	w.WriteHeader(204)
}

func Service(bucket *s3.Bucket) http.Handler {
	c := &appContext{bucket}
	m := mux.New()
	m.Get("/files/*path").ThenFunc(c.getFile)
	m.Put("/files/*path").ThenFunc(c.upsertFile)
	m.Delete("/files/*path").ThenFunc(c.deleteFile)
	return m
}
