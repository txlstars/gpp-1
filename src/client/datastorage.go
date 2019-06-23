package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gpplog"
	"sync"
)

type DocStaticInfo struct {
	docid       string
	src         string
	title       string
	typex       string
	pageUrl     string
	publishTime uint32
}

type DocDynamicInfo struct {
	docid      string
	src        string
	viewNum    int
	likeNme    int
	commentNum int
	updateTime uint32
}

func handleDocStaticInfo(parentWaitGroup *sync.WaitGroup) {
	defer parentWaitGroup.Done()

	for {
		docStaticInfo, ok := <-docStaticInfoTask
		if ok {
			fmt.Println(docStaticInfo)
			sqlContent := `insert into t_doc_info_test(docid, title, src, type, publish_time, pageurl) 
            values (?, ?, ?, ?, from_unixtime(?), ?)`
			_, err := db.Exec(sqlContent,
				docStaticInfo.docid,
				docStaticInfo.title,
				docStaticInfo.src,
				docStaticInfo.typex,
				docStaticInfo.publishTime,
				docStaticInfo.pageUrl)
			fmt.Println(err)
		} else {
			break
		}
	}
}

func handleDocDynamicInfo(parentWaitGroup *sync.WaitGroup) {
	defer parentWaitGroup.Done()

	for {
		docDynamicInfo, ok := <-docDynamicInfoTask
		if ok {
			fmt.Println(docDynamicInfo)
		} else {
			break
		}
	}
}

var db *sql.DB
var docStaticInfoTask = make(chan *DocStaticInfo, 100)
var docDynamicInfoTask = make(chan *DocDynamicInfo, 100)

func storageStart(parentWaitGroup *sync.WaitGroup) {
	defer parentWaitGroup.Done()

	// gpp数据库连接代理
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/gpp")
	if err != nil {
		gpplog.GetLogger("client").WithFields(log.Fields{"err": err}).Error("mysql client open fail")
		return
	}
	defer db.Close()

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(2 * 2)
	for i := 0; i < 2; i++ {
		go handleDocStaticInfo(waitGroup)
		go handleDocDynamicInfo(waitGroup)
	}
	waitGroup.Wait()
}
