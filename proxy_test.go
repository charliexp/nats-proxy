package natsproxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/nats-io/nats"
)

func TestProxy(t *testing.T) {

	clientConn, _ := nats.Connect(nats.DefaultURL)
	natsClient := NewNatsClient(clientConn)
	natsClient.Subscribe("POST", "/test", Handler)
	// defer clientConn.Close()

	proxyConn, _ := nats.Connect(nats.DefaultURL)
	proxyHandler := NewNatsProxy(proxyConn)
	http.Handle("/", proxyHandler)
	// defer proxyConn.Close()

	log.Println("initializing proxy")
	go http.ListenAndServe(":3000", nil)
	time.Sleep(1 * time.Second)

	log.Println("Posting request")
	reader := bytes.NewReader([]byte("testData"))
	resp, err := http.Post("http://localhost:3000/test", "multipart/form-data", reader)
	if err != nil {
		log.Println(err)
		t.Error("Cannot do post")
		return
	}

	out, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(out))
	log.Println(resp.Header)
}

func Handler(c *Context) {
	log.Println("Getting request")
	log.Println(c.request.URL)
	c.request.Form.Get("email")

	respStruct := struct {
		User string
	}{
		"Radek",
	}

	bytes, _ := json.Marshal(respStruct)
	c.response.Body = bytes
	c.response.Header.Add("X-AUTH", "12345")
}
