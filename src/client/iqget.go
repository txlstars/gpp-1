package main

import (
	"net/http"
	"gpplog"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"fmt"
	// "time"
	"net/url"
	"strings"
	"encoding/json"
	"strconv"
)

var docIdSet map[string]bool

// var releateDocIdSet []string

func httpProxy(*http.Request) (*url.URL, error) {
	httpProxyUrl, err := url.Parse("http://web-proxy.tencent.com:8080")
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"httpProxy err" : err}).Error("infoqCrawler")
		return nil, err
	}
	return httpProxyUrl, err
}

type InfoqThemeInfo struct {
	Id		int		`json:"id"`
	Name	string	`json:"name"`
}

type InfoqDocSimpleInfo struct {
	Uuid			string	`json:"uuid"`
	Article_title	string	`json:"article_title"`
	Views			int		`json:"views"`
	Publish_time	int64	`json:"publish_time"`
	Love			int		`json:"love"`
}

type InfoqIndexList struct {
	Book_list		[]InfoqDocSimpleInfo	`json:"book_list"`
	Hot_day_list	[]InfoqDocSimpleInfo	`json:"hot_day_list"`
	Hot_month_list	[]InfoqDocSimpleInfo	`json:"hot_month_list"`
	Hot_year_list	[]InfoqDocSimpleInfo	`json:"hot_year_list"`
	Recommend_list  []InfoqDocSimpleInfo	`json:"recommend_list"`
	Theme_list		[]InfoqThemeInfo		`json:"theme_list"`
}

type InfoqIndex struct {
	Code int			`json:"code"`
	Data InfoqIndexList `json:"data"`
}

type InfoqTheme struct {
	Code int					`json:"code"`
	Data []InfoqDocSimpleInfo	`json:"data"`
}

// 首页精选内容|热点|快讯|专题
func infoqCrawlerIndexList() {
	req, err := http.NewRequest("GET", "https://www.infoq.cn/public/v1/article/getIndexList", nil)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"NewRequest" : err}).Error("infoqCrawler")
		return
	}

	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Referer", "https://www.infoq.cn/")

	tr := &http.Transport{
		Proxy: httpProxy,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	rsp, err := client.Do(req)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err" : err}).Error("infoqCrawler")
		return
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err" : err}).Error("infoqCrawler")
		return
	}

	if json.Valid(body) == false {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json invalid" : err}).Error("infoqCrawlerIndex")
		return
	}

	var indexRsp InfoqIndex
	if err := json.Unmarshal(body, &indexRsp); err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json parse fail" : err}).Error("infoqCrawlerIndex")
		return
	}

	if indexRsp.Code != 0 {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"server error" : indexRsp.Code}).Error("infoqCrawlerIndex")
		return
	}

	fmt.Println("-------------------------------------book_list-------------------------------------")
	for _, v := range indexRsp.Data.Book_list {
		if _, ok := docIdSet[v.Uuid]; ok == false {
			docIdSet[v.Uuid] = true
			// infoqCrawlerDocAndReleate(v.Uuid)
		}

		fmt.Printf("docid:%s, title:%s, Views:%d\n", v.Uuid, v.Article_title, v.Views)
	}

	fmt.Println("-----------------------------------hot_day_list------------------------------------")
	for _, v := range indexRsp.Data.Hot_day_list {
		fmt.Printf("docid:%s, title:%s, Views:%d\n", v.Uuid, v.Article_title, v.Views)
	}

	fmt.Println("-----------------------------------hot_month_list-----------------------------------")
	for _, v := range indexRsp.Data.Hot_month_list {
		fmt.Printf("docid:%s, title:%s, Views:%d\n", v.Uuid, v.Article_title, v.Views)
	}

	fmt.Println("-----------------------------------hot_year_list------------------------------------")
	for _, v := range indexRsp.Data.Hot_year_list {
		fmt.Printf("docid:%s, title:%s, Views:%d\n", v.Uuid, v.Article_title, v.Views)
	}

	fmt.Println("-----------------------------------recommend_list------------------------------------")
	for _, v := range indexRsp.Data.Recommend_list {
		fmt.Printf("docid:%s, title:%s, Views:%d\n", v.Uuid, v.Article_title, v.Views)
	}

	for _, v := range indexRsp.Data.Theme_list {
		fmt.Println("-------------------------------------theme_list--------------------------------------")
		infoqCrawlerThemeList(v.Id) 
	}
}

// 主题列表数据
func infoqCrawlerThemeList(themeId int) {
	postData := `{"id":` + strconv.Itoa(themeId) + `}`
	req, err := http.NewRequest("POST", "https://www.infoq.cn/public/v1/theme/getArtList", strings.NewReader(postData))
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"NewRequest" : err}).Error("infoqCrawler")
		return
	}
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Referer", "https://www.infoq.cn/")

	tr := &http.Transport{
		Proxy: httpProxy,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	rsp, err := client.Do(req)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err" : err}).Error("infoqCrawler")
		return
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err" : err}).Error("infoqCrawler")
		return
	}

	// fmt.Println(string(body))

	if json.Valid(body) == false {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json invalid" : err}).Error("infoqCrawlerIndex")
		return
	}

	var themeRsp InfoqTheme
	if err := json.Unmarshal(body, &themeRsp); err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json parse fail" : err}).Error("infoqCrawlerIndex")
		return
	}

	if themeRsp.Code != 0 {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"server error" : themeRsp.Code}).Error("infoqCrawlerIndex")
		return
	}

	for _, v := range themeRsp.Data {
		fmt.Printf("docid:%s, title:%s, Views:%d\n", v.Uuid, v.Article_title, v.Views)
	}
}

// 首页推荐和垂类tab推荐列表数据
func infoqCrawlerGuidRecomList(postData string) {
	req, err := http.NewRequest("POST", "https://www.infoq.cn/public/v1/article/getList", strings.NewReader(postData))
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Referer", "https://www.infoq.cn/topic/architecture")

	tr := &http.Transport{
		Proxy: httpProxy,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	rsp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err" : err}).Error("infoqCrawler")
		return
	}

	if json.Valid(body) == false {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json invalid" : err}).Error("infoqCrawlerIndex")
		return
	}

	var themeRsp InfoqTheme
	if err := json.Unmarshal(body, &themeRsp); err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json parse fail" : err}).Error("infoqCrawlerIndex")
		return
	}

	if themeRsp.Code != 0 {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"server error" : themeRsp.Code}).Error("infoqCrawlerIndex")
		return
	}

	for _, v := range themeRsp.Data {
		fmt.Printf("docid:%s, title:%s, Views:%d\n", v.Uuid, v.Article_title, v.Views)
	}
}

type InfoqDocDetailInfo struct {
	Uuid			string					`json:"uuid"`
	Article_title	string					`json:"article_title"`
	Views			int						`json:"views"`
	Love			int						`json:"love"`
	Publish_time	int64					`json:"publish_time"`
	Recommend_list  []InfoqDocSimpleInfo	`json:"recommend_list"`
}

type InfoqDoc struct {
	Code int				`json:"code"`
	Data InfoqDocDetailInfo `json:"data"`
}

// 文章详情页数据|相关阅读文章列表
func infoqCrawlerDocAndReleate(docId string) {
	postData := `{"uuid":"` + docId + `"}`

	req, err := http.NewRequest("POST", "https://www.infoq.cn/public/v1/article/getDetail", strings.NewReader(postData))
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Referer", "https://www.infoq.cn")
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"infoqCrawlerDoc" : err}).Error("infoqCrawler")
		return
	}

	tr := &http.Transport{
		Proxy: httpProxy,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	rsp, err := client.Do(req)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"infoqCrawlerDoc Do" : err}).Error("infoqCrawler")
		return
	}
	defer rsp.Body.Close()

	body, _ := ioutil.ReadAll(rsp.Body)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err" : err}).Error("infoqCrawler")
		return
	}

	if json.Valid(body) == false {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json invalid" : err}).Error("infoqCrawlerIndex")
		return
	}

	var docRsp InfoqDoc
	if err := json.Unmarshal(body, &docRsp); err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json parse fail" : err}).Error("infoqCrawlerIndex")
		return
	}

	if docRsp.Code != 0 {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"server error" : docRsp.Code}).Error("infoqCrawlerIndex")
		return
	}

	// 详情数据写入DB
	infoqDocInsertDB(&docRsp.Data, "test")

	for _, v := range docRsp.Data.Recommend_list {
		fmt.Printf("docid:%s, title:%s, Views:%d, publishtime:%v\n", v.Uuid, v.Article_title, v.Views, v.Publish_time)
	}
}


func infoqCrawlerStart() {
	docIdSet = make(map[string]bool)
	// releateDocIdSet = make([]string)

	// 首页运营数据
	fmt.Println("-------------------------index op------------------------")
	infoqCrawlerIndexList()

	// 首页推荐列表数据
	fmt.Println("-------------------------index recom------------------------")
	infoqCrawlerGuidRecomList(`{"size":12}`)

	// 架构tab
	fmt.Println("-------------------------architecture------------------------")
	infoqCrawlerGuidRecomList(`{"type":1,"size":12,"id":8}`)

	// 云计算
	fmt.Println("-------------------------cloud computing ------------------------")
	infoqCrawlerGuidRecomList(`{"type":1,"size":12,"id":11`)

	// 前端
	fmt.Println("-------------------------front end------------------------")
	infoqCrawlerGuidRecomList(`{"type":1,"size":12,"id":33}`)

	// 运维
	fmt.Println("-------------------------operation------------------------")
	infoqCrawlerGuidRecomList(`{"type":1,"size":12,"id":38}`)

	// 文章相关阅读
	fmt.Println("-------------------------doc releate------------------------")
	infoqCrawlerDocAndReleate(`T3yPFdi88*GKZwHR2bPT`)
}

func infoqDocInsertDB(docInfo *InfoqDocDetailInfo, docType string) error {
	sqlContent := `insert into t_doc_info_test(docid, title, src, type, views, loves, publish_time) values("`
	sqlContent += docInfo.Uuid + `", "`
	sqlContent += docInfo.Article_title + `", "`
	sqlContent += `infoq", "`
	sqlContent += docType + `", `
	sqlContent += strconv.Itoa(docInfo.Views) + `, `
	sqlContent += strconv.Itoa(docInfo.Love) + `, `
	sqlContent += strconv.FormatInt(docInfo.Publish_time / 1000, 10) + `)`

	fmt.Println(sqlContent)

	return nil
}

/*
rows, err := db.Exec(

if err != nil {
	gpplog.GetLogger("mysql_client").WithFields(log.Fields{"err" : err}).Error("mysql client fail")
	return
}
*/
