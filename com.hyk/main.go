package main

import (
	//"github.com/PuerkitoBio/goquery"
	"../com.hyk/mg"
	//"fmt"
	"../com.hyk/query"
	"log"
	"time"
	"fmt"
	"strconv"
)

type User struct {
	Name string
	Age  int
}

var initDate = "20010101"//初始时间
var codes = []string{"000001", "002059", }
var layout = "20060102" //时间格式

func cacthData(stock *query.Stock, code string) {
	t, err := time.Parse(layout, initDate)
	if err != nil {
		log.Fatal("时间解析错误", err)
	}
	beginDate, err := time.Parse(layout, initDate)
	endDate := t.Add(30 * 24 * time.Hour)
	for n:=1; endDate.Before(time.Now()); n++ {
		stocks, err := stock.GetData(beginDate.Format(layout), endDate.Format(layout), code)
		if err != nil {
			log.Fatal(err)
		}
		insertMonggo(stocks)
		beginDate = endDate
		endDate = endDate.Add(30 * 24 * time.Hour)
		fmt.Println("inser " + strconv.Itoa(n) + " data")

	}
}

func insertMonggo(stocks []query.Stock) {
	mg := &mg.Mg{
		"mongodb://127.0.0.1:27017",
	}
	fmt.Println(stocks)
	inface := make([]interface{}, len(stocks))
	for n, s := range stocks {
		inface[n] = s
	}
	err := mg.Insert("admin", "stock", inface...)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	stock := query.Stock{}
	cacthData(&stock, "000001")
	//http.HandleFunc("")
	//http.ListenAndServe()
}



