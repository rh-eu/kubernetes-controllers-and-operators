package handler

import (
	"log"
	"net/http"

	admission "k8s.io/api/admission/v1beta1"
)

// AdmitFunc is a type for building Kubernetes admission webhooks. An AdmitFunc
// should check whether an admission request is valid, and shall return an
// admission response that sets AdmissionResponse.Allowed to true or false as
// needed.
//
// Users wishing to build their own admission handlers should satisfy the
// AdmitFunc type, and pass it to an AdmissionHandler for serving over HTTP.
//
// Note: this mirrors the type in k8s source:
// https://github.com/kubernetes/kubernetes/blob/v1.13.0/test/images/webhook/main.go#L43-L44
type AdmitFunc func(reviewRequest *admission.AdmissionReview) (*admission.AdmissionResponse, error)

// AdmissionHandler ...
type AdmissionHandler struct {
	// The AdmitFunc to invoke for this handler.
	AdmitFunc AdmitFunc
}

func (ah *AdmissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("... testing here ...")
}

// MyTestHandler ...
//func MyTestHandler() http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("Now I am here!"))
//		log.Println("MyTestHandler !!!")
//	})
//}
