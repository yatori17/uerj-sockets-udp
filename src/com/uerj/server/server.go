package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type JsonEventString struct {
	Tipo string `json:"tipo"`
	Val  string `json:"val"`
}
type JsonEventInt struct {
	Tipo string      `json:"tipo"`
	Val  interface{} `json:"val"`
}

func main() {
	port := ":3000"
	protocol := "udp"
	var result JsonEventInt
	udpAddr, err := net.ResolveUDPAddr(protocol, port)
	if err != nil {
		fmt.Println("Wrong Address")
		return
	}

	conn, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	for {
		message := make([]byte, 4000)
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		data := strings.TrimSpace(string(message[:rlen]))
		error := json.Unmarshal(message[:rlen], &result)
		if error != nil {
			fmt.Println(err)
		}
		resultado := defineOperation(&result)
		result = nil
		fmt.Printf("received: %s from %s\n", data, remote)
		messageToSend := []byte(resultado)
		go serve(conn, remote, messageToSend[:])
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
	if value == strings.ToUpper(value) {
		value = strings.ToLower(value)
		return value
	}
	value = strings.ToUpper(value)
	return value
}

func defineOperation(obj1 *JsonEventInt) string {
	if obj1.Tipo == "string" {
		obj1.Val = inverteString(obj1.Val.(string))
	}
	if obj1.Tipo == "int" {
		obj1.Val = incrementValue(int(obj1.Val.(float64)))
		fmt.Println(obj1.Val)
	}
	if obj1.Tipo == "char" {
		obj1.Val = changeToLowerCaseOrUpperCase(obj1.Val.(string))
	}

	result, err := json.Marshal(obj1)

	if err == nil {
		fmt.Println(err)
	}
	return string(result)
}
