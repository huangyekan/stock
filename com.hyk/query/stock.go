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
	"golang.org/x/text/encoding/simplifiedchinese"
	"../util"
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
	Code          string
	Name          string
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


func (stock *Stock) GetStockCodes() []StockCode {
	url := "http://quote.eastmoney.com/stocklist.html#sh"
	decoder := simplifiedchinese.GBK.NewDecoder()
	docs, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("查询股票代码失败 :(" + url +")", err)
	}
	size := len(docs.Find("#quotesearch ul li").Nodes)
	stockCodes := make([]StockCode, size)
	docs.Find("#quotesearch ul li").Each(
		func(i int, contentSelection *goquery.Selection) {
			dst := make([]byte, 2*len(contentSelection.Text()))
			decoder.Transform(dst, []byte(contentSelection.Text()), true)
			codeStr := util.ByteString(dst)
			arry := strings.Split(codeStr, "(")
			stockCodes[i] = StockCode{
				Name : arry[0],
				Code : arry[1][:len(arry[1])-1],
			}
		})

	return stockCodes
}





