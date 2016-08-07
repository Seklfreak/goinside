package goinside

import (
	"io"
	"net/http"
)

// web urls
const (
	gallerysURL    = "http://m.dcinside.com/category_gall_total.html"
	commentMoreURL = "http://m.dcinside.com/comment_more_new.php"
)

// api urls
const (
	appID            = "blM1T09mWjRhQXlZbE1ML21xbkM3QT09"
	loginAPI         = "https://dcid.dcinside.com/join/mobile_app_login.php"
	articleWriteAPI  = "http://upload.dcinside.com/_app_write_api.php"
	articleDeleteAPI = "http://m.dcinside.com/api/gall_del.php"
	commentWriteAPI  = "http://m.dcinside.com/api/comment_ok.php"
	commentDeleteAPI = "http://m.dcinside.com/api/comment_del.php"
	recommendUpAPI   = "http://m.dcinside.com/api/_recommend_up.php"
	recommendDownAPI = "http://m.dcinside.com/api/_recommend_down.php"
	reportAPI        = "http://m.dcinside.com/api/report_upload.php"
)

// content types
const (
	defaultContentType    = "application/x-www-form-urlencoded; charset=UTF-8"
	nonCharsetContentType = "application/x-www-form-urlencoded"
)

var (
	apiRequestHeader = map[string]string{
		"User-Agent": "dcinside.app",
		"Referer":    "http://m.dcinside.com",
		"Host":       "m.dcinside.com",
	}
	mobileRequestHeader = map[string]string{
		"User-Agent": "Linux Android",
		"Referer":    "http://m.dcinside.com",
	}
)

type connector interface {
	connection() *Connection
}

func post(c connector, URL string, cookies []*http.Cookie, form io.Reader, contentType string) (*http.Response, error) {
	return do(c, "POST", URL, cookies, form, contentType, mobileRequestHeader)
}

func get(c connector, URL string) (*http.Response, error) {
	return do(c, "GET", URL, nil, nil, defaultContentType, mobileRequestHeader)
}

func api(c connector, URL string, form io.Reader, contentType string) (*http.Response, error) {
	return do(c, "POST", URL, nil, form, contentType, apiRequestHeader)
}

func do(c connector, method, URL string, cookies []*http.Cookie, form io.Reader, contentType string, requestHeader map[string]string) (*http.Response, error) {
	URL = _MobileURL(URL)
	req, err := http.NewRequest(method, URL, form)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range requestHeader {
		req.Header.Set(k, v)
	}
	client := func() *http.Client {
		proxy := c.connection().proxy
		if proxy != nil {
			return &http.Client{Transport: &http.Transport{Proxy: proxy}}
		}
		return &http.Client{}
	}()
	if c.connection().timeout != 0 {
		client.Timeout = c.connection().timeout
	}
	return client.Do(req)
}
