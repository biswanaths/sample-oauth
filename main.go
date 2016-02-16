package main

import (
    "fmt"
    "io"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func main() {
    fmt.Println("Starting sample auth...")
    r := mux.NewRouter()
    r.HandleFunc("/version", Version)
    r.HandleFunc("/login", Login)
    http.Handle("/", r)
    http.ListenAndServe("127.0.0.1:8090", nil)
    fmt.Println("Listening now")
}

func Version(w http.ResponseWriter, r *http.Request) {
    for _, cookie := range r.Cookies() {
        fmt.Println(cookie.Name)
    }
    io.WriteString(w, "1.0")
}

func Login(w http.ResponseWriter, r *http.Request) { 

  	conf := &oauth2.Config{
  	    ClientID:     "grafana",
  	    ClientSecret: "10b54f7c-a8ed-4a61-abd7-eb993d12367b",
  	    RedirectURL: "http://localhost:8090/login", 
  	    Scopes:       []string{"name","email"},
  	    Endpoint: oauth2.Endpoint{
  	      AuthURL:  "http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/auth", 
  	      TokenURL: "http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/token" } }

    apiUrl :=  "http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/userinfo"

  	code := r.URL.Query().Get("code")

  	if code == "" {
        url := conf.AuthCodeURL("")
        http.Redirect(w, r, url, 301)
        return
  	}

  	token, err := conf.Exchange(oauth2.NoContext, code)
  	if err != nil {
  	  http.Error(w, err.Error(), http.StatusInternalServerError)
  	  return
  	}

	var data struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	client := conf.Client(oauth2.NoContext, token)

	res, err := client.Get(apiUrl)
	if err != nil {
  	    http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil { 
  	    http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

    io.WriteString(w, "Login successful for " + data.Email)
}


