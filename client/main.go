package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	connectNewWorker()
	time.Sleep(time.Minute)
}

func connectNewWorker() {
	go func() {
		rconn, err := createWorkerConn("test")
		if err != nil {
			log.Println(err)
			return
		}
		defer rconn.Close()

		// Connect to backend server
		conn, err := net.Dial("tcp", "golang.org:80")
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		go func() {
			defer conn.Close()
			defer rconn.Close()
			_, err = io.Copy(rconn, conn)
			if err != nil {
				log.Println(err)
				return
			}
		}()
		_, err = io.Copy(conn, rconn)
		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func createWorkerConn(id string) (net.Conn, error) {
	// Connect to reverse server
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}
	rconn, err := tls.Dial("tcp", "google.com:443", conf)
	if err != nil {
		return nil, err
	}

	// Send HTTP request
	req, err := http.NewRequest("GET", "https://google.com/", nil)
	if err != nil {
		return nil, err
	}
	err = req.Write(rconn)
	if err != nil {
		return nil, err
	}

	return rconn, nil
}
