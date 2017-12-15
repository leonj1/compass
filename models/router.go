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
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Message not found"}
		respondWithJSON(w, 404, response)
		return
	}
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
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Message not found"}
		respondWithJSON(w, 404, response)
		return
	}
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
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Message not found"}
		respondWithJSON(w, 404, response)
		return
	}
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
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Message not found"}
		respondWithJSON(w, 404, response)
		return
	}
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
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Message not found"}
		respondWithJSON(w, 404, response)
		return
	}
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
