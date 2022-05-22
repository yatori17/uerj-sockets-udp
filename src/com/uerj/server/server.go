package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SocketConstruct struct {
	port    string
	address string
}
type JsonEvent struct {
	Tipo string      `json:"tipo"`
	Val  interface{} `json:"val"`
}

type peer struct {
	stop  func()
	since time.Time
}

var DIVISORIA string = "=================================================\n"
var MSG_IP string = "Informe o endereço IP que rodará seu servidor: "
var MSG_NUMERO_PORTA string = "Informe o Número da porta que rodará seu servidor: "
var MSG_WRONG_ADDRES string = "Wrong Address"
var protocol string = "udp"
var MESSAGE_ERROR_CLOSING_CONN = "Error in closing the UDP Connection: "
var MSG_SERVER_INIT = "Servidor UDP iniciado!"
var MSG_SERVER_LIST = "server listening "
var MSG_ERROR_CHAR = "Não foi possível converter o caracter!"
var MSG_ERROR_JSON = "Não foi possível transformar a mensagem para json"
var GENERATE_UDP_MESSAGE = "Gerando a mensagem para o cliente UDP "

func main() {
	var porta int
	var address string
	fmt.Print(MSG_IP)
	fmt.Scan(&address)
	fmt.Print(MSG_NUMERO_PORTA)
	fmt.Scan(&porta)
	conexao := address + ":" + strconv.Itoa(porta)
	var result JsonEvent
	peers := map[string]peer{}
	udpAddr, err := net.ResolveUDPAddr(protocol, conexao)
	if err != nil {
		fmt.Println(MSG_WRONG_ADDRES)
		return
	}

	conn, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		panic(err)
	}

	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(MESSAGE_ERROR_CLOSING_CONN, err)
		}
	}(conn)
	message := make([]byte, 2048)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Print(DIVISORIA + MSG_SERVER_INIT)
	fmt.Println(MSG_SERVER_LIST, conn.LocalAddr().String())
	fmt.Print(DIVISORIA)

	for {
		//recebendo a mensagem do cliente
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}
		//log de mensagem
		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Printf("received: %s from %s\n", data, remote)
		error := json.Unmarshal([]byte(data), &result)
		if error != nil {
			fmt.Println(error)
			continue
		}
		peer, ok := peers[remote.String()]
		if ok {
			continue
		}
		pctx, pcancel := context.WithCancel(ctx)
		peer.stop = pcancel
		peer.since = time.Now()
		peers[remote.String()] = peer
		resultado := defineOperation(&result, data)
		messageToSend := []byte(resultado)
		go generateMessageToClientUdp(messageToSend, pctx, conn, remote)
		for remote, p := range peers {
			if time.Since(p.since) > time.Minute {
				fmt.Println("Peer timedout")
				p.stop()
				delete(peers, remote)
			}
		}
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

func defineOperation(obj1 *JsonEvent, data string) string {
	if obj1.Tipo == "string" {
		obj1.Val = inverteString(obj1.Val.(string))
	}
	if obj1.Tipo == "int" {
		obj1.Val = incrementValue(int(obj1.Val.(float64)))
		fmt.Println(obj1.Val)
	}
	if obj1.Tipo == "char" {
		if obj1.Val == nil || len(obj1.Val.(string)) != 1 {
			println(MSG_ERROR_CHAR)
			return ""
		}
		obj1.Val = changeToLowerCaseOrUpperCase(obj1.Val.(string))
	}
	result, err := json.Marshal(obj1)
	if err != nil {
		fmt.Println(MSG_ERROR_JSON)
		fmt.Println(err)
	}
	clear(obj1)
	return string(result)
}

func clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func generateMessageToClientUdp(data []byte, ctx context.Context, conn *net.UDPConn, addr *net.UDPAddr) {
	fmt.Println(GENERATE_UDP_MESSAGE, addr)
	go func() {
		if len(data) == 0 {
			context.Background().Err()
		}
		serve(conn, addr, data)
		time.Sleep(time.Second * 1)
	}()
	<-ctx.Done()
	fmt.Println("Parando de escrever para cliente UDP", addr)
}
