package app

import (
	"log"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/julienschmidt/httprouter"
	admissioncontrol "github.com/rh-eu/kubernetes-controllers-and-operators/pkg/admission-control"
	"github.com/rh-eu/kubernetes-controllers-and-operators/pkg/helper"
)

// App ...
type App struct {
	//c Config
	r *httprouter.Router
}

func myHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Handling request")
	w.Write([]byte("Hello World!"))
}

// NewApp ...
func NewApp() *App {

	// Set up logging
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "loc", kitlog.DefaultCaller)

	k := &App{
		r: httprouter.New(),
	}

	router := k.r

	router.GET("/hello", myHandler)

	router.POST("/api/user/create", helper.StdToJulienMiddleware(helper.StdHandler()))
	router.GET("/api/user/create", helper.JulienToJulienMiddleware(helper.JulienHandler()))

	router.GET("/mytest", helper.MyTestToJulienMiddleware(helper.MyTestHandler()))

	router.GET("/admissioncontrol", helper.MyTestToJulienMiddleware(&admissioncontrol.AdmissionHandler{
		AdmitFunc: admissioncontrol.EnforcePodAnnotations(
			[]string{"kube-system"},
			map[string]func(string) bool{
				"k8s.questionable.services/hostname": func(string) bool { return true },
			}),
		Logger: logger,
	}))

	return k
}

// Run ...
func (k *App) Run() {
	log.Printf("app is up and running")
	log.Fatal(http.ListenAndServe(":5051", k.r))
}
