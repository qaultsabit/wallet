package route

import (
	"net/http"

	"github.com/qaultsabit/wallet/src/interface/rest/handlers"

	"github.com/go-chi/chi/v5"
)

func HealthRouter(h handlers.IHealthHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", h.Ping)

	return r
}
