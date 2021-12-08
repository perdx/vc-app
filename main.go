package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/perdx/vc-app/handler"
)

// func init() {
// 	os.Setenv("TENANT", "9125264c-86cb-45fe-baa2-e022db0590d6")
// 	// issuer & verifier authority (could be different)
// 	os.Setenv("AUTHORITY", "did:ion:EiBuwQM4Yu-r3NV2qQsaeu2ziZ03D4TUTKRCAZjDVVteIg:eyJkZWx0YSI6eyJwYXRjaGVzIjpbeyJhY3Rpb24iOiJyZXBsYWNlIiwiZG9jdW1lbnQiOnsicHVibGljS2V5cyI6W3siaWQiOiJzaWdfM2ZlZjk4ZDQiLCJwdWJsaWNLZXlKd2siOnsiY3J2Ijoic2VjcDI1NmsxIiwia3R5IjoiRUMiLCJ4IjoiSjZDeEE5U2QzeUV4Z2hTTDJ6OUx0YzYzMXZxbEJfRFV6bEo4QlU3WWZORSIsInkiOiJwWmhYRG9LbVNNc2FlcnY4N3V3ME5zOWZLZkFEN1hLZmZQMjJreXBPVXZNIn0sInB1cnBvc2VzIjpbImF1dGhlbnRpY2F0aW9uIiwiYXNzZXJ0aW9uTWV0aG9kIl0sInR5cGUiOiJFY2RzYVNlY3AyNTZrMVZlcmlmaWNhdGlvbktleTIwMTkifV0sInNlcnZpY2VzIjpbeyJpZCI6ImxpbmtlZGRvbWFpbnMiLCJzZXJ2aWNlRW5kcG9pbnQiOnsib3JpZ2lucyI6WyJodHRwczovL3BlcmR4LmlvLyJdfSwidHlwZSI6IkxpbmtlZERvbWFpbnMifV19fV0sInVwZGF0ZUNvbW1pdG1lbnQiOiJFaUJJeTFTTXhzWjRwLS1uNkI1MVRXQjFDZWl4bDZqa3V6QVNUZmc2QWF4bFdBIn0sInN1ZmZpeERhdGEiOnsiZGVsdGFIYXNoIjoiRWlEcjlUbVlhSTEwamEtMFpzWkJ5ODJuWmcyT2F1MGpSbUFWSHE3Z09renBRQSIsInJlY292ZXJ5Q29tbWl0bWVudCI6IkVpQlJldTY5RXN0UVRLeDVCV1VfUzhyNDRoXzZpY20tVlEyWXZTNU54NkwtV3cifX0")
// 	// App registration credentials
// 	os.Setenv("CLIENT_ID", "275c62f4-75b0-40aa-839c-6fd2760e8e7d")
// 	os.Setenv("CLIENT_SECRET", "xSk7Q~awIY4CBl.E5FjS081tUVs8ckXHjQo6R")
// }
func main() {
	// s := store.Connect()
	// email.Send()

	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.StripSlashes,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "QUERY"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	i := handler.NewIssuer()

	r.Route("/issuer", func(r chi.Router) {
		r.Get("/issuance", i.Issue)
		r.Post("/callback", i.Callback)
		r.Get("/status/{state}", i.Status)
	})

	p := handler.NewVerifier()

	r.Route("/verifier", func(r chi.Router) {
		r.Get("/presentation", p.Present)
		r.Post("/callback", p.Callback)
		r.Get("/status/{state}", p.Status)
	})

	u := handler.NewUserRegistration()

	r.Route("/register", func(r chi.Router) {
		r.Post("/", u.Create)
		r.Get("/status", u.Status)
	})

	// start server
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Print(err)
	}
}
