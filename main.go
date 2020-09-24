package main

import (
	"log"
	"net/http"
	"os"

	admissioncontrol "github.com/elithrar/admission-control"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}

func main() {

	// Set up logging
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "loc", kitlog.DefaultCaller)

	r := mux.NewRouter()

	admissions := r.PathPrefix("/admission-control").Subrouter()

	r.HandleFunc("/", myHandler)

	admissions.Handle("/enforce-pod-annotations", &admissioncontrol.AdmissionHandler{
		AdmitFunc: admissioncontrol.EnforcePodAnnotations(
			[]string{"kube-system"},
			map[string]func(string) bool{
				"k8s.questionable.services/hostname": func(string) bool { return true },
			}),
		Logger: logger,
	}).Methods(http.MethodPost)

	log.Println("Server is listening on port 8443.")
	log.Println("...TEST")
	//log.Fatal(http.ListenAndServeTLS(":8443", "certs/mifommcrt.pem", "certs/mifommkey.pem", r))
	log.Fatal(http.ListenAndServeTLS(":8443", "certs/mifommcrt.pem", "certs/mifommkey.pem", admissioncontrol.LoggingMiddleware(logger)(r)))
}
