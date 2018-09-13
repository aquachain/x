package main

import (
	"bytes"
	"net/http"
	"time"
)

const CacheDuration = time.Second * 30

type Cache struct {
	b bytes.Buffer
	t time.Time
}

func NewCache() *Cache {
	return &Cache{bytes.Buffer{}, time.Time{}}
}

func (h *Handler) cached(w http.ResponseWriter, r *http.Request) (handled bool) {
	if cache := h.cache[r.URL.Path]; cache != nil && time.Since(cache.t) < CacheDuration {
		w.Write(cache.b.Bytes())
		return true
	}
	if h.cache[r.URL.Path] == nil {
		h.cache[r.URL.Path] = NewCache()
	}
	h.cache[r.URL.Path].t = time.Now()
	h.cache[r.URL.Path].b.Reset()
	return false
}
