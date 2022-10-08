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

func GetRouterHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/create-hash-sum", HandleHashSum)
	router.HandleFunc("/get-hash-sum", nil)
	router.HandleFunc("/health", tools.HandleHealthRequest)
	return router
}

func validateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	}
}

func HandleHashSum(w http.ResponseWriter, r *http.Request) {
	validateRequest(w, r)
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
	resp := HashSumResponse{
		MimeType:    contentType,
		HashSum:     fmt.Sprintf("%x", hashSum),
		HashedAt:    time.Now().UTC().Format(time.UnixDate),
		HashingTime: hashingTime.String(),
	}
	tools.WriteResponse(w, r, resp)
}
