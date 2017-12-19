package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/leonj1/compass/models"
	"github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	var httpPort = flag.String("http-port", "80", "Port which HTTP rest endpoint should listen on")
	flag.Parse()

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
	s.HandleFunc("/clusters/{name}/namespaces/{namespace}", clusters.UpdateNamespaceHandler).Methods("PUT")
	s.HandleFunc("/clusters/{name}/events", clusters.SetEventsHandler).Methods("POST")

	s.HandleFunc("/clusters", clusters.GetAllClustersHandler).Methods("GET")
	s.HandleFunc("/clusters/{name}", clusters.GetOneClustersHandler).Methods("GET")

	s.HandleFunc("/clusters/{name}", clusters.DeleteClustersHandler).Methods("DELETE")

	s.HandleFunc("/public/health", clusters.HealthCheckHandler).Methods("GET")

	port := fmt.Sprintf(":%s", *httpPort)
	log.Infof("Staring HTTPS service on %s ...\n", port)
	if err := http.ListenAndServe(port, s); err != nil {
		panic(err)
	}
}
