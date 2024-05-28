package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

var sender = false

// Client --
func Client() {
	register()
}

func register() {
	signalAddress := os.Args[2]

	localAddress := ":9595" // default port
	if len(os.Args) > 3 {
		localAddress = os.Args[3]
	}
	if len(os.Args) > 4 {
		sender = os.Args[4] == "sender"
	}

	remote, _ := net.ResolveUDPAddr("udp", signalAddress)
	local, _ := net.ResolveUDPAddr("udp", localAddress)
	conn, _ := net.ListenUDP("udp", local)
	go func() {
		bytesWritten, err := conn.WriteTo([]byte("register"), remote)
		if err != nil {
			panic(err)
		}

		fmt.Println(bytesWritten, " bytes written")
	}()

	listen(conn, local.String())
}

func listen(conn *net.UDPConn, local string) {
	bar := progressbar.DefaultBytes(8*1024*1024*20, "listening")
	for {
		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		bar.Add(bytesRead)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}

		if strings.Count(string(buffer[0:bytesRead]), ".") == 3 {
			for _, a := range strings.Split(string(buffer[0:bytesRead]), ",") {
				if a != local {
					go chatter(conn, a)
				}
			}
		}
	}
}

func chatter(conn *net.UDPConn, remote string) {
	addr, _ := net.ResolveUDPAddr("udp", remote)

	for {
		// fmt.Println("sent Hello! to ", remote)
		if sender {
			var buffer [500]byte
			rand.Read(buffer[:])
			conn.WriteTo(buffer[:], addr)
			time.Sleep(2 * time.Microsecond)
		} else {
			conn.WriteTo([]byte("Hello!"), addr)
			time.Sleep(100 * time.Millisecond)

		}
	}
}
