package app

import (
	"net/http"

	"vibe-golang-template/pkg/response"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				response.WriteErrorFrom(w, r, response.InternalError("error.internal"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func chainMiddlewares(handler http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	wrapped := handler
	for i := len(mws) - 1; i >= 0; i-- {
		wrapped = mws[i](wrapped)
	}
	return wrapped
}
