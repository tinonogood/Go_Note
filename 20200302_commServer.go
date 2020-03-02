package main

import (
    "net/http"
    "encoding/json"
    "strconv"
    "github.com/gorilla/mux"
    "html/template"
    _ "github.com/lib/pq"
    "database/sql"
    "fmt"
)


const (
    host        = "127.0.0.1"
    port        = 5432
    user        = "user"
    password    = "password"
    dbname      = "dbname"
)

var Db *sql.DB

type User struct {
    Id  int
    Name        string
    Age int
}

func init() {
    var err error
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    Db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
}


// ===================
// Serialize
// ===================

func handleConnection(conn net.Conn) {
    dec := gob.NewDecoder(conn)
    p := &User{}

    dec.Decode(p)
    _, err = Db.Exec("UPDATE users SET name = $2, age = $3 WHERE id = $1", p.Id, p.Name, p.Age)
    if err != nil{
        panic(err)
    }
    log.Println("Hello ",p.Name,", Your Age is ",p.Age);

    conn.Close()
}




func main() {




    // Serialize
    // 
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            continue
        }
        go handleConnection(conn)
    }
}
