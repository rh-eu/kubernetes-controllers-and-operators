package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
)

func main() {
	server := getServer()
	http.HandleFunc("/validate", myHandler)
	log.Println("Server is listening on port 8443.")
	server.ListenAndServeTLS("certs/mifommcrt.pem", "certs/mifommkey.pem")
}

func getServer() *http.Server {
	server := &http.Server{
		Addr: ":8443",
	}
	return server
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request")

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	//log.Printf("Body: %+v", body)

	ar := v1beta1.AdmissionReview{}
	//log.Printf("AR: %+v", ar)

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		payload, err := json.Marshal(&v1beta1.AdmissionResponse{
			UID:     ar.Request.UID,
			Allowed: false,
			Result: &metav1.Status{
				Message: err.Error(),
			},
		})
		if err != nil {
			fmt.Println(err)
		}
		w.Write(payload)
	}
	//
	//log.Printf("AR: %+v", ar)
	//
	////object, _ := json.Marshal(ar.Request.Object)
	////log.Printf("Object: %s", string(object))
	//
	////payload, _ := json.Marshal(ar)
	//
	//message := "This is a test message"
	//
	//response := v1beta1.AdmissionResponse{
	//	UID:     ar.Request.UID,
	//	Allowed: false,
	//	Result: &metav1.Status{
	//		Message: message,
	//	},
	//}
	//
	////log.Printf("Response: %v", response)
	//json.Marshal(response)
	//ar.Response = &response
	//json.Marshal(ar)
	//
	//log.Printf("AR: %+v", ar)
	////w.Header().Set("Content-Type", "application/json")
	////w.Write(body)
	////w.Write(payload)
	//
	//serveJSON(w, ar)
	//

	response := v1beta1.AdmissionResponse{
		UID:     ar.Request.UID,
		Allowed: true,
	}

	if ar.Request.Kind.Kind == "Pod" {
		pod := v1.Pod{}
		json.Unmarshal(ar.Request.Object.Raw, &pod)
		log.Printf("Pod: %+v", pod)
		for _, container := range pod.Spec.Containers {
			if container.Name == "spam" {
				log.Println("BLOCK container from running...")
				log.Printf("Container resources: %+v", container.Resources)
				response.Allowed = false
				response.Result = &metav1.Status{
					Message: "no SPAM please!",
				}
				break
			} else {
				log.Println("Container is a-ok!")
			}
		}
	}

	//payload, err := json.Marshal(admitResponse)
	//if err != nil {
	//	log.Println(err)
	//}
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(payload)
	ar.Response = &response
	json.Marshal(ar)
	serveJSON(w, ar)

}

func serveJSON(w http.ResponseWriter, o interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}
