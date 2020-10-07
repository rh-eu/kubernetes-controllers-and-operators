package app

import (
	"flag"
	"log"
	"net/http"

	//kitlog "github.com/go-kit/kit/log"
	"github.com/julienschmidt/httprouter"

	"github.com/rh-eu/kubernetes-controllers-and-operators/pkg/handler"
	//admissioncontrol "github.com/rh-eu/kubernetes-controllers-and-operators/pkg/admission-control"
	"github.com/rh-eu/kubernetes-controllers-and-operators/pkg/helper"
)

type conf struct {
	TLSCertPath string
	TLSKeyPath  string
	//HTTPOnly    bool
	Port string
	Host string
}

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
	//var logger kitlog.Logger
	//logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	//log.SetOutput(kitlog.NewStdlibAdapter(logger))
	//logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "loc", kitlog.DefaultCaller)

	k := &App{
		r: httprouter.New(),
	}

	router := k.r

	router.GET("/hello", myHandler)

	router.POST("/api/user/create", helper.StdToJulienMiddleware(helper.StdHandler()))
	router.GET("/api/user/create", helper.JulienToJulienMiddleware(helper.JulienHandler()))

	router.GET("/mytest", helper.MyTestToJulienMiddleware(helper.MyTestHandler()))

	//router.POST("/mutation", helper.MyTestToJulienMiddleware(helper.MutatingHandler()))

	router.POST("/mutation", helper.MyTestToJulienMiddleware(&handler.AdmissionHandler{
		AdmitFunc: handler.MutatingLittleTesting(),
	}))

	router.POST("/admission-control/enforce-pod-annotations", helper.MyTestToJulienMiddleware(&handler.AdmissionHandler{
		AdmitFunc: handler.EnforcePodAnnotations(
			[]string{"kube-system"},
			map[string]func(string) bool{
				"k8s.questionable.services/hostname": func(string) bool { return true },
			}),
		//Logger: logger,
	}))

	return k
}

// Run ...
func (k *App) Run() {

	// Get config
	conf := &conf{}
	flag.StringVar(&conf.TLSCertPath, "cert-path", "./certs/mifomm.mutation.svc/cert.pem", "The path to the PEM-encoded TLS certificate")
	flag.StringVar(&conf.TLSKeyPath, "key-path", "./certs/mifomm.mutation.svc/key.pem", "The path to the unencrypted TLS key")
	//flag.BoolVar(&conf.HTTPOnly, "http-only", false, "Only listen on unencrypted HTTP (e.g. for proxied environments)")
	flag.StringVar(&conf.Port, "port", ":8443", "The port to listen on (HTTPS).")
	flag.StringVar(&conf.Host, "host", "admissiond.questionable.services", "The hostname for the service")
	flag.Parse()

	log.Printf("app is up and running and listening on port %s", conf.Port)
	log.Fatal(http.ListenAndServeTLS(conf.Port, conf.TLSCertPath, conf.TLSKeyPath, k.r))
}
