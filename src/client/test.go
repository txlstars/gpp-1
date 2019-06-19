package main

import (
	"database/sql"
	"fmt"
	"gpplog"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/go-redis/redis"
)

func mysqlClientTest(ch chan int) {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql")
	if err != nil {
		gpplog.GetLogger("mysql_client").WithFields(log.Fields{"err" : err}).Error("mysql client fail")
		return
	}

	defer db.Close()

	rows, err := db.Query("select User, Host, plugin from user limit 1")
	if err != nil {
		gpplog.GetLogger("mysql_client").WithFields(log.Fields{"err" : err}).Error("mysql client fail")
		return
	}

	for rows.Next() {
		var user, host, plugin string
		if err := rows.Scan(&user, &host, &plugin); err != nil {
			gpplog.GetLogger("mysql_client").WithFields(log.Fields{"err" : err}).Error("mysql client fail")
			continue
		}
		gpplog.GetLogger("mysql_client").WithFields(log.Fields{
			"user" : user,
			"host" : host,
			"plugin" : plugin,
		}).Info("mysql client succ")
	}

	ch <- 1
}

func redisClientTest(ch chan int) {
	client := redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password:	"",
		DB:			0,
	})

	key := "foo"
	val, err := client.Get(key).Result()
	if err == nil {
		gpplog.GetLogger("redis-client").WithFields(log.Fields{"key" : key,
		"val" : val}).Info("redis client get:")
	}

	ch <- 1
}

func main() {
	// ch := make(chan int, 4)

	// go mysqlClientTest(ch)

	// go redisClientTest(ch)

	infoqCrawlerStart()

	fmt.Printf("hello world\n")

	/*
	for i :=0; i < 2; i++ {
		<-ch
	}
	*/
}
