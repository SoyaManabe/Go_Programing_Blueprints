package main

import (
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
  "flag"
  "os"
  "app/trace"
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
  t.templ.Execute(w, r)
}
func main() {
  var addr = flag.String("addr", ":8080", "ApplicationAddress")
  flag.Parse() //Analyse flag
  r := newRoom()
  r.tracer = trace.New(os.Stdout)
  http.Handle("/", &templateHandler{filename: "chat.html"})
  http.Handle("/room", r)
  // start chat room
  go r.run()
  // start webServer
  log.Println("Starting Webserver... PORT: ", *addr)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe;", err)
  }
}

