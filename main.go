package main

import (
    "fmt"
    "io"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/codegangsta/negroni"
    oauth2 "github.com/goincremental/negroni-oauth2"
    sessions "github.com/goincremental/negroni-sessions"
    "github.com/goincremental/negroni-sessions/cookiestore"
)

var cookies map[string]string

func KeyCloak(config *oauth2.Config) negroni.Handler { 
    authUrl     := "http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/auth"
    tokenUrl    := "http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/token"
    return oauth2.NewOAuth2Provider(config, authUrl, tokenUrl)
}

func main() {
    fmt.Println("Starting sample auth...")

    n := negroni.Classic()

    n.Use(sessions.Sessions("mysession", cookiestore.New([]byte("secret12"))))

    n.Use(KeyCloak(&oauth2.Config{
            ClientID : "grafana",
            ClientSecret : "10b54f7c-a8ed-4a61-abd7-eb993d12367b",
            RedirectURL : "http://127.0.0.1:8090/oauth2callback",
            Scopes : []string{"name","email"} }))

    router := mux.NewRouter()

    router.HandleFunc("/", Home) 
    router.HandleFunc("/version", Version) 

    router.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(w,"World !") 
    })

    n.Use(oauth2.LoginRequired())

    n.UseHandler(router)
    n.Run(":8090")
}

func Home(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Welcome to Sample Oauth")
}

func Version(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "1.0")
}

