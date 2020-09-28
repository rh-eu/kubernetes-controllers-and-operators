package helper

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// StdToStdMiddleware ...
// Middleware without "github.com/julienschmidt/httprouter"
func StdToStdMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do stuff
		next.ServeHTTP(w, r)
	})
}

// StdToJulienMiddleware ...
// Middleware for a standard handler returning a "github.com/julienschmidt/httprouter" Handle
// https://stackoverflow.com/questions/43964461/how-to-use-middlewares-when-using-julienschmidt-httprouter-in-golang/43964572
func StdToJulienMiddleware(next http.Handler) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// do stuff
		w.Write([]byte("Hello Julien!"))
		next.ServeHTTP(w, r)
	}
}

// JulienToJulienMiddleware ...
// Pure "github.com/julienschmidt/httprouter" middleware
func JulienToJulienMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// do stuff
		next(w, r, ps)
	}
}

// JulienHandler ...
func JulienHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// do stuff
	}
}

// StdHandler ...
func StdHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do stuff
	})
}

// MyTestToJulienMiddleware ....
func MyTestToJulienMiddleware(next http.Handler) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// do stuff
		log.Printf("MyTest to Julien ...%v", w)
		//w.Write([]byte("Admission Controller!"))
		next.ServeHTTP(w, r)
	}
}

// MyTestHandler ...
func MyTestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("Now I am here!"))
		log.Println("MyTestHandler !!!")
	})
}
