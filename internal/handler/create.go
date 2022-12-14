package handler

import (
  "net/http"
  "sha256service/internal/tools"
  "fmt"
  "sha256service/internal/request"
)

func (h *Handler) HandleCreateHash(w http.ResponseWriter, r *http.Request) {
  req := &request.CreateHashRequest{}
  var err error
  if err = tools.ReadRequest(r, req); err != nil {
    tools.WriteRequestError(w, r, err)
    return
  }
  if err = validateRequest(req); err != nil {
    tools.WriteRequestError(w, r, err)
    return
  }
  rawBytes, err := request.PreparePayload(req.Payload, req.PayloadType)
  if err != nil {
    tools.WriteInternalError(w, r, err)
    return
  }
  resp, err := h.GetHash(rawBytes, req.Secret)
  if err != nil {
    tools.WriteInternalError(w, r, fmt.Errorf("cannot get item hash: %v", err))
    return
  }
  if err := tools.WriteResponse(w, r, resp); err != nil {
    tools.WriteInternalError(w, r, err)
    return
  }
}
