package handler

import (
	"log"

	admission "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newDefaultDenyResponse returns an AdmissionResponse with a Result sub-object,
// and defaults to allowed = false.
func newDefaultDenyResponse() *admission.AdmissionResponse {
	return &admission.AdmissionResponse{
		Allowed: false,
		Result:  &metav1.Status{},
	}
}

// EnforcePodAnnotations ...
func EnforcePodAnnotations(ignoredNamespaces []string, requiredAnnotations map[string]func(string) bool) AdmitFunc {
	return func(admissionReview *admission.AdmissionReview) (*admission.AdmissionResponse, error) {

		log.Println("... inside AdmitFunc ...")
		resp := newDefaultDenyResponse()
		return resp, nil
	}
}
