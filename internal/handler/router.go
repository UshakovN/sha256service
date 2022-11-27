package handler

import (
  "github.com/gorilla/mux"
  "net/http"
  "sha256service/internal/tools"
)

const (
  htmlTemplateMain       = "main.html"
  htmlTemplateAbout      = "about.html"
  htmlTemplatePrefixPath = "./templates/"
)

func GetRoutesHandler(h *Handler) http.Handler {
  r := mux.NewRouter()
  r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("assets/css/"))))
  r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("assets/js/"))))
  r.HandleFunc("/", h.HandleMainPage)
  r.HandleFunc("/about", h.HandleAboutPage)
  r.HandleFunc("/create", h.HandleCreateHash)
  r.HandleFunc("/compare", h.HandleCompareHash)
  r.HandleFunc("/health", tools.HandleHealthRequest)
  return r
}
