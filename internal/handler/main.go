package handler

import (
  "net/http"
  "fmt"
  "sha256service/internal/tools"
  "html/template"
)

const (
  htmlTemplatePrefixPath = "./templates/"
  htmlTemplateMain       = "main.html"
)

func (h *Handler) HandleMainPage(w http.ResponseWriter, r *http.Request) {
  path := fmt.Sprint(htmlTemplatePrefixPath, htmlTemplateMain)
  t, err := template.ParseFiles(path)
  if err != nil {
    tools.WriteInternalError(w, r, fmt.Errorf("cannot parse %s", path))
    return
  }
  if err := t.Execute(w, nil); err != nil {
    tools.WriteInternalError(w, r, fmt.Errorf("cannot execute %s", path))
    return
  }
}
