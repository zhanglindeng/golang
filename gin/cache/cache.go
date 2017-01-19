package cache

import (
	"bytes"
	"crypto/sha1"
	"io"
	"net/http"
	"net/url"
	"time"
	"log"

	"dengzhanglin/gin/cache/persistence"
	"gopkg.in/gin-gonic/gin.v1"
)

const (
	CACHE_MIDDLEWARE_KEY = "gincontrib.cache"
)

var (
	PageCachePrefix = "gincontrib.page.cache"
)

type responseCache struct {
	status int
	header http.Header
	data   []byte
}

type cachedWriter struct {
	gin.ResponseWriter
	status  int
	written bool
	store   persistence.CacheStore
	expire  time.Duration
	key     string
}

func urlEscape(prefix string, u string) string {
	key := url.QueryEscape(u)
	if len(key) > 200 {
		h := sha1.New()
		io.WriteString(h, u)
		key = string(h.Sum(nil))
	}
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}

func newCachedWriter(store persistence.CacheStore, expire time.Duration, writer gin.ResponseWriter, key string) *cachedWriter {
	return &cachedWriter{writer, 0, false, store, expire, key}
}

func (w *cachedWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachedWriter) Status() int {
	return w.status
}

func (w *cachedWriter) Written() bool {
	return w.written
}

func (w *cachedWriter) Write(data []byte) (int, error) {
	log.Println("hello")
	ret, err := w.ResponseWriter.Write(data)
	if err == nil {
		log.Println("nil")
		//cache response
		store := w.store
		val := responseCache{
			w.status,
			w.Header(),
			data,
		}
		err = store.Set(w.key, val, w.expire)
		if err != nil {
			// need logger
			log.Println("err", err)
		}
	}
	return ret, err
}

func (w *cachedWriter) WriteString(str string) (int, error) {
	log.Println("hello")
	ret, err := w.ResponseWriter.WriteString(str)
	if err == nil {
		log.Println("nil")
		//cache response
		store := w.store
		val := responseCache{
			w.status,
			w.Header(),
			[]byte(str),
		}
		err = store.Set(w.key, val, w.expire)
		if err != nil {
			// need logger
			log.Println("err", err)
		}
	}
	return ret, err
}

// Cache Middleware
func Cache(store *persistence.CacheStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(CACHE_MIDDLEWARE_KEY, store)
		c.Next()
	}
}

func SiteCache(store persistence.CacheStore, expire time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		var cache responseCache
		url1 := c.Request.URL
		key := urlEscape(PageCachePrefix, url1.RequestURI())
		if err := store.Get(key, &cache); err != nil {
			c.Next()
		} else {
			c.Writer.WriteHeader(cache.status)
			for k, vals := range cache.header {
				for _, v := range vals {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Write(cache.data)
		}
	}
}

// Cache Decorator
func CachePage(store persistence.CacheStore, expire time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		var cache responseCache
		url1 := c.Request.URL
		key := urlEscape(PageCachePrefix, url1.RequestURI())
		if err := store.Get(key, &cache); err != nil {
			// replace writer
			writer := newCachedWriter(store, expire, c.Writer, key)
			c.Writer = writer
			handle(c)
			log.Println("abc")
		} else {
			c.Writer.WriteHeader(cache.status)
			for k, vals := range cache.header {
				for _, v := range vals {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Write(cache.data)
		}
	}
}
