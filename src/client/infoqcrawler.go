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
)

func httpProxy(*http.Request) (*url.URL, error) {
	httpProxyUrl, err := url.Parse("http://web-proxy.tencent.com:8080")
	if err != nil {
		gpplog.GetLogger("infoq").WithFields(logrus.Fields{"httpProxy err" : err}).Error("infoqCrawler")
		return nil, err
	}

	return httpProxyUrl, err
}

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

	body, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(body))
}

func infoqCrawlerStart() {
	infoqCrawlerGuidRecomList("{\"size\":20}")
}
