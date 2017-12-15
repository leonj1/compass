package models

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/orcaman/concurrent-map"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type RootRouter struct {
	sync.RWMutex
	Clusters *cmap.ConcurrentMap
}

func (m *RootRouter) AddClusterHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	payload := &Cluster{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem unmarshalling cluster"}
		respondWithJSON(w, 404, response)
		return
	}

	err = m.AddCluster(*payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem adding cluster: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 201, &QueryResponse{Message: "created"})
}

func (m *RootRouter) UpdateClusterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["name"]
	body, err := ioutil.ReadAll(r.Body)
	payload := &Cluster{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem unmarshalling cluster"}
		respondWithJSON(w, 404, response)
		return
	}

	err = m.UpdateCluster(clusterName, *payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem updating cluster: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 200, &QueryResponse{Message: "updated"})
}

func (m *RootRouter) AddCustomResourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["name"]
	body, err := ioutil.ReadAll(r.Body)
	payload := &Crd{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem unmarshalling custom resource"}
		respondWithJSON(w, 404, response)
		return
	}

	err = m.AddCustomResource(clusterName, *payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem adding custom resource: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 201, &QueryResponse{Message: "created"})
}

func (m *RootRouter) UpdateCustomResourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["name"]
	crdName := vars["crd"]
	body, err := ioutil.ReadAll(r.Body)
	payload := &Crd{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem unmarshalling custom resource"}
		respondWithJSON(w, 404, response)
		return
	}

	err = m.UpdateCustomResource(clusterName, crdName, *payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem updating custom resource: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 200, &QueryResponse{Message: "updated"})
}

func (m *RootRouter) AddNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["name"]
	body, err := ioutil.ReadAll(r.Body)
	payload := &Node{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem unmarshalling node"}
		respondWithJSON(w, 404, response)
		return
	}

	err = m.AddNode(clusterName, *payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem adding node: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 201, &QueryResponse{Message: "created"})
}

func (m *RootRouter) UpdateNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["name"]
	nodeName := vars["node"]
	body, err := ioutil.ReadAll(r.Body)
	payload := &Node{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem unmarshalling node"}
		respondWithJSON(w, 404, response)
		return
	}

	err = m.UpdateNode(clusterName, nodeName, *payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem updating node: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 200, &QueryResponse{Message: "updated"})
}

func (m *RootRouter) AddNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["name"]
	body, err := ioutil.ReadAll(r.Body)
	payload := &Namespace{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem unmarshalling namespace"}
		respondWithJSON(w, 404, response)
		return
	}

	err = m.AddNamespace(clusterName, *payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem adding namespace: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 201, &QueryResponse{Message: "created"})
}

func (m *RootRouter) UpdateNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["name"]
	namespace := vars["namespace"]
	body, err := ioutil.ReadAll(r.Body)
	payload := &Namespace{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem unmarshalling namespace"}
		respondWithJSON(w, 404, response)
		return
	}

	err = m.UpdateNamespace(clusterName, namespace, *payload)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem updating namespace: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 200, &QueryResponse{Message: "updated"})
}

func (m *RootRouter) SetEventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterName := vars["name"]
	body, err := ioutil.ReadAll(r.Body)
	err = m.SetEvents(clusterName, string(body))
	if err != nil {
		response := &ErrorResponse{ErrorMessage: fmt.Sprintf("Problem adding events: %s", err.Error())}
		respondWithJSON(w, 404, response)
		return
	}

	respondWithJSON(w, 201, &QueryResponse{Message: "set"})
}

func (m *RootRouter) GetAllClustersHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, m.Clusters)
}

type PostResponse struct {
	Digest string `json:"digest"`
}

type QueryResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"err_msg"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
