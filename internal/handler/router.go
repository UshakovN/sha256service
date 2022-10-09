package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"sha256service/internal/tools"
	"sha256service/pkg/sha256"
	"time"
)

func GetRoutesHandler(handler *Handler) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/create-hash-sum", handler.HandleCreateHashSum)
	router.HandleFunc("/get-hash-sum", handler.HandleGetHashSum)
	router.HandleFunc("/health", tools.HandleHealthRequest)
	return router
}

func validateCreateHashSumRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandleCreateHashSum(w http.ResponseWriter, r *http.Request) {
	validateCreateHashSumRequest(w, r)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		tools.WriteRequestError(w, r, fmt.Errorf("cannot read request data"))
		return
	}
	startTime := time.Now()
	hash := sha256.New()
	hashSum := hash.Sum(data)
	hashingTime := time.Since(startTime)
	contentType := http.DetectContentType(data)
	err = r.Body.Close()
	if err != nil {
		tools.WriteInternalError(w, r, fmt.Errorf("cannot close request body"))
	}
	resp := &ItemHash{
		MimeType:    contentType,
		HashSum:     fmt.Sprintf("%x", hashSum),
		HashedAt:    time.Now().UTC().Format(time.UnixDate),
		HashingTime: hashingTime.String(),
	}
	if err := h.PutItemHashInStore(resp); err != nil {
		tools.WriteInternalError(w, r, err)
		return
	}
	tools.WriteResponse(w, r, resp)
}

func validateGetHashSumRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	}
	sum := r.Header.Get("hash_sum")
	if sum == "" {
		tools.WriteRequestError(w, r, fmt.Errorf("parameter 'hash_sum' is mandatory"))
	}
}

func (h *Handler) HandleGetHashSum(w http.ResponseWriter, r *http.Request) {
	validateGetHashSumRequest(w, r)
	sum := r.Header.Get("hash_sum")
	item, err := h.GetItemHashFromStore(sum)
	if err != nil {
		tools.WriteInternalError(w, r, fmt.Errorf(""))
	}
	tools.WriteResponse(w, r, item)
}
