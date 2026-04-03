package main

import (
	"context"
	"fmt"
	"net/http"

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
	if r, ok := ctx.Value(ctxKey{}).(uax.Result); ok {
		return r
	}
	return uax.Result{}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ua := GetUA(r.Context())
	fmt.Fprintf(w, "Hello, %s user on %s!\n", ua.ShortBrowser(), ua.ShortOS())
	if ua.IsBot {
		fmt.Fprintf(w, "(Detected bot: %s)\n", ua.Bot.Name)
	}
}

func main() {
	parser, _ := uax.NewParser()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	wrapped := UAMiddleware(parser)(mux)

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", wrapped)
}
