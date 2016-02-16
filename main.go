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
    io.WriteString(w, "1.0")
}

func Login(w http.ResponseWriter, r *http.Request) { 


  	// api_url = http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/userinfo
  	// Instantiating the OAuth2 package to exchange the Code for a Token
  	conf := &oauth2.Config{
  	  ClientID:     "grafana",
  	  ClientSecret: "10b54f7c-a8ed-4a61-abd7-eb993d12367b",
  	  RedirectURL: "http://localhost:8090/login", 
  	  Scopes:       []string{"name","email"},
  	  Endpoint: oauth2.Endpoint{
  	    AuthURL:  "http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/auth", 
  	    TokenURL: "http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/token" } }

  	// Getting the Code that we got from Auth0
  	code := r.URL.Query().Get("code")
    fmt.Println("reached up to the code ")

  	if code == "" {
  	    http.Redirect(w,r,"http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/auth?access_type=online&client_id=grafana&redirect_uri=http%3A%2F%2Flocalhost%3A8090%2Flogin&response_type=code",301)
        return
  	}

  	// Exchanging the code for a token
  	token, err := conf.Exchange(oauth2.NoContext, code)
  	if err != nil {
        fmt.Println("this is the error")
  	  http.Error(w, err.Error(), http.StatusInternalServerError)
  	  return
  	}

	var data struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	client := conf.Client(oauth2.NoContext, token)

    apiUrl :=  "http://indiadevres3.cloudapp.net:8080/auth/realms/Waygum/protocol/openid-connect/userinfo"

	res, err := client.Get(apiUrl)
	if err != nil {
		fmt.Println("Error")
		return 
	}

	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil { 
		fmt.Println("Some Error")
		return 
	}

	fmt.Println("User Name " +  data.Name)
	fmt.Println("User Id " +  data.Id)
	fmt.Println("User Email" +  data.Email)


    io.WriteString(w, data.Email)

}


