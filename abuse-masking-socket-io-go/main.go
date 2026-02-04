package main

import (
	"abuse-masking-go/abuse-masker"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"os"
	"github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main(){
	AbuseTrie, err := masker.LoadAbuseTrie("abuse_words.txt")
	if err != nil{
        log.Fatal("Failed to load abuse words:", err)
	}
	
	server := socketio.NewServer(&engineio.Options{
		PingTimeout:  60 * time.Second,
		PingInterval: 25 * time.Second,
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: func(r *http.Request) bool { return true },
			},
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool { return true },
			},
		},
	})
	server.OnConnect("/",func(s socketio.Conn) error {
		log.Println("connected:",s.ID())
		return nil
	})

	server.OnEvent("/","join_room",func (s socketio.Conn, room string) {
		s.Join(room)
		log.Printf("Client %s joined room %s\n",s.ID(), room)
		server.BroadcastToRoom("",room,"server_msg",s.ID()+" joined the room")
	})

	server.OnEvent("/","leave_room",func (s socketio.Conn, room string) {
		s.Leave(room)
		log.Printf("Client %s left room %s\n",s.ID(),room)
        server.BroadcastToRoom("",room,s.ID() + " left the room")
	})

	server.OnEvent("/", "send_msg", func(s socketio.Conn, room, msg string) {
		log.Printf("Message from %s to %s: %s\n", s.ID(), room, msg)
		maskedMsg := masker.MaskText(msg,AbuseTrie)
		server.BroadcastToRoom("", room, "receive_msg", s.ID(), maskedMsg)
	})

	server.OnError("/",func(s socketio.Conn, err error) {
		log.Println("Socket error:", err)
	})

	server.OnDisconnect("/",func(c socketio.Conn, reason string) {
		log.Println("Client disconnected: ",reason)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", corsMiddleware(server))
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	type AddAbuseReq struct{
		Abuse string `json:"abuse"`
	}
	http.HandleFunc("/add-abuse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var req AddAbuseReq
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		abuse := req.Abuse
		if abuse == "" {
			http.Error(w, "abuse field is required", http.StatusBadRequest)
			return
		}
		masker.UpdateTrie(AbuseTrie, abuse)
		f, err := os.OpenFile("abuse_words.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Failed to persist abuse word:", err)
		} else {
			f.WriteString("\n" + abuse)
			f.Close()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "added", "word": abuse})
	})
	log.Println("Server started on :9000")
	log.Fatal(http.ListenAndServe(":9000",nil))
}