package main

import (
	"bufio"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	. "server_go/funcs"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	_ "github.com/lib/pq"
	godpi "github.com/mushorg/go-dpi"
	"github.com/mushorg/go-dpi/types"
)

type patients struct {
	id          int
	fio         string
	birthdate   string
	mestoprojiv string
	number      string
	email       string
}

var signals []string
var (
	count, idCount int
	protoCounts    map[types.Protocol]int
	packetChannel  <-chan gopacket.Packet
	err            error
)

func main() {
	loadFile, err := os.ReadFile("ssl_medassist.pcap")
	_ = loadFile
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":2025", config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	protoCounts = make(map[types.Protocol]int)

	device := flag.String("device", "", "Device to watch for packets")

	flag.Parse()

	if *device != "" {

		handle, deverr := pcap.OpenLive(*device, 1024, false, time.Duration(-1))
		if deverr != nil {
			fmt.Println("Error opening device:", deverr)
			return
		}
		packetChannel = gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	}

	initErrs := godpi.Initialize()
	if len(initErrs) != 0 {
		for _, err := range initErrs {
			fmt.Println(err)
		}
		return
	}
	fmt.Println("Init done")

	defer func() {
		godpi.Destroy()

	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	intSignal := false

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	count = 0
	for packet := range packetChannel {
		fmt.Printf("Packet #%d: ", count+1)
		flow, isNew := godpi.GetPacketFlow(packet)
		result := godpi.ClassifyFlow(flow)
		if result.Protocol != types.Unknown {

			GetSignal([]string{result.String()})
		} else {
			fmt.Print("Could not identify")
		}
		if isNew {
			fmt.Println(" (new flow)")
		} else {
			fmt.Println()
		}

		select {
		case <-signalChannel:
			fmt.Println("Received interrupt signal")
			intSignal = true
		default:
		}
		if intSignal {
			break
		}
		count++
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		newmsg, err := hex.DecodeString(msg)
		n, err := conn.Write(newmsg)
		if err != nil {
			log.Println(n, err)
			return
		}

	}

}
func GetSignal(args []string) {
	switch args[0] {
	case "GetRoles":
		val := GetRoles()
		_ = val
	case "GetPatients":
		val := GetPatients()
		_ = val
	case "GetUsers":
		val := GetUsers()
		_ = val
	case "GetOne":
		val := GetOne("login")
		_ = val
	case "CreateUser":
		CreateUser()
	case "CreateRole":
		CreateRole()
	case "GetStorage":
		val := GetStorage()
		_ = val

	}

}
