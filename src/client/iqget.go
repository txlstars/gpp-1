package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gpplog"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

func httpProxy(*http.Request) (*url.URL, error) {
	// return nil, nil
	httpProxyUrl, err := url.Parse("http://web-proxy.tencent.com:8080")
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"httpProxy err": err}).Error("infoqCrawler")
		return nil, err
	}
	return httpProxyUrl, err
}

type InfoqThemeInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type InfoqDocSimpleInfo struct {
	Uuid			string `json:"uuid"`
	Article_title	string `json:"article_title"`
	Article_summary string `json:"article_summary"`
	Views			int    `json:"views"`
	Publish_time	int64  `json:"publish_time"`
	Love			int    `json:"love"`
}

type InfoqIndexList struct {
	Book_list      []InfoqDocSimpleInfo `json:"book_list"`
	Hot_day_list   []InfoqDocSimpleInfo `json:"hot_day_list"`
	Hot_month_list []InfoqDocSimpleInfo `json:"hot_month_list"`
	Hot_year_list  []InfoqDocSimpleInfo `json:"hot_year_list"`
	Recommend_list []InfoqDocSimpleInfo `json:"recommend_list"`
	Theme_list     []InfoqThemeInfo     `json:"theme_list"`
}

type InfoqIndex struct {
	Code int            `json:"code"`
	Data InfoqIndexList `json:"data"`
}

type InfoqTheme struct {
	Code int                  `json:"code"`
	Data []InfoqDocSimpleInfo `json:"data"`
}

// 首页精选内容|热点|快讯|专题
func infoqCrawlerIndexList() {
	req, err := http.NewRequest("GET", "https://www.infoq.cn/public/v1/article/getIndexList", nil)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"NewRequest": err}).Error("infoqCrawler")
		return
	}

	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Referer", "https://www.infoq.cn/")

	tr := &http.Transport{
		Proxy:              httpProxy,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	rsp, err := client.Do(req)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err": err}).Error("infoqCrawler")
		return
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err": err}).Error("infoqCrawler")
		return
	}

	if json.Valid(body) == false {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json invalid": err}).Error("infoqCrawlerIndex")
		return
	}

	var indexRsp InfoqIndex
	if err := json.Unmarshal(body, &indexRsp); err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json parse fail": err}).Error("infoqCrawlerIndex")
		return
	}

	if indexRsp.Code != 0 {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"server error": indexRsp.Code}).Error("infoqCrawlerIndex")
		return
	}

	for _, v := range indexRsp.Data.Book_list {
		infoqDocInsertDB(&v, "book_list")
	}

	for _, v := range indexRsp.Data.Hot_day_list {
		infoqDocInsertDB(&v, "hot_day")
	}

	for _, v := range indexRsp.Data.Hot_month_list {
		infoqDocInsertDB(&v, "hot_month")
	}

	for _, v := range indexRsp.Data.Hot_year_list {
		infoqDocInsertDB(&v, "hot_year")
	}

	for _, v := range indexRsp.Data.Recommend_list {
		infoqDocInsertDB(&v, "high_quality")
	}

	for _, v := range indexRsp.Data.Theme_list {
		infoqCrawlerThemeList(v.Id)
	}
}

// 主题列表数据
func infoqCrawlerThemeList(themeId int) {
	postData := `{"id":` + strconv.Itoa(themeId) + `}`
	req, err := http.NewRequest("POST", "https://www.infoq.cn/public/v1/theme/getArtList", strings.NewReader(postData))
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"NewRequest": err}).Error("infoqCrawler")
		return
	}
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Referer", "https://www.infoq.cn/")

	tr := &http.Transport{
		Proxy:              httpProxy,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	rsp, err := client.Do(req)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err": err}).Error("infoqCrawler")
		return
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err": err}).Error("infoqCrawler")
		return
	}

	// fmt.Println(string(body))

	if json.Valid(body) == false {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json invalid": err}).Error("infoqCrawlerIndex")
		return
	}

	var themeRsp InfoqTheme
	if err := json.Unmarshal(body, &themeRsp); err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json parse fail": err}).Error("infoqCrawlerIndex")
		return
	}

	if themeRsp.Code != 0 {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"server error": themeRsp.Code}).Error("infoqCrawlerIndex")
		return
	}

	for _, v := range themeRsp.Data {
		infoqDocInsertDB(&v, "theme")
	}
}

// 首页推荐和垂类tab推荐列表数据
func infoqCrawlerGuidRecomList(postData string, dataType string) {
	req, err := http.NewRequest("POST", "https://www.infoq.cn/public/v1/article/getList", strings.NewReader(postData))
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Referer", "https://www.infoq.cn/topic/architecture")

	tr := &http.Transport{
		Proxy:              httpProxy,
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
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err": err}).Error("infoqCrawler")
		return
	}

	if json.Valid(body) == false {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json invalid": err}).Error("infoqCrawlerIndex")
		return
	}

	var themeRsp InfoqTheme
	if err := json.Unmarshal(body, &themeRsp); err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json parse fail": err}).Error("infoqCrawlerIndex")
		return
	}

	if themeRsp.Code != 0 {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"server error": themeRsp.Code}).Error("infoqCrawlerIndex")
		return
	}

	for _, v := range themeRsp.Data {
		infoqDocInsertDB(&v, dataType)
	}
}

type InfoqDocDetailInfo struct {
	Uuid           string               `json:"uuid"`
	Article_title  string               `json:"article_title"`
	Views          int                  `json:"views"`
	Love           int                  `json:"love"`
	Publish_time   int64                `json:"publish_time"`
	Recommend_list []InfoqDocSimpleInfo `json:"recommend_list"`
}

type InfoqDoc struct {
	Code int                `json:"code"`
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
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"infoqCrawlerDoc": err}).Error("infoqCrawler")
		return
	}

	tr := &http.Transport{
		Proxy:              httpProxy,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	rsp, err := client.Do(req)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"infoqCrawlerDoc Do": err}).Error("infoqCrawler")
		return
	}
	defer rsp.Body.Close()

	body, _ := ioutil.ReadAll(rsp.Body)
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"err": err}).Error("infoqCrawler")
		return
	}

	if json.Valid(body) == false {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json invalid": err}).Error("infoqCrawlerIndex")
		return
	}

	var docRsp InfoqDoc
	if err := json.Unmarshal(body, &docRsp); err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"json parse fail": err}).Error("infoqCrawlerIndex")
		return
	}

	if docRsp.Code != 0 {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"server error": docRsp.Code}).Error("infoqCrawlerIndex")
		return
	}

	// 详情数据写入DB
	// infoqDocInsertDB(&docRsp.Data, "test")

	for _, v := range docRsp.Data.Recommend_list {
		fmt.Printf("docid:%s, title:%s, Views:%d, publishtime:%v\n", v.Uuid, v.Article_title, v.Views, v.Publish_time)
	}
}

func infoqCrawlerStart(parentWaitGroup *sync.WaitGroup) {
	defer parentWaitGroup.Done()

	for ;; {
		// 首页运营数据
		// fmt.Println("-------------------------index op------------------------")
		infoqCrawlerIndexList()

		// 首页推荐列表数据
		// fmt.Println("-------------------------index recom------------------------")
		infoqCrawlerGuidRecomList(`{"size":30}`, "index recom")

		// 架构tab
		// fmt.Println("-------------------------architecture------------------------")
		infoqCrawlerGuidRecomList(`{"type":1,"size":30,"id":8}`, "architecture")

		// 云计算
		// fmt.Println("-------------------------cloud computing ------------------------")
		infoqCrawlerGuidRecomList(`{"type":1,"size":12,"id":11`, "cloud-computing")

		// 前端
		// fmt.Println("-------------------------front end------------------------")
		infoqCrawlerGuidRecomList(`{"type":1,"size":30,"id":33}`, "front-end")

		// 运维
		// fmt.Println("-------------------------operation------------------------")
		infoqCrawlerGuidRecomList(`{"type":1,"size":12,"id":38}`, "operation")

		// 文章相关阅读
		// fmt.Println("-------------------------doc releate------------------------")
		// infoqCrawlerDocAndReleate(`T3yPFdi88*GKZwHR2bPT`)

		time.Sleep(12 * time.Hour)
	}
}

func infoqDocInsertDB(docInfo *InfoqDocSimpleInfo, docType string) {
	docStaticInfo := &DocStaticInfo{
		docid:       docInfo.Uuid,
		src:         "infoq",
		title:       docInfo.Article_title,
		summary:	 docInfo.Article_summary,
		typex:       docType,
		pageUrl:     "https://www.infoq.cn/article/" + docInfo.Uuid,
		publishTime: uint32(docInfo.Publish_time / 1000),
	}
	docStaticInfoTask <- docStaticInfo
}
