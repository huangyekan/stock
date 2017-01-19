package query

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	//"container/list"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"log"
)

type StockResult struct {
	Status int
	Hq     [][]string
	Code   string
}
type Stock struct {
	Code          string
	Name          string
	Date          string
	Start         string
	End           string
	Low           string
	High          string
	Change        string
	ChangePercent string
	DealCount     string
	DealAmount    string

}

type StockCode struct {
	Code string
	Name string
	TradeLocation string
}

func (stock *Stock) GetData(startDate string, endDate string, stockCode string) ([]Stock, error) {

	req, err := http.NewRequest("GET", stock.GetUrl(stockCode, startDate, endDate), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("http request error")
		return nil, err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("param : " + startDate + " " + endDate + " " + stockCode)
	fmt.Println("result : " + string(result))
	result = []byte(string(result)[strings.Index(string(result), "{"):strings.LastIndex(string(result), "]")])
	stockResult := &StockResult{}
	err = json.Unmarshal(result, stockResult)
	if err != nil {
		return nil, errors.New("解析报文异常")
	}
	if stockResult.Status != 0 {
		return nil, errors.New("get data error")
	}
	if len(stockResult.Hq) <= 0 {
		return nil, errors.New("not get any data")
	}

	stocks := make([]Stock, len(stockResult.Hq))
	for n, hq := range stockResult.Hq {
		s := Stock{
			Date:string(hq[0]),
			Start:string(hq[1]),
			End:string(hq[2]),
			Change:string(hq[3]),
			ChangePercent:string(hq[4]),
			Low:string(hq[5]),
			High:string(hq[6]),
			DealCount:string(hq[7]),
			DealAmount:string(hq[8]),
			Code:stockCode,
		}
		stocks[n] = s
	}
	return stocks, nil
}

func (stock *Stock)GetUrl(stockCode string, startDate string, endDate string) string {
	return "http://q.stock.sohu.com/hisHq?" +
		"code=" + "zs_" + stockCode +
		"&start=" + startDate + "&end=" + endDate +
		"&stat=1&order=D&period=d&callback=historySearchHandler&rt=jsonp&r=0.09105574639477387&0.021587371893673213"

}

func (stock Stock)GetAllStocks() []StockCode {
	shUrl := "http://quote.eastmoney.com/stocklist.html#sh"
	szUrl := "http://quote.eastmoney.com/stocklist.html#sz"
	docs, err := goquery.NewDocument(shUrl)
	if err != nil {
		log.Fatal("查询上海所有股票代码失败", err)
	}
	docs, err = goquery.NewDocument(szUrl)
	if err != nil {
		log.Fatal("查询深圳所有股票代码失败", err)
	}
	return nil

}


