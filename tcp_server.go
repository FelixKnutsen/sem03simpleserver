package main

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/FelixKnutsen/is105sem03/mycrypt"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekryptert melding: ", string(dekryptertMelding))
					switch msg := string(dekryptertMelding); msg {
					case "ping":
						_, err = c.Write([]byte(mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)))
					default:
						_, err = c.Write(mycrypt.Krypter([]rune(msg), mycrypt.ALF_SEM03, 4))
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
