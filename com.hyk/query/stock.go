package query

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	//"container/list"
	"log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"../util"
	"time"
	"strconv"
	"github.com/PuerkitoBio/goquery"
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
	Start         string //开盘
	End           string //收盘
	Low           string //最低
	High          string //最高
	Change        string //涨幅
	ChangePercent string //涨幅百分比
	DealCount     string //成交量
	DealAmount    string //成交额
}

type StockCode struct {
	Code string
	Name string
}

func (stock *Stock) GetData(startDate string, endDate string, stockCode string) ([]Stock, error) {

	req, err := http.NewRequest("GET", stock.GetUrl(stockCode, startDate, endDate), nil)
	fmt.Println("url :" + stock.GetUrl(stockCode, startDate, endDate))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)")
	fmt.Println("do")
	res, err := http.DefaultClient.Do(req)
	fmt.Println("f")
	if err != nil {
		fmt.Println("http request error")
		return nil, err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("result : " + string(result))
	if len(result) < 100 {
		return nil, nil
	}
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
			Name:"上证指数",
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
		log.Println("查询股票代码失败 :(" + url + ")", err)
	}
	size := len(docs.Find("#quotesearch ul li").Nodes)
	stockCodes := make([]StockCode, size)
	docs.Find("#quotesearch ul li").Each(
		func(i int, contentSelection *goquery.Selection) {
			dst := make([]byte, 2 * len(contentSelection.Text()))
			decoder.Transform(dst, []byte(contentSelection.Text()), true)
			codeStr := util.ByteString(dst)
			arry := strings.Split(codeStr, "(")
			stockCodes[i] = StockCode{
				Name : arry[0],
				Code : arry[1][:len(arry[1]) - 1],
			}
		})

	return stockCodes
}

func (stock *Stock) GetStocks(code string, name string) ([]Stock, error) {
	url := "http://www.aigaogao.com/tools/history.html?s=" + code
	fmt.Println(url)
	//decoder := simplifiedchinese.GBK.NewDecoder()
	docs, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	selections := docs.Find("#ctl16_contentdiv table tr")

	if selections == nil || len(selections.Nodes) == 0 {
		return nil, errors.New("not find data " + code)
	}

	arrys := make([]string, len(selections.Nodes))
	var index = 0
	selections.Each(func(i int, selection *goquery.Selection) {
		tds := selection.Find(".altertd")
		if tds != nil && len(tds.Nodes) > 0 {
			var str string
			tds.Each(func(i int, selection *goquery.Selection) {
				str += (selection.Text() + ";")

			})
			if len(strings.TrimSpace(str)) != 0 {
				arrys[index] = str
				index++
			}
		}
	})
	stocks := make([]Stock, index)
	for i:=0; i < index; i++ {
		//fmt.Println(i)
		if i == 733 {
			fmt.Println(arrys[i] + "****" + strconv.Itoa(len(strings.TrimSpace(arrys[i]))))
		}
		tds := strings.Split(arrys[i], ";")
		t, _ := time.Parse(util.Layout, tds[0])
		stocks[i] = Stock{
			Code : code,
			Name : name,
			Date : t.Format(util.Layout_2),
			Start: tds[1],
			High : tds[2],
			Low  : tds[3],
			End  : tds[4],
			DealCount : tds[5],
			DealAmount: tds[6],
			Change : tds[7],
			ChangePercent : tds[8],
		}
	}

	return stocks, nil
}





