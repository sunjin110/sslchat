package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"sslchat/pkg/common/chk"
	"strings"
)

var sc = bufio.NewScanner(os.Stdin)

const closeErrMsg = "use of closed network connection"

// client クライアント認証を使用してみる
func main() {

	// tls setting
	cert, err := tls.LoadX509KeyPair("../../cert.pem", "../../key.pem")
	chk.SE(err)
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true} // オレオレ証明書のためInsecureSkipをtrue

	// make connection
	conn, err := tls.Dial("tcp4", "localhost:1234", tlsConfig)
	chk.SE(err)
	defer conn.Close()

	// serverからのWriteを監視する
	go watchChat(conn)

	// 書き込み
	for {
		// 標準入力をまつ
		if sc.Scan() {

			msg := sc.Text()
			if msg == "q" {
				conn.Close()
				break
			}

			conn.Write([]byte(sc.Text()))
		}
	}

	fmt.Println("logout...")
}

// watchChat チャットをReadし続けるもの
func watchChat(conn net.Conn) {
	var buf [1024]byte
	for {
		i, err := conn.Read(buf[:])
		if err != nil {
			if strings.Contains(err.Error(), closeErrMsg) {
				fmt.Println("connect close...")
				return
			}
			chk.SE(err)
		}
		fmt.Println(string(buf[:i]))
	}
}
