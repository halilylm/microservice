package integration

import (
	"fmt"
	"github.com/halilylm/microservice/server"
	"net/http"
	"testing"
	"time"
)

func CreateServer() func() {
	fmt.Println("creating the server")
	srv := server.New(&server.Options{
		Host: "localhost",
		Port: 9000,
	})
	go func() {
		if err := srv.Start(); err != nil {
			panic(err)
		}
	}()
	for {
		_, err := http.Get("http://localhost:9000")
		if err == nil {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return func() {
		if err := srv.Stop(); err != nil {
			panic(err)
		}
	}
}

func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
}
