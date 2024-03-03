package main

import (
  "net/http"
  "os"
  // "log"
  "fmt"
  "github.com/a-h/templ"
  T "johtotimes.com/templates"
  // Types "johtotimes.com/internal/types"
  Globals "johtotimes.com/internal/globals"
)

func generateHeader() templ.Component {
  fmt.Println("indexHandler")

  // tmpl.ExecuteTemplate(w, "base", data)
  // return T.Head("Johto Times", data)
  return T.Base("Johto Times")
  
}

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = "3000"
  }
  fmt.Println("Listening on port " + port)

  assets := http.FileServer(http.Dir(Globals.AssetPath))

  mux := http.NewServeMux()
  // mux.HandleFunc("/", indexHandler)
  // // mux.HandleFunc("/", hello)
  mux.Handle("/", templ.Handler(generateHeader()))
	prefix := "/"+Globals.AssetPath+"/"
  mux.Handle(prefix, http.StripPrefix(prefix, assets))
  http.ListenAndServe(":"+port, mux)
}
