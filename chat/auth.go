package main
import (
  "net/http"
  "strings"
  "log"
  "fmt"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/objx"
)
type authHandler struct {
  next http.Handler
}
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
    // Not authorized
    w.Header().Set("Location", "/login")
    w.WriteHeader(http.StatusTemporaryRedirect)
  } else if err != nil {
    // Some errors
    panic(err.Error())
  } else {
    // Success
    h.next.ServeHTTP(w, r)
  }
}
func MustAuth(handler http.Handler) http.Handler {
  return &authHandler{next: handler}
}

// loginHandler wait login action with third party application
// /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
  segs := strings.Split(r.URL.Path, "/")
  action := segs[2]
  provider := segs[3]
  switch action {
  case "login":
    provider, err := gomniauth.Provider(provider)
    if err != nil {
      log.Fatalln("Failed to get auth provider:", provider, "-", err)
    }
    loginUrl, err := provider.GetBeginAuthURL(nil, nil)
    if err != nil{
      log.Fatalln("Error occured during GetBeginAuthURL calling:", provider, "-", err)
    }
    w.Header().Set("Location", loginUrl)
    w.WriteHeader(http.StatusTemporaryRedirect)
  case "callback":
    provider, err := gomniauth.Provider(provider)
    if err != nil {
      log.Fatalln("Failed to get auth provider:", provider, "-", err)
    }

    creds, err :=
        provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
    if err != nil {
      log.Fatalln("Failue to finish authentication", provider, "-", err)
    }

    user, err := provider.GetUser(creds)
    if err != nil {
      log.Fatalln("Failue to get User", provider, "-", err)
    }
    authCookieValue := objx.New(map[string]interface{}{
      "name": user.Name(),
      "avatar_url": user.AvatarURL(),
    }).MustBase64()
    http.SetCookie(w, &http.Cookie{
      Name: "auth",
      Value: authCookieValue,
      Path: "/"})
    w.Header()["Location"] = []string{"/chat"}
    w.WriteHeader(http.StatusTemporaryRedirect)

  default:
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Action %s not supported", action)
  }
}
