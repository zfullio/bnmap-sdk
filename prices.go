package bnmap

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Client) GetPriceLists(ctx context.Context) (RespPrices, error) {
	q := ParamsBuilder(PbiPrices, c.token, 1)

	resp, err := SendRequest(ctx, q)
	if err != nil {
		return RespPrices{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
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

	t, err := dec.Token()
	if err != nil {
		return DecodeTokenError
	}

	if t != json.Delim('{') {
		return errors.New("expected '{', got " + fmt.Sprint(t))
	}

	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return DecodeTokenError
		}

		switch t {
		case "status":
			err = dec.Decode(&r.Status)
			if err != nil {
				return &DecodeError{msg: fmt.Sprintf("field 'status' failed to decode: %v", err)}
			}
		case "auth":
			err = dec.Decode(&r.Auth)
			if err != nil {
				return &DecodeError{msg: fmt.Sprintf("field 'auth' failed to decode: %v", err)}
			}
		case "content":
			_, err := dec.Token()
			if err != nil {
				return DecodeTokenError
			}

			r.Content = make([]map[string]string, 0, 50)

			for dec.More() {
				rawRow := make(map[string]interface{}, 0)

				convRow := make(map[string]string, 50)

				err = dec.Decode(&rawRow)
				if err != nil {
					return &DecodeError{msg: err.Error()}
				}

				for k, v := range rawRow {
					if v == nil {
						continue
					}

					switch v.(type) {
					case string:
						val, ok := v.(string)
						if !ok {
							return &DecodeError{msg: fmt.Sprintf("%v", v)}
						}
						convRow[k] = val
					case int64:
						val, ok := v.(int64)
						if !ok {
							return &DecodeError{msg: fmt.Sprintf("%v", v)}
						}
						convRow[k] = fmt.Sprintf("%d", val)
					case float64:
						val, ok := v.(float64)
						if !ok {
							return &DecodeError{msg: fmt.Sprintf("%v", v)}
						}

						convRow[k] = fmt.Sprintf("%f", val)
					default:
						convRow[k] = fmt.Sprintf("%v", v)
					}

					r.Content = append(r.Content, convRow)
				}
			}

			_, err = dec.Token()
			if err != nil {
				return DecodeTokenError
			}
		}
	}

	return nil
}
