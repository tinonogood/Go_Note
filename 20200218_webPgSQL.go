package main

import (
    "net/http"
    "encoding/json"
    "strconv"
    "github.com/gorilla/mux"
    "html/template"
    "github.com/lib/pg"
    "database/sql"
)

const (
    host	= "127.0.0.1"
    port	= 5432
    user	= "user"
    password	= "password"
    dbname	= "db"
)


type User struct {
    Id  int
    Name        string
    Age int
}

func getUsersHTML(w http.ResponseWriter, r *http.Request){
    tpl, _ := template.ParseFiles("users.html")
    tpl.Execute(w, allUsers)
}

func getUser(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    for _, user := range allUsers {
        if user.Id == id {
            js, err :=json.Marshal(user)
            if err != nil{
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Header().Set("Content-Type","application/json")
            w.Write(js)
        }
    }
}

func getNewUserID() int {
    newID := 0
    for _, user := range allUsers {
        if (newID <= user.Id) {
            newID = user.Id + 1
        }
    }
    return newID
}

func createUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    var user User
    _ = json.NewDecoder(r.Body).Decode(&user)
    user.Id = getNewUserID()
    allUsers = append(allUsers, user)
    js, err := json.Marshal(user)
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(js)
}

func updateUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    for index, originalUser := range allUsers{
        if originalUser.Id == id{
            allUsers = append(allUsers[:index], allUsers[index+1:]...)
            var updateUser User
            _ = json.NewDecoder(r.Body).Decode(&updateUser)
            updateUser.Id = id
            allUsers = append(allUsers, updateUser)
            js, err := json.Marshal(updateUser)
            if err != nil{
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Write(js)
          }
    }
}

func deleteUser(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
    vars := mux.Vars(r)
    idURL := vars["id"]
    id, _ := strconv.Atoi(idURL)
    for index, originalUser := range allUsers{
        if originalUser.Id == id{
            allUsers = append(allUsers[:index], allUsers[index+1:]...)
            break
        }
    }
}


func main() {
    r := mux.NewRouter()

    r.HandleFunc("/users", getUsersHTML).Methods("GET")
    r.HandleFunc("/users/{id}", getUser).Methods("GET")
    r.HandleFunc("/users/create", createUser).Methods("POST")
    r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

    http.ListenAndServe(":8080", r)
}

/* users.html

<!DOCTYPE html>
<html>
    <head>
        <meta charset='utf-8'>
        <title>Users Infomation</title>
    </head>
    <body>
        {{range .}}
            <h1>ID: {{.Id}}</h1>
            <h2>Name: {{.Name}}</h2>
            <h2>Age: {{.Age}}</h2>
        {{end}}
    </body>
</html>

 */
