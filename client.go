package bnmap

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const BnHost = "api.bndev.it"

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func ParamsBuilder(method Method, token string, page int64) url.Values {
	params := url.Values{}
	params.Add("act", string(method))
	params.Add("pbi", token)

	if page > 0 {
		params.Add("page", strconv.FormatInt(page, 10))
	}

	return params
}

func SendRequest(ctx context.Context, params url.Values) (*http.Response, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     BnHost,
		Path:     "cmap/analytics.json",
		RawQuery: params.Encode(),
	}
	reqURL := u.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	cl := http.Client{}

	response, reqErr := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", reqErr)
	}

	return response, nil
}
