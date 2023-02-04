package bnmap

import (
	"encoding/json"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (c *Client) GetPriceLists() (prices []Price, err error) {
	q := ParamsBuilder(PbiPrices, c.token, 1)
	resp, err := SendRequest(q)
	defer resp.Body.Close()
	if err != nil {
		return prices, err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return prices, err
	}
	data := RespPrices{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return prices, err
	}
	return data.Content, nil
}

type RespPrices struct {
	Status  string  `json:"status"`
	Content []Price `json:"content"`
	Auth    bool    `json:"auth"`
}
type Price struct {
	Agreement           string  `json:"agreement"`
	Area                string  `json:"area"`
	Building            string  `json:"building"`
	CarSpacesNum        string  `json:"car_spaces_num"`
	Construction        string  `json:"construction"`
	CreateTime          string  `json:"create_time"`
	Decoration          string  `json:"decoration"`
	DecorationDesc      string  `json:"decoration_desc"`
	Developer           string  `json:"developer"`
	Discount            string  `json:"discount"`
	DiscountDesc        string  `json:"discount_desc"`
	DiscountInPrice     string  `json:"discountInPrice"`
	DiscountNote        string  `json:"discountNote"`
	DiscountValue       string  `json:"discount_value"`
	District            string  `json:"district"`
	Dsc                 string  `json:"dsc"`
	FirstAppearance     string  `json:"first_appearance"`
	FirstLotDate        string  `json:"first_lot_date"`
	FirstLotPrice       float64 `json:"first_lot_price"`
	FirstLotPriceSquare float64 `json:"first_lot_price_square"`
	Floor               int64   `json:"floor"`
	FloorsFrom          string  `json:"floors_from"`
	FloorsTo            string  `json:"floors_to"`
	HcName              string  `json:"hc_name"`
	ID                  string  `json:"id"`
	Interior            string  `json:"interior"`
	LivingSquare        float64 `json:"living_square"`
	LotPriceRise        float64 `json:"lot_price_rise"`
	LotPriceSquareRise  float64 `json:"lot_price_square_rise"`
	LotTerm             string  `json:"lot_term"`
	NumApartment        string  `json:"numApartment"`
	NumInPlatform       int64   `json:"numInPlatform"`
	ObjectClass         string  `json:"object_class"`
	ObjectType          string  `json:"object_type"`
	ParkingType         string  `json:"parking_type"`
	Price               float64 `json:"price"`
	Region              string  `json:"region"`
	Rooms               string  `json:"rooms"`
	Section             int     `json:"section"`
	SourceURL           string  `json:"source_url"`
	Square              float64 `json:"square"`
	SquarePrice         float64 `json:"squarePrice"`
	Stage               string  `json:"stage"`
	StageDesc           string  `json:"stage_desc"`
	StartSales          string  `json:"start_sales"`
}

func (c *Client) GetFullDeals() (deals []Deal, err error) {
	q := ParamsBuilder(PbiFullDeals, c.token, 1)
	data := RespFullDeals{}
	err = PrepareResponse(q, &data)
	if err != nil {
		return deals, err
	}
	deals = append(deals, data.Content.Data...)
	pages := data.Content.TotalPages
	bar := progressbar.Default(pages, "Обработка страниц")
	err = bar.Add(1)
	if err != nil {
		log.Printf("Ошибка прогрессбара: %s", err)
	}
	//TODO Исправить page <= int(pages)
	for page := 2; page <= int(pages); page++ {
		err := bar.Add(1)
		if err != nil {
			log.Printf("Ошибка прогрессбара: %s", err)
		}
		q := ParamsBuilder(PbiFullDeals, c.token, int64(page))
		data := RespFullDeals{}
		err = PrepareResponse(q, &data)
		if err != nil {
			return deals, err
		}
		deals = append(deals, data.Content.Data...)
	}
	return deals, nil
}

type RespFullDeals struct {
	Status  string `json:"status"`
	Content struct {
		Page       int64  `json:"page"`
		PerPage    int64  `json:"per_page"`
		Total      string `json:"total"`
		TotalPages int64  `json:"total_pages"`
		Data       []Deal `json:"data"`
	} `json:"content"`
	Auth bool `json:"auth"`
}
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

func (c *Client) GetFullObjects() (objects []Object, err error) {
	q := ParamsBuilder(PbiFullObjects, c.token, 1)
	data := RespFullObjects{}
	err = PrepareResponse(q, &data)
	if err != nil {
		return objects, err
	}
	objects = append(objects, data.Content.Data...)
	pages := data.Content.TotalPages
	bar := progressbar.Default(int64(pages), "Обработка страниц")
	err = bar.Add(1)
	if err != nil {
		log.Printf("Ошибка прогрессбара: %s", err)
	}
	for page := 2; page <= pages; page++ {
		err := bar.Add(1)
		if err != nil {
			log.Printf("Ошибка прогрессбара: %s", err)
		}
		q := ParamsBuilder(PbiFullObjects, c.token, int64(page))
		data := RespFullObjects{}
		err = PrepareResponse(q, &data)
		if err != nil {
			return objects, err
		}
		objects = append(objects, data.Content.Data...)
		time.Sleep(2 * time.Second)
	}
	return objects, nil
}

type RespFullObjects struct {
	Status  string `json:"status"`
	Content struct {
		Page       int      `json:"page,string"`
		PerPage    int      `json:"per_page"`
		Total      int      `json:"total"`
		TotalPages int      `json:"total_pages"`
		Data       []Object `json:"data"`
	} `json:"content"`
	Auth bool `json:"auth"`
}
type Object struct {
	ID            string  `json:"id"`
	HcName        string  `json:"hc_name"`
	Address       string  `json:"address"`
	SaleStage     string  `json:"salestage"`
	Section       int64   `json:"section"`
	PboNumber     string  `json:"pbo_number"`
	OtName        string  `json:"ot_name"`
	DotName       string  `json:"dot_name"`
	Floor         int64   `json:"floor"`
	Rooms         string  `json:"rooms"`
	BRooms        string  `json:"b_rooms"`
	Square        float64 `json:"square"`
	IsDeals       string  `json:"is_deals"`
	AgreementDate string  `json:"agreement_date"`
	WholeSale     string  `json:"wholesale"`
	Concession    string  `json:"concession"`
	PdsType       string  `json:"pds_type"`
	BtName        string  `json:"bt_name"`
	DBFullName    string  `json:"db_fullname"`
}

func (c *Client) GetPriceLayers() (corpses []CorpusLayer, err error) {
	q := ParamsBuilder(Table, c.token, 1)
	resp, err := SendRequest(q)
	if err != nil {
		return corpses, err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return corpses, err
	}
	data := RespCorpusLayers{}
	if err != nil {
		return corpses, err
	}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return corpses, err
	}
	return data.Content.Data, nil
}

type RespCorpusLayers struct {
	Status  string `json:"status"`
	Content struct {
		Data []CorpusLayer `json:"data"`
	} `json:"content"`
	Auth bool `json:"auth"`
}
type CorpusLayer struct {
	Project             string      `json:"project"`
	Address             string      `json:"address"`
	Coords              string      `json:"coords"`
	Area                string      `json:"area"`
	District            string      `json:"district"`
	Region              string      `json:"region"`
	City                string      `json:"city"`
	Metro               string      `json:"metro"`
	Apartments          string      `json:"apartments"`
	Developer           string      `json:"developer"`
	Builder             string      `json:"builder"`
	TotalSquare         float64     `json:"total_square"`
	LivingSquare        float64     `json:"living_square"`
	Agreement           string      `json:"agreement"`
	Interior            string      `json:"interior"`
	FloorsFrom          string      `json:"floors_from"`
	FloorsTo            string      `json:"floors_to"`
	FloorType           string      `json:"floor_type"`
	Seller              string      `json:"seller"`
	StartSalesDate      string      `json:"start_sales_date"`
	InitialDsc          string      `json:"initial_dsc"`
	DscCount            string      `json:"dsc_count"`
	DateStateCommission string      `json:"date_state_commission"`
	Stage               string      `json:"stage"`
	StageAlias          string      `json:"stage_alias"`
	Class               string      `json:"class"`
	Construction        string      `json:"construction"`
	ProjectFrom         string      `json:"project_from"`
	LayoutFrom          string      `json:"layout_from"`
	ParkingType         string      `json:"parking_type"`
	ParkingProject      int         `json:"parking_project"`
	ParkingExpo         interface{} `json:"parking_expo"`
	CarSpacesPriceMin   int         `json:"carSpacesPriceMin"`
	CarSpacesPriceMax   int         `json:"carSpacesPriceMax"`
	CarSpacesSquareMin  float64     `json:"carSpacesSquareMin"`
	CarSpacesSquareMax  float64     `json:"carSpacesSquareMax"`
	CarSpAvg            float64     `json:"car_sp_avg"`
	MetrPriceRMin       float64     `json:"metrPriceRMin"`
	MetrPriceRAvg       float64     `json:"metrPriceRAvg"`
	MetrPriceRMax       float64     `json:"metrPriceRMax"`
	Discount            string      `json:"discount"`
	DiscountDesc        string      `json:"discount_desc"`
	Available           string      `json:"available"`
	Deals               bool        `json:"deals"`
	Price               bool        `json:"price"`
	Pantry              bool        `json:"pantry"`
	CarSpace            bool        `json:"carspace"`
	Commercial          bool        `json:"commercial"`
	PboNum              string      `json:"pbo_num"`
	MonthsInSales       int         `json:"months_in_sales"`
	MonthsBeforeDsc     int         `json:"months_before_dsc"`
	ApartTotal          Apart       `json:"apart_total"`
	ApartSt             Apart       `json:"apart_st"`
	Apart1              Apart       `json:"apart_1"`
	Apart2              Apart       `json:"apart_2"`
	Apart3              Apart       `json:"apart_3"`
	Apart4              Apart       `json:"apart_4"`
	Apart0              Apart       `json:"apart_0"`
	B                   struct {
		ApartSt      ApartSquare `json:"apart_st"`
		Apart1       ApartSquare `json:"apart_1"`
		Apart2       ApartSquare `json:"apart_2"`
		Apart3       ApartSquare `json:"apart_3"`
		Apart4       ApartSquare `json:"apart_4"`
		Apart0       ApartSquare `json:"apart_0"`
		LivingSquare float64     `json:"living_square"`
		RoomsNum     int         `json:"rooms_num"`
	} `json:"b"`
	Ds struct {
		Amount              string      `json:"amount"`
		Flats               string      `json:"flats"`
		Apart               string      `json:"apart"`
		CarSpaces           string      `json:"car_spaces"`
		NonResidential      string      `json:"non_residential"`
		Pantry              string      `json:"pantry"`
		FlatsFl             string      `json:"flats_fl"`
		FlatsUl             string      `json:"flats_ul"`
		ApartFl             string      `json:"apart_fl"`
		ApartUl             string      `json:"apart_ul"`
		FlatsFlMortgage     string      `json:"flats_fl_mortgage"`
		FlatsFlNonMortgage  string      `json:"flats_fl_non_mortgage"`
		ApartFlMortgage     string      `json:"apart_fl_mortgage"`
		ApartFlNonMortgage  string      `json:"apart_fl_non_mortgage"`
		PaceFlats           interface{} `json:"pace_flats"`
		PaceFlatsFl         interface{} `json:"pace_flats_fl"`
		PaceFlatsUl         interface{} `json:"pace_flats_ul"`
		PaceApart           interface{} `json:"pace_apart"`
		PaceApartFl         interface{} `json:"pace_apart_fl"`
		PaceApartUl         interface{} `json:"pace_apart_ul"`
		PaceFlatsPre1       interface{} `json:"pace_flats_pre_1"`
		PaceFlatsFlPre1     interface{} `json:"pace_flats_fl_pre_1"`
		PaceFlatsUlPre1     interface{} `json:"pace_flats_ul_pre_1"`
		PaceApartPre1       interface{} `json:"pace_apart_pre_1"`
		PaceApartFlPre1     interface{} `json:"pace_apart_fl_pre_1"`
		PaceApartUlPre1     interface{} `json:"pace_apart_ul_pre_1"`
		PaceFlatsPre3       interface{} `json:"pace_flats_pre_3"`
		PaceFlatsFlPre3     interface{} `json:"pace_flats_fl_pre_3"`
		PaceFlatsUlPre3     interface{} `json:"pace_flats_ul_pre_3"`
		PaceApartPre3       interface{} `json:"pace_apart_pre_3"`
		PaceApartFlPre3     interface{} `json:"pace_apart_fl_pre_3"`
		PaceApartUlPre3     interface{} `json:"pace_apart_ul_pre_3"`
		PaceFlatsPre6       interface{} `json:"pace_flats_pre_6"`
		PaceFlatsFlPre6     interface{} `json:"pace_flats_fl_pre_6"`
		PaceFlatsUlPre6     interface{} `json:"pace_flats_ul_pre_6"`
		PaceApartPre6       interface{} `json:"pace_apart_pre_6"`
		PaceApartFlPre6     interface{} `json:"pace_apart_fl_pre_6"`
		PaceApartUlPre6     interface{} `json:"pace_apart_ul_pre_6"`
		PaceFlatsPre12      interface{} `json:"pace_flats_pre_12"`
		PaceFlatsFlPre12    interface{} `json:"pace_flats_fl_pre_12"`
		PaceFlatsUlPre12    interface{} `json:"pace_flats_ul_pre_12"`
		PaceApartPre12      interface{} `json:"pace_apart_pre_12"`
		PaceApartFlPre12    interface{} `json:"pace_apart_fl_pre_12"`
		PaceApartUlPre12    interface{} `json:"pace_apart_ul_pre_12"`
		FlatsFlSquare       string      `json:"flats_fl_square"`
		ApartFlSquare       string      `json:"apart_fl_square"`
		FlatsFlSquareAvg    string      `json:"flats_fl_square_avg"`
		ApartFlSquareAvg    string      `json:"apart_fl_square_avg"`
		FlatsFlMetrPriceAvg string      `json:"flats_fl_metrprice_avg"`
		ApartFlMetrPriceAvg string      `json:"apart_fl_metrprice_avg"`
		FlatsFlSumAvg       string      `json:"flats_fl_sum_avg"`
		ApartFlSumAvg       string      `json:"apart_fl_sum_avg"`
	} `json:"ds"`
	Ub struct {
		UnrealizedAmount       string      `json:"unrealized_amount"`
		UnrealizedSquare       string      `json:"unrealized_square"`
		AvgAnnualAmount        string      `json:"avg_annual_amount"`
		AvgAnnualSquare        string      `json:"avg_annual_square"`
		AbsorptionByLotAmount1 interface{} `json:"absorption_by_lot_amount_1"`
		AbsorptionBySquare1    interface{} `json:"absorption_by_square_1"`
		AbsorptionByLotAmount2 interface{} `json:"absorption_by_lot_amount_2"`
		AbsorptionBySquare2    interface{} `json:"absorption_by_square_2"`
		UnrealizedAmount1      string      `json:"unrealized_amount_1"`
		UnrealizedSquare1      string      `json:"unrealized_square_1"`
		AvgAnnualAmount1       string      `json:"avg_annual_amount_1"`
		AvgAnnualSquare1       string      `json:"avg_annual_square_1"`
		UnrealizedAmount2      string      `json:"unrealized_amount_2"`
		UnrealizedSquare2      string      `json:"unrealized_square_2"`
		AvgAnnualAmount2       string      `json:"avg_annual_amount_2"`
		AvgAnnualSquare2       string      `json:"avg_annual_square_2"`
		UnrealizedAmount3      string      `json:"unrealized_amount_3"`
		UnrealizedSquare3      string      `json:"unrealized_square_3"`
		AvgAnnualAmount3       string      `json:"avg_annual_amount_3"`
		AvgAnnualSquare3       string      `json:"avg_annual_square_3"`
		UnrealizedAmount4      string      `json:"unrealized_amount_4"`
		UnrealizedSquare4      string      `json:"unrealized_square_4"`
		AvgAnnualAmount4       string      `json:"avg_annual_amount_4"`
		AvgAnnualSquare4       string      `json:"avg_annual_square_4"`
		UnrealizedAmount0      string      `json:"unrealized_amount_0"`
		UnrealizedSquare0      string      `json:"unrealized_square_0"`
		AvgAnnualAmount0       string      `json:"avg_annual_amount_0"`
		AvgAnnualSquare0       string      `json:"avg_annual_square_0"`
		UnrealizedAmountSt     string      `json:"unrealized_amount_st"`
		UnrealizedSquareSt     string      `json:"unrealized_square_st"`
		AvgAnnualAmountSt      string      `json:"avg_annual_amount_st"`
		AvgAnnualSquareSt      string      `json:"avg_annual_square_st"`
		ConversionFactor       string      `json:"conversion_factor"`
		ConversionFactor1      string      `json:"conversion_factor_1"`
		ConversionFactor2      string      `json:"conversion_factor_2"`
		ConversionFactor3      string      `json:"conversion_factor_3"`
		ConversionFactor4      string      `json:"conversion_factor_4"`
		ConversionFactor0      string      `json:"conversion_factor_0"`
		ConversionFactorSt     string      `json:"conversion_factor_st"`
	} `json:"ub"`
	Ps struct {
		MetrPrice   string      `json:"metr_price"`
		MetrPriceSt interface{} `json:"metr_price_st"`
		MetrPrice1  string      `json:"metr_price_1"`
		MetrPrice2  string      `json:"metr_price_2"`
		MetrPrice3  string      `json:"metr_price_3"`
		MetrPrice4  interface{} `json:"metr_price_4"`
		MetrPrice0  interface{} `json:"metr_price_0"`
		Sum         string      `json:"sum"`
		SumSt       interface{} `json:"sum_st"`
		Sum1        string      `json:"sum_1"`
		Sum2        string      `json:"sum_2"`
		Sum3        string      `json:"sum_3"`
		Sum4        interface{} `json:"sum_4"`
		Sum0        interface{} `json:"sum_0"`
	} `json:"ps"`
	IP struct {
		IncPrice   string `json:"inc_price"`
		IncPriceSt string `json:"inc_price_st"`
		IncPrice1  string `json:"inc_price_1"`
		IncPrice2  string `json:"inc_price_2"`
		IncPrice3  string `json:"inc_price_3"`
		IncPrice4  string `json:"inc_price_4"`
		IncPrice0  string `json:"inc_price_0"`
		IncSum     string `json:"inc_sum"`
		IncSumSt   string `json:"inc_sum_st"`
		IncSum1    string `json:"inc_sum_1"`
		IncSum2    string `json:"inc_sum_2"`
		IncSum3    string `json:"inc_sum_3"`
		IncSum4    string `json:"inc_sum_4"`
		IncSum0    string `json:"inc_sum_0"`
	} `json:"ip"`
	DTRowID     string `json:"DT_RowId"`
	LayerID     string `json:"layer_id"`
	Expert      bool   `json:"expert"`
	LayerStatus string `json:"layer_status"`
}
type Apart struct {
	Expo          int     `json:"expo"`
	SquareAll     float64 `json:"squareAll"`
	SquareMin     float64 `json:"squareMin"`
	SquareAvg     float64 `json:"squareAvg"`
	SquareMax     float64 `json:"squareMax"`
	MetrPriceRMin float64 `json:"metrPriceRMin"`
	MetrPriceRAvg float64 `json:"metrPriceRAvg"`
	MetrPriceRMax float64 `json:"metrPriceRMax"`
	SumRmin       float64 `json:"sumRmin"`
	SumRavg       float64 `json:"sumRavg"`
	SumRmax       float64 `json:"sumRmax"`
}
type ApartSquare struct {
	Rooms     string  `json:"rooms"`
	SquareMin float64 `json:"square_min"`
	SquareAvg float64 `json:"square_avg"`
	SquareMax float64 `json:"square_max"`
}
type Column struct {
	Data     string `json:"data"`
	Name     string `json:"name"`
	Visible  bool   `json:"visible"`
	Title    string `json:"title"`
	Order    string `json:"order"`
	IsSystem string `json:"is_system"`
	Type     string `json:"type"`
}

func ParamsBuilder(method Method, token string, page int64) (params url.Values) {
	params = url.Values{}
	params.Add("act", string(method))
	params.Add("pbi", token)
	if page > 0 {
		params.Add("page", strconv.FormatInt(page, 10))
	}
	return params
}

func SendRequest(params url.Values) (response *http.Response, err error) {
	u := url.URL{
		Scheme:   "https",
		Host:     "api.bndev.it",
		Path:     "cmap/analytics.json",
		RawQuery: params.Encode(),
	}
	reqURL := u.String()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return response, err
	}
	cl := http.Client{}
	response, err = cl.Do(req)
	if err != nil {
		return response, err
	}
	return response, nil
}

func PrepareResponse(q url.Values, template any) (err error) {
	resp, err := SendRequest(q)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseBody, &template)
	if err != nil {
		return err
	}
	return nil
}
