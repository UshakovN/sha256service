package handler

import (
  "net/http"
  "fmt"
  "sha256service/internal/tools"
  "html/template"
)

const htmlTemplateAbout = "about.html"

func (h *Handler) HandleAboutPage(w http.ResponseWriter, r *http.Request) {
  path := fmt.Sprint(htmlTemplatePrefixPath, htmlTemplateAbout)
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
