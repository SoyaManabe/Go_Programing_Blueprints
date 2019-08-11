package main

import (
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
  "flag"
  //"fmt"
  //"os"
  //"app/trace"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/gomniauth/providers/facebook"
  "github.com/stretchr/gomniauth/providers/github"
  "github.com/stretchr/gomniauth/providers/google"
  "github.com/stretchr/objx"
)

//temp1 shows 1 template
type templateHandler struct {
  once     sync.Once
  filename string
  templ    *template.Template
  }

// ServeHTTP
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func() {
    t.templ = 
      template.Must(template.ParseFiles(filepath.Join("templates",
        t.filename)))
  })
  data := map[string]interface{}{
    "Host": r.Host,
  }
  if authCookie, err := r.Cookie("auth"); err == nil {
    data["UserData"] = objx.MustFromBase64(authCookie.Value)
    //fmt.Println(data["UserData"])
  }

  t.templ.Execute(w, data)
}
func main() {
  var addr = flag.String("addr", ":8080", "ApplicationAddress")
  flag.Parse() //Analyse flag
  // Gomniauth set up
  gomniauth.SetSecurityKey("soya1374")
  gomniauth.WithProviders(
    facebook.New("", "", "http://localhost:8080/auth/callback/facebook"),
    github.New("", "", "http://localhost:8080/auth/callback/github"),
    google.New("683114142464-odunei9fljnjeaitjfv7f75uhpd0ui9h.apps.googleusercontent.com", "0EyHWaf1gGnH2tUsuSL5tHf3", "http://localhost:8080/auth/callback/google"),
  )
  r := newRoom()
  // r.tracer = trace.New(os.Stdout)
  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/login", &templateHandler{filename: "login.html"})
  http.Handle("/room", r)
  http.HandleFunc("/auth/", loginHandler)
  http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
      Name: "auth",
      Value: "",
      Path: "/",
      MaxAge: -1,
    })
    w.Header()["Location"] = []string{"/chat"}
    w.WriteHeader(http.StatusTemporaryRedirect)
  })
  // start chat room
  go r.run()
  // start webServer
  log.Println("Starting Webserver... PORT: ", *addr)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe;", err)
  }
}
