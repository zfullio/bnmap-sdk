package bnmap

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetPriceLists(ctx context.Context) (RespPrices, error) {
	q := ParamsBuilder(PbiPrices, c.token, 1)

	resp, err := SendRequest(ctx, q)
	if err != nil {
		return RespPrices{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RespPrices{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var data RespPrices

	err = data.Decode(bufio.NewReader(resp.Body))
	if err != nil {
		return RespPrices{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return data, nil
}

type RespPrices struct {
	Status  string              `json:"status"`
	Content []map[string]string `json:"content"`
	Auth    bool                `json:"auth"`
}

func (r *RespPrices) Decode(reader *bufio.Reader) error {
	dec := json.NewDecoder(reader)

	var token json.Token

	var err error

	if token, err = dec.Token(); err != nil {
		return ErrDecodeToken
	}

	if json.Delim('{') != token {
		return fmt.Errorf("expected {, got %v", token)
	}

	for dec.More() {
		switch token {
		case "status":
			if err := dec.Decode(&r.Status); err != nil {
				return &DecodeError{msg: err.Error()}
			}
		case "auth":
			if err := dec.Decode(&r.Auth); err != nil {
				return &DecodeError{msg: err.Error()}
			}
		case "content":
			var contentToken json.Token

			if contentToken, err = dec.Token(); err != nil {
				return ErrDecodeToken
			}

			if delim, ok := contentToken.(json.Delim); !ok || delim != '[' {
				return fmt.Errorf("expected [, got %v", contentToken)
			}

			for dec.More() {
				row := make(map[string]string, 0)

				var raw map[string]interface{}

				var conv string

				err = dec.Decode(&raw)
				if err != nil {
					return &DecodeError{msg: err.Error()}
				}

				for k, val := range raw {
					if val == nil {
						continue
					}

					switch val := val.(type) {
					case string:
						if val == "" {
							continue
						}

						conv = val
					case int:
						conv = fmt.Sprintf("%d", val)
					case float64:
						conv = fmt.Sprintf("%f", val)
					default:
						conv = fmt.Sprintf("%v", val)
					}

					row[k] = conv
				}

				r.Content = append(r.Content, row)
			}

			_, err = dec.Token()
			if err != nil {
				return ErrDecodeToken
			}
		}

		if token, err = dec.Token(); err != nil {
			return ErrDecodeToken
		}
	}

	return nil
}
