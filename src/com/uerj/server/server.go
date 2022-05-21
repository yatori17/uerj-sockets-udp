t adpackage main

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strings"
)

type SocketConstruct struct {
	port    string
	address string
}
type JsonEvent struct {
	Tipo string      `json:"tipo"`
	Val  interface{} `json:"val"`
}

func main() {
	port := ":3000"
	protocol := "udp"
	address := "0.0.0.0"
	var result JsonEvent
	udpAddr, err := net.ResolveUDPAddr(protocol, address+port)
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
		message := make([]byte, 2048)
		//recebendo a mensagem do cliente
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}
		//log de mensagem
		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Printf("received: %s from %s\n", data, remote)

		error := json.Unmarshal(message[:rlen], &result)
		if error != nil {
			fmt.Println(err)
			continue
		}
		resultado := defineOperation(&result)
		messageToSend := []byte(resultado)
		go serve(conn, remote, messageToSend)
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
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

func defineOperation(obj1 *JsonEvent) string {
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
	clear(obj1)
	return string(result)
}

func clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}
