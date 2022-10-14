package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"sha256service/internal/tools"
	"time"
)

const (
	authHeaderName = "x-session-id"
)

func (h *Handler) readAuthRequest(r *http.Request) (*AuthRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read request body: %v", err)
	}
	req := &AuthRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		return nil, fmt.Errorf("cannot parse request body: %v", err)
	}
	return req, nil
}

func (h *Handler) validateAuthRequest(r *AuthRequest) error {
	if r.Login == "" || r.Password == "" {
		return fmt.Errorf("login or password was empty")
	}
	return nil
}

func (h *Handler) createNewUser(login string, password string) *User {
	return &User{
		Login:    login,
		Password: fmt.Sprintf("%x", h.hashClient.Sum([]byte(password))),
	}
}

func (h *Handler) createNewUserSession(user *User) *UserSession {
	return &UserSession{
		Id:            uuid.NewString(),
		UserLogin:     user.Login,
		Expired:       false,
		LastLoginDate: time.Now().UTC().Format(time.UnixDate),
	}
}

func (h *Handler) HandlerSignUp(w http.ResponseWriter, r *http.Request) {
	req, err := h.readAuthRequest(r)
	if err != nil {
		tools.WriteRequestError(w, r, fmt.Errorf("cannot read sign in request: %v", err))
		return
	}
	if err := h.validateAuthRequest(req); err != nil {
		tools.WriteRequestError(w, r, fmt.Errorf("invalid sign in request: %v", err))
		return
	}
	user := h.createNewUser(req.Login, req.Password)
	session := h.createNewUserSession(user)
	g, _ := errgroup.WithContext(h.ctx)
	g.Go(func() error {
		return h.PutUserToStore(user)
	})
	g.Go(func() error {
		return h.PutUserSessionToStore(session)
	})
	if err := g.Wait(); err != nil {
		tools.WriteInternalError(w, r, fmt.Errorf("cannot put new user data to storage: %v", err))
	}
	tools.WriteResponse(w, r, &AuthResponse{
		SessionId: session.Id,
	})
}

func (h *Handler) compareHashAndPassword(hash, password string) bool {
	return hash == fmt.Sprintf("%x", h.hashClient.Sum([]byte(password)))
}

func (h *Handler) HandleLogIn(w http.ResponseWriter, r *http.Request) {
	req, err := h.readAuthRequest(r)
	if err != nil {
		tools.WriteRequestError(w, r, fmt.Errorf("cannot read log in request: %v", err))
		return
	}
	if err := h.validateAuthRequest(req); err != nil {
		tools.WriteRequestError(w, r, fmt.Errorf("invalid log in request: %v", err))
		return
	}
	user, found, err := h.GetUserFromStore(req.Login)
	if err != nil {
		tools.WriteInternalError(w, r, fmt.Errorf("cannot get user from storage: %v", err))
		return
	}
	if !found {
		tools.WriteRequestError(w, r, fmt.Errorf("user '%s' not found in storage", req.Login))
		return
	}
	if !h.compareHashAndPassword(user.Password, req.Password) {
		tools.WriteRequestError(w, r, fmt.Errorf("wrond user password for user '%s'", req.Login))
		return
	}

	//todo getting session by user login
	tools.WriteResponse(w, r, &AuthResponse{
		SessionId: "",
	})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(authHeaderName)
		if authToken == "" {
			http.Error(w, fmt.Sprintf("auth token must be passed as '%s' header", authHeaderName), http.StatusForbidden)
			return
		}
		_, found, err := h.GetUserSessionFromStore(authToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot get user session: %v", err), http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, "user session not found", http.StatusForbidden)
			return
		}
		// todo check expired session

		//ctx := context.WithValue(r.Context(), ApiClientContextField, apiClient)
		//ctx = context.WithValue(ctx, UserContextField, response.User)
		//next.ServeHTTP(w, r.WithContext(ctx))
	})
}
