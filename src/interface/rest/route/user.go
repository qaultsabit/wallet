package route

import (
	"net/http"

	"github.com/qaultsabit/wallet/src/interface/rest/handlers"

	"github.com/go-chi/chi/v5"
)

func UserRouter(h handlers.UserHandlerInterface) http.Handler {
	r := chi.NewRouter()

	r.Post("/create_user", h.Register)
	r.Post("/login", h.Login)

	return r
}
