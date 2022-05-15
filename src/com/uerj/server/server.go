package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	for {
		message := make([]byte, 20)
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Printf("received: %s from %s\n", data, remote)
		messageToSend := []byte(changeToLowerCaseOrUpperCase(data))
		go serve(conn, remote, messageToSend[:len(messageToSend)])
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	// 0 - 1: ID
	// 2: QR(1): Opcode(4)
	//buf[2] |= 0x80 // Set QR bit

	pc.WriteTo(buf, addr)
}

func inverteString(value string) string {
	runes := []rune(value)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func incrementValue(value int) int {
	return value + 1
}

func changeToLowerCaseOrUpperCase(value string) string {
	if (value == strings.ToUpper(value)){
		value = strings.ToLower(value)
		return value
	}
	value = strings.ToUpper(value)
	return value
}