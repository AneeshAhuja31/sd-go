package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sha-go/hash"
	"sha-go/node"
	"sha-go/ring"
)

type GetRequest struct {
	Key string `json:"key"`
}

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DeleteRequest struct {
	Key string `json:"key"`
}

type Server struct {
	Node *node.Node
	Ring *ring.Ring
}

func SetupRouter(srv *Server) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		srv.GetHandler(w, r)
	})
	mux.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		srv.PutHandler(w, r)
	})
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		srv.DeleteHandler(w, r)
	})

	return mux
}

func (srv *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var req GetRequest

	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := req.Key

	targetNode := ring.FindNode(key, *srv.Ring, 16)

	if targetNode.ID == srv.Node.ID {
		slotQuery := fmt.Sprintf("SELECT value, hash FROM files_%d WHERE key = $1", srv.Node.Slot)

		var value string
		var hash uint64

		err := srv.Node.DB.QueryRow(slotQuery, key).Scan(&value, &hash)
		if err == sql.ErrNoRows {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		response := map[string]string{"key": key, "value": value}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		srv.forwardRequest(w, r, targetNode, "/get")
	}
}

func (srv *Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var req PutRequest

	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := req.Key
	value := req.Value

	targetNode := ring.FindNode(key, *srv.Ring, 16)

	if targetNode.ID == srv.Node.ID {
		keyHash := hash.Hash(key)
		insertQuery := fmt.Sprintf("INSERT INTO files_%d (key, value, hash) VALUES ($1, $2, $3) ON CONFLICT (key) DO UPDATE SET value = $2, hash = $3", srv.Node.Slot)

		_, err := srv.Node.DB.Exec(insertQuery, key, value, keyHash)
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		response := map[string]string{"status": "success", "key": key}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		srv.forwardRequest(w, r, targetNode, "/put")
	}
}

func (srv *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var req DeleteRequest

	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := req.Key

	targetNode := ring.FindNode(key, *srv.Ring, 16)

	if targetNode.ID == srv.Node.ID {
		deleteQuery := fmt.Sprintf("DELETE FROM files_%d WHERE key = $1", srv.Node.Slot)

		result, err := srv.Node.DB.Exec(deleteQuery, key)
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		response := map[string]string{"status": "deleted", "key": key}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		srv.forwardRequest(w, r, targetNode, "/delete")
	}
}

func (srv *Server) forwardRequest(w http.ResponseWriter, r *http.Request, targetNode *node.Node, endpoint string) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}

	targetURL := fmt.Sprintf("http://%s%s", targetNode.HttpServer.Addr, endpoint)

	proxyReq, err := http.NewRequest(r.Method, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Failed to forward request to %s: %v", targetURL, err)
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
