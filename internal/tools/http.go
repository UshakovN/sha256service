package tools

import (
  "encoding/json"
  "log"
  "net/http"
  "fmt"
  "net/url"
  "io"
)

func ReadRequest(r *http.Request, message any) error {
  var payload []byte
  var err error
  switch r.Method {

  case http.MethodPost:
    defer r.Body.Close()
    payload, err = io.ReadAll(r.Body)
    if err != nil {
      return fmt.Errorf("cannot read request payload: %v", err)
    }

  case http.MethodGet:
    m, err := convertQueryToMap(r.URL.Query())
    if err != nil {
      return fmt.Errorf("cannot convert query to map: %v", err)
    }
    payload, err = json.Marshal(m)
    if err != nil {
      return fmt.Errorf("cannot convert request payload to json: %v", err)
    }

  case http.MethodOptions:
    return nil

  default:
    return fmt.Errorf("unsupported method")
  }

  if err := json.Unmarshal(payload, message); err != nil {
    return fmt.Errorf("cannot parse request payload to message")
  }
  return nil
}

func convertQueryToMap(query url.Values) (map[string]interface{}, error) {
  m := make(map[string]interface{}, len(query))
  for key, values := range query {
    for idx, value := range values {
      if idx > 0 {
        return nil, fmt.Errorf("repeated query parameters unsopported")
      }
      m[key] = value
    }
  }
  return m, nil
}

func WriteResponse(w http.ResponseWriter, _ *http.Request, p any) error {
  b, err := json.MarshalIndent(p, "", "  ")
  if err != nil {
    return fmt.Errorf("cannot convert data to json: %v", err)
  }
  w.Header().Add("Content-Type", "application/json")
  _, err = w.Write(b)
  if err != nil {
    return fmt.Errorf("cannot write to message writer: %v", err)
  }
  return nil
}

func WriteRequestError(w http.ResponseWriter, r *http.Request, err error) {
  log.Printf("error in request %s: %v", r.RequestURI, err)
  http.Error(w, err.Error(), http.StatusBadRequest)
}

func WriteInternalError(w http.ResponseWriter, r *http.Request, err error) {
  log.Printf("internal error in request %s: %v", r.RequestURI, err)
  http.Error(w, err.Error(), http.StatusInternalServerError)
}

func HandleHealthRequest(w http.ResponseWriter, _ *http.Request) {
  w.WriteHeader(http.StatusOK)
  _, err := w.Write([]byte("/ok"))
  if err != nil {
    log.Printf("cannot write to response writer: %v", err)
  }
}
