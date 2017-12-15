package models

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/orcaman/concurrent-map"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/kataras/go-errors"
	"io/ioutil"
	"net/http"
	"sync"
)

type RootRouter struct {
	sync.RWMutex
	Clusters *cmap.ConcurrentMap
}

func (m *RootRouter) AddCluster(cluster Cluster) error {
	if m.Clusters.Has(cluster.Name) {
		return errors.New("cluster already exists")
	}
	m.Clusters.Set(cluster.Name, cluster)
	return nil
}

func (m *RootRouter) UpdateCluster(clusterName string, cluster Cluster) error {
	if !m.Clusters.Has(cluster.Name) {
		return errors.New("cluster does not exist")
	}

	if existing, ok := m.Clusters.Get(clusterName); ok {
		current := existing.(Cluster)
		current.Name = cluster.Name
		current.Status = cluster.Status
		current.Personality = cluster.Personality
		m.Clusters.Set(cluster.Name, current)
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}

func (m *RootRouter) AddCustomResource(clusterName string, crd Crd) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if existing, ok := m.Clusters.Get(clusterName); ok {
		current := existing.(Cluster)
		if current.Crds.Has(crd.Name) {
			return errors.New("crd already exists in cluster")
		}
		current.Crds.Set(crd.Name, crd)
		m.Clusters.Set(clusterName, current)
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
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

func (m *RootRouter) SecureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	w.Header().Set("Content-Type", "application/json")
	//if value, ok := m.cMap.Get(hash); ok {
	//	valueAsString := value.(string)
	//	response := &QueryResponse{Message: string(valueAsString)}
	//	f, _ := json.Marshal(response)
	//	w.Write([]byte(f))
	//	return
	//}
	response := &ErrorResponse{ErrorMessage: "Message not found"}
	respondWithJSON(w, 404, response)
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
