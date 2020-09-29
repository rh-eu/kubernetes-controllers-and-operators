package handler

import (
	"fmt"
	"log"

	"golang.org/x/xerrors"
	admission "k8s.io/api/admission/v1beta1"

	//admission "k8s.io/api/admission/v1"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	podDeniedError       = "the submitted Pods are missing required annotations:"
	unsupportedKindError = "the submitted Kind is not supported by this admission handler:"
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

		kind := admissionReview.Request.Kind.Kind
		log.Printf("Kind: %+v", kind)

		deserializer := serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()

		// We handle all built-in Kinds that include a PodTemplateSpec, as described here:
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#pod-v1-core
		var namespace string
		annotations := make(map[string]string)
		// Extract the necessary metadata from our known Kinds

		log.Printf("Kind: %+v", kind)

		switch kind {
		case "Pod":
			pod := core.Pod{}
			log.Printf("Pod: %+v", pod)
			if _, _, err := deserializer.Decode(admissionReview.Request.Object.Raw, nil, &pod); err != nil {
				return nil, err
			}

			namespace = pod.GetNamespace()
			annotations = pod.GetAnnotations()
			log.Printf("Namespace: %s", namespace)
			log.Printf("Annotations: %v", annotations)
		case "Deployment":
			deployment := apps.Deployment{}
			if _, _, err := deserializer.Decode(admissionReview.Request.Object.Raw, nil, &deployment); err != nil {
				return nil, err
			}

			deployment.GetNamespace()
			annotations = deployment.Spec.Template.GetAnnotations()
		case "StatefulSet":
			statefulset := apps.StatefulSet{}
			if _, _, err := deserializer.Decode(admissionReview.Request.Object.Raw, nil, &statefulset); err != nil {
				return nil, err
			}

			namespace = statefulset.GetNamespace()
			annotations = statefulset.Spec.Template.GetAnnotations()
		case "DaemonSet":
			daemonset := apps.DaemonSet{}
			if _, _, err := deserializer.Decode(admissionReview.Request.Object.Raw, nil, &daemonset); err != nil {
				return nil, err
			}

			namespace = daemonset.GetNamespace()
			annotations = daemonset.Spec.Template.GetAnnotations()
		case "Job":
			job := batch.Job{}
			if _, _, err := deserializer.Decode(admissionReview.Request.Object.Raw, nil, &job); err != nil {
				return nil, err
			}

			namespace = job.Spec.Template.GetNamespace()
			annotations = job.Spec.Template.GetAnnotations()
		default:
			// TODO(matt): except for whitelisted namespaces
			log.Printf("Kind is not %v", kind)
			return nil, xerrors.Errorf("the submitted Kind is not supported by this admission handler: %s", kind)
		}

		// Ignore objects in whitelisted namespaces.
		for _, ns := range ignoredNamespaces {
			if namespace == ns {
				resp.Allowed = true
				resp.Result.Message = fmt.Sprintf("allowing admission: %s namespace is whitelisted", namespace)
				return resp, nil
			}
		}

		missing := make(map[string]string)
		log.Printf("Missing: %v", missing)
		// We check whether the (strictly matched) annotation key exists, and then run
		// our user-provided matchFunc against it. If we're missing any keys, or the
		// value for a key does not match, admission is rejected.
		for requiredKey, matchFunc := range requiredAnnotations {
			if matchFunc == nil {
				return resp, xerrors.Errorf("cannot validate annotations (%s) with a nil matchFunc", requiredKey)
			}

			if existingVal, ok := annotations[requiredKey]; !ok {
				// Key does not exist; add it to the missing annotations list
				missing[requiredKey] = "key was not found"
			} else {
				if matched := matchFunc(existingVal); !matched {
					missing[requiredKey] = "value did not match"
				}
				// Key exists & matchFunc returned OK.
			}
		}

		if len(missing) > 0 {
			return resp, xerrors.Errorf("%s %v", podDeniedError, missing)
		}

		// No missing or invalid annotations; allow admission
		resp.Allowed = true

		log.Printf("Response: %+v", resp)
		return resp, nil
	}
}

// ensureHasAnnotations checks whether the provided ObjectMeta has the required
// annotations. It returns both a map of missing annotations, and a boolean
// value if the meta had all of the provided annotations.
//
// The required annotations are case-sensitive; an empty string for the map
// value will match on key (only) and thus allow any value.
func ensureHasAnnotations(required map[string]string, annotations map[string]string) (map[string]string, bool) {
	missing := make(map[string]string)
	for requiredKey, requiredVal := range required {
		if existingVal, ok := annotations[requiredKey]; !ok {
			// Missing a required annotation; add it to the list
			missing[requiredKey] = requiredVal
		} else {
			// The key exists; does the value match?
			if existingVal != requiredVal {
				missing[requiredKey] = requiredVal
			}
		}
	}

	// If we have any missing annotations, report them to the caller so the user
	// can take action.
	if len(missing) > 0 {
		return missing, false
	}

	return nil, true
}
