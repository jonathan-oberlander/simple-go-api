package admin

import (
	"net/http"
	"os"
)

// Portal '
type Portal struct {
	password string
}

// Handler '
func (a *Portal) Handler(w http.ResponseWriter, r *http.Request) {
	usr, pwd, ok := r.BasicAuth()
	if !ok || pwd != a.password || usr != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("<html><h1>Admin Protal</h1></html>"))
}

// NewAdminPortal '
func NewAdminPortal() Portal {
	pwd := os.Getenv("ADMIN_PASSWORD")
	if pwd == "" {
		panic("required env var ADMIN_PASSWORD")
	}
	return Portal{pwd}
}
