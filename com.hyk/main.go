package main

import (
	//"github.com/PuerkitoBio/goquery"
	"../com.hyk/mg"
	//"fmt"
	"../com.hyk/query"
	"log"
	"time"
	"../com.hyk/util"
)

type User struct {
	Name string
	Age  int
}

var cacthStockFlag = true

func cacthDapanData(stock *query.Stock) {
	for {
		log.Println("catch dapan data begin ....")
		code := "000001"//上证指数
		initDate := "20010101"//初始时间
		stock, err := getDpData(code)
		if err != nil {
			log.Println("查询大盘数据失败", err)
		}
		if stock != nil {
			d, err := time.Parse(util.Layout_2, stock.Date)
			if err != nil {
				log.Fatal("时间解析错误", err)
			} else {
				initDate = d.Format(util.Layout)
			}
		}
		t, err := time.Parse(util.Layout, initDate)
		if err != nil {
			log.Fatal("时间解析错误", err)
		}
		beginDate, err := time.Parse(util.Layout, initDate)
		endDate := t.Add(1 * 24 * time.Hour)
		for n := 1; endDate.Before(time.Now()); n++ {
			stocks, err := stock.GetData(beginDate.Format(util.Layout), endDate.Format(util.Layout), code)
			if err != nil {
				log.Println(err)
			}
			if stocks != nil {
				insertMonggo(stocks)
			}
			beginDate = endDate
			endDate = endDate.Add(1 * 24 * time.Hour)
		}
		log.Println("catch dapan data end ....")
		time.Sleep(time.Hour * 24)
	}
}

func getDpData(code string) (*query.Stock, error) {
	mg := getMongo()
	result := &query.Stock{}
	err := mg.FindSortLimit("admin", "stock", map[string]interface{}{"code":code}, "-date", 1, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getMongo() *mg.Mg {
	return &mg.Mg{"mongodb://127.0.0.1:27017", }
}

func insertMonggo(stocks []query.Stock) {
	mg := getMongo()
	inface := make([]interface{}, len(stocks))
	for n, s := range stocks {
		inface[n] = s
	}
	err := mg.Insert("admin", "stock", inface...)
	if err != nil {
		log.Println(err)
	}
}

func initStockCode(stock *query.Stock) {
	for {
		log.Println("catch StockCode data begin ....")
		cacthStockFlag = false
		stockCodes := stock.GetStockCodes()
		codeMap := make(map[string]string, len(stockCodes))
		inface := make([]interface{}, len(stockCodes))
		i := 0
		for _, s := range stockCodes {
			if codeMap[s.Code] == "" {
				codeMap[s.Code] = s.Name
				inface[i] = s
				i++
			}
		}

		mg := getMongo()
		err := mg.RemoveAll("admin", "stockCode")
		if err != nil {
			log.Println("删除数据失败", err)
		}
		err = mg.Insert("admin", "stockCode", inface[:i]...)
		if err != nil {
			log.Println("插入数据失败", err)
		}
		//time.Sleep(time.Second * 5)
		cacthStockFlag = true
		log.Println("catch StockCode data end ....")
		time.Sleep(time.Hour * 24 * 30)

	}
}

func cacthStock(stock *query.Stock) {
	log.Println(cacthStockFlag)
	for ; cacthStockFlag; {
		log.Println("catch cacthStock data begin ....")
		mg := getMongo()
		stockCodes := make([]query.StockCode, 1024)
		err := mg.FindAll("admin", "stockCode", nil, &stockCodes)
		if err != nil {
			log.Println("查询stockCodes失败")
		}
		for _, code := range stockCodes {
			stocks, err := stock.GetStocks(code.Code, code.Name)
			if err != nil {
				log.Println("查询 stock.GetStocks【" + code.Code + " " + code.Name + "】失败", err)
				continue
			}
			result := &query.Stock{}
			err = mg.FindSortLimit("admin", "stock", map[string]interface{}{"code":code.Code}, "-date", 1, result)
			if err != nil {
				log.Println(" mg.FindSortLimit【" + code.Code + " " + code.Name + "】未查到数据", err)
			}else {
				stocks = filterStock(stocks, result)
			}
			if len(stocks) == 0 {
				continue
			}
			inface := make([]interface{}, len(stocks))
			for i, stk := range stocks {
				inface[i] = stk
			}
			err = mg.Insert("admin", "stock", inface...)
			if err != nil {
				log.Println("插入【" + code.Code + " " + code.Name + "】失败", err)
			}

		}
		log.Println("catch cacthStock data end ....")
		time.Sleep(time.Hour * 12)

	}
}

func filterStock(stocks []query.Stock, st *query.Stock) []query.Stock {
	result := make([]query.Stock, len(stocks))
	stTime, _ := time.Parse(util.Layout_2, st.Date)
	i := 0
	for _, stock := range stocks {
		stockTime, _ := time.Parse(util.Layout_2, stock.Date)
		if stockTime.After(stTime) {
			result[i] = stock
			i ++
		}
	}
	return result[:i]
}

func main() {
	mg := getMongo()
	result := &query.Stock{}
	err := mg.FindSortLimit("admin", "stock", map[string]interface{}{"code": "1"}, "-date", 1, result)
	if err != nil {
		log.Println(err)
	}
	stock := query.Stock{}
	go cacthDapanData(&stock)
	go initStockCode(&stock)
	time.Sleep(time.Second * 10)
	go cacthStock(&stock)

	for {
		time.Sleep(time.Hour * 1000)
	}
	//http.HandleFunc("")
	//http.ListenAndServe()
	//for {
	//	//req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	//	//http.DefaultClient.Do(req)
	//	i := 0
	//	i++
	//	fmt.Println(strconv.Itoa(i))
	//	time.Sleep(time.Hour * 1000)
	//}
}



