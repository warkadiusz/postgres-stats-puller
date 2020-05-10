package http_service

import (
	"ask/db"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type HttpService struct {
	clients    []*websocket.Conn
	clientsMux sync.Mutex
}

type Response struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   string `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(_ *http.Request) bool {
	return true
}

func CreateHttpService() *HttpService {
	hs := HttpService{}
	hs.clients = make([]*websocket.Conn, 0, 5)
	hs.clientsMux = sync.Mutex{}
	httpServer := mux.NewRouter()
	httpServer.HandleFunc("/data/{name}", hs.serveData)
	httpServer.HandleFunc("/ws", hs.serveWS)

	go http.ListenAndServe(":8090", httpServer)
	return &hs
}

func (hs *HttpService) serveData(w http.ResponseWriter, req *http.Request) {
	requiredParameters := []string{"before", "after"}
	parametersValues := make(map[string]string)
	dataName := mux.Vars(req)["name"]

	for _, param := range requiredParameters {
		paramValue, ok := req.URL.Query()[param]

		if !ok || len(paramValue[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(jsonify(Response{Status: "error", Code: 400, Data: "\"" + param + "\" parameter is missing"}))
			return
		}
		parametersValues[param] = paramValue[0]
	}

	z, _ := time.Now().Local().Zone()

	before, err := time.Parse("2006-01-02T15:04:05 MST", parametersValues["before"]+" "+z)
	after, err2 := time.Parse("2006-01-02T15:04:05 MST", parametersValues["after"]+" "+z)

	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(jsonify(Response{Status: "error", Code: 400, Data: "\"after\" or \"before\" is not a valid date time"}))
		return
	}

	//db.QueryData("")
	//log.Print(dataName, before, after)
	data, err := db.QueryData(dataName, after, before)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonify(Response{Status: "error", Code: 500, Data: "Internal server error"}))
	}

	responseDataPoints := map[string]float32{}

	for _, dataPoint := range data.Series[0].Values {
		val := 0.0
		if dataPoint[1] != nil {
			val, _ = dataPoint[1].(json.Number).Float64()
		}

		responseDataPoints[dataPoint[0].(string)] = float32(val)
	}

	byteArrDataPoints, _ := json.Marshal(responseDataPoints)
	stringifiedDataPoints := string(byteArrDataPoints)

	response := Response{
		Status: "success",
		Code:   200,
		Data:   stringifiedDataPoints,
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonify(response))
}

func (hs *HttpService) serveWS(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	hs.appendClient(conn)
	defer hs.removeClient(conn)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			hs.removeClient(conn)
			break
		}

		log.Printf("recv: %s", message)
	}
}

func (hs *HttpService) BroadcastDatapoint(name string, value int) {
	defer hs.clientsMux.Unlock()

	data, _ := json.Marshal(map[string]int{name: value})
	dataString := string(data)

	hs.clientsMux.Lock()
	log.Print(hs.clients)
	for _, client := range hs.clients {
		_ = client.WriteMessage(websocket.TextMessage, jsonify(Response{
			Status: "datapoint",
			Code:   200,
			Data:   dataString,
		}))
	}
}

func (hs *HttpService) appendClient(conn *websocket.Conn) {
	hs.clientsMux.Lock()
	defer hs.clientsMux.Unlock()
	hs.clients = append(hs.clients, conn)
}

func (hs *HttpService) removeClient(conn *websocket.Conn) {
	hs.clientsMux.Lock()
	defer hs.clientsMux.Unlock()

	for index, clientToCheck := range hs.clients {
		if clientToCheck == conn {
			hs.clients[index] = hs.clients[len(hs.clients)-1]
			hs.clients[len(hs.clients)-1] = nil
			hs.clients = hs.clients[:len(hs.clients)-1]
			conn.Close()
		}
	}
}

func jsonify(response Response) []byte {
	msg, _ := json.Marshal(response)
	return msg
}
