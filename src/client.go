package main

import (
    "fmt"
    "database/sql"
    _ "mysql"
)

func printUser(ch chan int) {
    db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql")
    rows, err := db.Query("select user from user")

    if err != nil {

    }

    for rows.Next() {
        var user string
        if err := rows.Scan(&user); err != nil {

        }
        fmt.Printf("%s\n", user)
    }

    ch <- 1
}

func main() {
    ch := make(chan int, 4)

    go printUser(ch)

    fmt.Printf("hello world\n")

    <-ch
    <-ch
}
