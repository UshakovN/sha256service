package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"sha256service/internal/tools"
	"strings"
)

const (
	contentTypePlainText = "text/plain"
)

func GetRoutesHandler(handler *Handler) http.Handler {
	rootRouter := mux.NewRouter()
	apiRouter := rootRouter.PathPrefix("/").Subrouter()
	apiRouter.HandleFunc("/get-hash", handler.HandleGetHash).Methods(http.MethodGet)
	apiRouter.HandleFunc("/create-hash", handler.HandleCreateHash).Methods(http.MethodPost)
	apiRouter.HandleFunc("/create-http-content-hash", handler.HandleCreateHttpContentHash).Methods(http.MethodPost)
	rootRouter.HandleFunc("/health", tools.HandleHealthRequest).Methods(http.MethodGet)
	return rootRouter
}

func (h *Handler) HandleCreateHash(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		tools.WriteRequestError(w, r, fmt.Errorf("cannot read request data"))
		return
	}
	itemHash, err := h.GetItemHash(data)
	if err != nil {
		tools.WriteInternalError(w, r, fmt.Errorf("cannot get item hash: %v", err))
	}
	err = r.Body.Close()
	if err != nil {
		tools.WriteInternalError(w, r, fmt.Errorf("cannot close request body"))
	}
	resp := itemHash
	if err := h.PutItemHashInStore(resp); err != nil {
		tools.WriteInternalError(w, r, err)
		return
	}
	tools.WriteResponse(w, r, resp)
}

func (h *Handler) getHashQueryParam(r *http.Request) (string, error) {
	sum := r.URL.Query().Get("sum")
	if sum == "" {
		return "", fmt.Errorf("parameter 'sum' is mandatory")
	}
	return sum, nil
}

func (h *Handler) HandleGetHash(w http.ResponseWriter, r *http.Request) {
	sum, err := h.getHashQueryParam(r)
	if err != nil {
		tools.WriteRequestError(w, r, err)
		return
	}
	item, found, err := h.GetItemHashFromStore(sum)
	if found {
		item.HashFound = true
	}
	if err != nil {
		tools.WriteInternalError(w, r, err)
		return
	}
	tools.WriteResponse(w, r, item)
}

func (h *Handler) getContentUrl(r *http.Request) (string, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read request body: %v", err)
	}
	if err := r.Body.Close(); err != nil {
		return "", fmt.Errorf("cannot close request body: %v", err)
	}
	contentType := http.DetectContentType(data)
	if !strings.HasPrefix(contentType, contentTypePlainText) {
		return "", fmt.Errorf("invalid request body content type")
	}
	url := string(data)
	if !tools.MatchUrl(url) {
		return "", fmt.Errorf("uncorrect content url")
	}
	return tools.StripWebPrefix(url), nil
}

func (h *Handler) HandleCreateHttpContentHash(w http.ResponseWriter, r *http.Request) {
	contentUrl, err := h.getContentUrl(r)
	if err != nil {
		tools.WriteRequestError(w, r, fmt.Errorf("cannot get content url: %v", err))
		return
	}
	data, err := h.httpClient.Get(contentUrl, r.Header)
	if err != nil {
		tools.WriteInternalError(w, r, fmt.Errorf("cannot load content from url: %v", err))
		return
	}
	itemHash, err := h.GetItemHash(data)
	if err != nil {
		tools.WriteInternalError(w, r, fmt.Errorf("cannot get item hash: %v", err))
	}
	resp := itemHash
	if err := h.PutItemHashInStore(resp); err != nil {
		tools.WriteInternalError(w, r, err)
		return
	}
	tools.WriteResponse(w, r, resp)
}
