package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"gregvader/triple-tac-toe/database/database_type"
	tokens "gregvader/triple-tac-toe/security/helpers/token"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func DBTransactionMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txHandle := db.Begin()

			defer func() {
				if r := recover(); r != nil {
					txHandle.Rollback()
				}
			}()
			params := mux.Vars(r)
			cntx := context.WithValue(context.Background(), database_type.Transactional, txHandle)
			newR := r.WithContext(cntx)
			newR = mux.SetURLVars(newR, params)
			newR.Header.Set("Authorization", r.Header.Get("Authorization"))
			next.ServeHTTP(w, newR)

			txHandle.Commit() //If not explicitly rollback-ed in handlers, just commit.
		})
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenValidErr := tokens.TokenValid(r)

		if tokenValidErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(tokenValidErr.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
