package main

import (
//    "fmt"
    "net/http"
    "encoding/json"
    "strconv"
    "github.com/gorilla/mux"
)

type User struct {
    Id	int
    Name	string
    Age	int
}

var allUsers = []User {
    {Id: 1,Name: "john", Age: 20},
    {Id: 2,Name: "mary", Age: 18},
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idURL := vars["id"]
        id,_ := strconv.Atoi(idURL)
        for _, user := range allUsers {
            if user.Id == id {
                js, err := json.Marshal(user)
                if err != nil{
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }

                w.Header().Set("Content-Type","application/json")
                w.Write(js)
            }
        }
    })

    http.ListenAndServe(":80", r)
}

