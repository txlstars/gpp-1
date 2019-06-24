package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gpplog"
	"sync"
)

type DocStaticInfo struct {
	docid       string
	src         string
	title       string
	summary     string
	typex       string
	pageUrl     string
	publishTime uint32
}

type DocDynamicInfo struct {
	docid   string
	src     string
	viewNum int
	likeNme int
}

func addDocStaticInfoToDB(docStaticInfo *DocStaticInfo) {
	tx, err := db.Begin()
	if err != nil {
		gpplog.GetLogger("mysql").WithFields(logrus.Fields{
			"err":   err,
			"docid": docStaticInfo.docid,
			"src":   docStaticInfo.src,
		}).Error("addDocStaticInfoToDB")
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			gpplog.GetLogger("mysql").WithFields(logrus.Fields{
				"err":   err,
				"docid": docStaticInfo.docid,
				"src":   docStaticInfo.src,
			}).Error("addDocStaticInfoToDB")
		}
	} ()

	// 文章静态信息
	sqlContent := `insert ignore into t_doc_info_test(docid, title, summary, src, type, publish_time, pageurl) 
	values (?, ?, ?, ?, ?, from_unixtime(?), ?)`
	if _, err = tx.Exec(sqlContent,
		docStaticInfo.docid,
		docStaticInfo.title,
		docStaticInfo.summary,
		docStaticInfo.src,
		docStaticInfo.typex,
		docStaticInfo.publishTime,
		docStaticInfo.pageUrl); err != nil {
		return
	}

	// 文章动态信息
	sqlContent = `insert ignore into t_doc_dynamic_info_test(docid, src, view_num, like_num) values (?, ?, ?, ?)`
	if _, err = tx.Exec(sqlContent, docStaticInfo.docid, docStaticInfo.src, 0, 0); err != nil {
		return
	}

	err = tx.Commit()
}

func handleDocStaticInfo(parentWaitGroup *sync.WaitGroup) {
	defer parentWaitGroup.Done()

	for {
		docStaticInfo, ok := <-docStaticInfoTask
		if ok {
			addDocStaticInfoToDB(docStaticInfo)
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

// 数据更新任务

func storageStart(parentWaitGroup *sync.WaitGroup) {
	defer parentWaitGroup.Done()

	// gpp数据库连接代理
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/gpp")
	if err != nil {
		gpplog.GetLogger("client").WithFields(logrus.Fields{"err": err}).Error("mysql client open fail")
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
