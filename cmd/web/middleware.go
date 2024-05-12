package main

import (
	"fmt"
	"net/http"
)

// middleware functions act as a "chain" of handlers that pass the request
// to the next handler in the chain

// secureHeaders is a middleware that adds security headers to all responses
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

// logRequest is a middleware that logs the details of each request
// this middleware is on the application struct, giving it access
// to all the application logic
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// recoverPanic is a middleware that recovers from any panics and
// response with a 500 internal server error
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defered functions are executed in last-in, first-out order
		// and will always be run in the event of a panic
		// as GO unwind the call stack
		defer func() {
			// use recover() to check if there has been a panic
			// if so, respond with internal server error
			if err := recover(); err != nil {
				// close the connection
				w.Header().Set("Connection", "close")
				// whatever type the err is, here we format it into an error
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
