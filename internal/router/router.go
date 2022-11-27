package router

import (
  "github.com/gorilla/mux"
  "net/http"
  "sha256service/internal/tools"
  "sha256service/internal/handler"
)

func GetRoutesHandler(h *handler.Handler) http.Handler {
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
