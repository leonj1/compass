package compass

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/leonj1/compass/models"
	"github.com/orcaman/concurrent-map"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
)

func main() {
	var httpPort = flag.String("http-port", ":80", "Port which HTTP rest endpoint should listen on")
	flag.Parse()

	log.SetOutput(&lumberjack.Logger{
		Filename:   "/tmp/compass.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     3, //days
	})

	// Concurrent HashMap
	bar := cmap.New()
	clusters := &models.RootRouter{Clusters: &bar}

	s := mux.NewRouter()
	s.HandleFunc("/clusters", clusters.AddClusterHandler).Methods("POST")
	s.HandleFunc("/clusters/{name}", clusters.UpdateClusterHandler).Methods("PUT")
	s.HandleFunc("/clusters/{name}/crds", clusters.AddCustomResourceHandler).Methods("POST")
	s.HandleFunc("/clusters/{name}/crds/{crd}", clusters.UpdateCustomResourceHandler).Methods("PUT")
	s.HandleFunc("/clusters/{name}/nodes", clusters.AddNodeHandler).Methods("POST")
	s.HandleFunc("/clusters/{name}/nodes/{node}", clusters.UpdateNodeHandler).Methods("PUT")
	s.HandleFunc("/clusters/{name}/namespaces", clusters.AddNamespaceHandler).Methods("POST")
	//s.HandleFunc("/clusters/{name}/namespaces/{namespace}", clusters.SecureHandler).Methods("PUT")
	//s.HandleFunc("/clusters/{name}/events", clusters.SecureHandler).Methods("POST")
	//s.HandleFunc("/clusters/{name}/events/{event}", clusters.SecureHandler).Methods("PUT")

	//s.HandleFunc("/clusters", clusters.SecureHandler).Methods("GET")
	//s.HandleFunc("/clusters/{name}", clusters.SecureHandler).Methods("GET")

	log.Printf("Staring HTTPS service on %s .../n", *httpPort)
	if err := http.ListenAndServe(*httpPort, s); err != nil {
		panic(err)
	}
}
