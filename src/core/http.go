package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ============================================================================

var http_client = &http.Client{
	Timeout: time.Second * 5,
}

// ============================================================================

func HttpGet(addr string) (ret string) {
	res, err := http_client.Get(addr)
	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return string(body)
}

func HttpPost(addr string, data url.Values) (ret string) {
	res, err := http_client.PostForm(addr, data)
	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return string(body)
}

func HttpPostJson(addr string, J string) (ret string) {
	res, err := http_client.Post(addr, "application/json", bytes.NewBufferString(J))
	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return string(body)
}
