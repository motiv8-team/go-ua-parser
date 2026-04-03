package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	uax "github.com/motiv8-team/go-ua-parser"
)

type ctxKey struct{}

func UAMiddleware(parser *uax.Parser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := parser.ParseString(r.Header.Get("User-Agent"))
			ctx := context.WithValue(r.Context(), ctxKey{}, result)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUA(ctx context.Context) uax.Result {
	if v, ok := ctx.Value(ctxKey{}).(uax.Result); ok {
		return v
	}
	return uax.Result{}
}

func main() {
	parser, _ := uax.NewParser()
	r := chi.NewRouter()
	r.Use(UAMiddleware(parser))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ua := GetUA(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"browser": ua.ShortBrowser(),
			"os":      ua.ShortOS(),
			"device":  ua.DeviceClass(),
			"bot":     ua.IsBot,
		})
	})

	http.ListenAndServe(":8080", r)
}
