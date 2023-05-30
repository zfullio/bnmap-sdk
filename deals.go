package bnmap

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetFullDeals(ctx context.Context, page int64) (RespFullDeals, error) {
	q := ParamsBuilder(PbiFullDeals, c.token, page)

	resp, err := SendRequest(ctx, q)
	if err != nil {
		return RespFullDeals{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RespFullDeals{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var data RespFullDeals

	err = data.Decode(bufio.NewReader(resp.Body))
	if err != nil {
		return data, err
	}

	return data, nil
}

type RespFullDeals struct {
	Status  string           `json:"status"`
	Content FullDealsContent `json:"content"`
	Auth    bool             `json:"auth"`
}

func (r *RespFullDeals) Decode(reader *bufio.Reader) error {
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

			if delim, ok := contentToken.(json.Delim); !ok || delim != '{' {
				return fmt.Errorf("expected {, got %v", contentToken)
			}

			if _, err = dec.Token(); err != nil {
				return ErrDecodeToken
			}

			for dec.More() {
				contentToken, err := dec.Token()
				if err != nil {
					return ErrDecodeToken
				}

				switch contentToken {
				case "page":
					if err := dec.Decode(&r.Content.Page); err != nil {
						return &DecodeError{msg: err.Error()}
					}
				case "per_page":
					if err := dec.Decode(&r.Content.PerPage); err != nil {
						return &DecodeError{msg: err.Error()}
					}
				case "total_pages":
					if err := dec.Decode(&r.Content.TotalPages); err != nil {
						return &DecodeError{msg: err.Error()}
					}
				case "total":
					if err := dec.Decode(&r.Content.Total); err != nil {
						return &DecodeError{msg: err.Error()}
					}
				case "data":
					dataToken, err := dec.Token()
					if err != nil {
						return ErrDecodeToken
					}

					if delim, ok := dataToken.(json.Delim); !ok || delim != '[' {
						return fmt.Errorf("expected [, got %v", dataToken)
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

						r.Content.Data = append(r.Content.Data, row)
					}

					_, err = dec.Token()
					if err != nil {
						return ErrDecodeToken
					}
				}
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

type FullDealsContent struct {
	Page       int64               `json:"page"`
	PerPage    int64               `json:"per_page"`
	Total      string              `json:"total"`
	TotalPages int64               `json:"total_pages"`
	Data       []map[string]string `json:"data"`
}

type Row []map[string]string

type Field map[string]string

type Deal struct {
	BStartSales  string  `json:"b_start_sales"`
	BtName       string  `json:"bt_name"`
	Builder      string  `json:"builder"`
	Class        string  `json:"class"`
	Concession   string  `json:"concession"`
	DealsSeller  string  `json:"deals_seller"`
	Developer    string  `json:"developer"`
	DocumentDate string  `json:"document_date"`
	DoSquare     float64 `json:"do_square"`
	EstBudget    int     `json:"est_budget"`
	Floor        int64   `json:"floor"`
	HcName       string  `json:"hc_name"`
	ID           string  `json:"id"`
	LocAddress   string  `json:"loc_address"`
	LocArea      string  `json:"loc_area"`
	LocDistrict  string  `json:"loc_district"`
	Mortgage     string  `json:"mortgage"`
	MortgageTerm int     `json:"mortgage_term"`
	ObjectID     string  `json:"object_id"`
	OtName       string  `json:"ot_name"`
	PboNumber    string  `json:"pbo_number"`
	PriceSquareR float64 `json:"price_square_r"`
	RegDate      string  `json:"reg_date"`
	RegionName   string  `json:"region_name"`
	Rooms        string  `json:"rooms"`
	RoomsPricesT string  `json:"rooms_prices_t"`
	Section      int64   `json:"section"`
	Square       float64 `json:"square"`
	Wholesale    string  `json:"wholesale"`
}
