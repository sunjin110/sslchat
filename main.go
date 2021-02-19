package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sslchat/pkg/common/chk"
	"sync/atomic"
)

type User struct {
	conn net.Conn
	name string
}

// userMap
// TODO mulock
var userMap map[uint64]*User = map[uint64]*User{}

// userID
var userID uint64

func getAtomicID() uint64 {
	return atomic.AddUint64(&userID, 1)
}

func main() {

	// tls setting
	cert, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	chk.SE(err)
	tlsConfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	// make listener
	listener, err := tls.Listen("tcp4", "localhost:1234", tlsConfg)
	chk.SE(err)
	defer listener.Close()

	for {
		log.Println("Wating for clients")
		conn, err := listener.Accept()
		chk.SE(err)
		go handle(conn)
	}

}

func handle(conn net.Conn) {

	defer func() {
		conn.Close()
	}()

	conn.Write([]byte("hello"))
	conn.Write([]byte("Please input your name>"))
	var buf [1024]byte

	nameLength, err := conn.Read(buf[:]) // ここで、bufよりも大きな名前が来てしまった場合あふれるからどうにかする必要があるのか
	chk.SE(err)
	name := string(buf[:nameLength])
	id := getAtomicID()
	userMap[id] = &User{
		conn: conn,
		name: name,
	}
	conn.Write([]byte(name + "さん、ようこそ"))

	// 全serverに通知
	pushMessage(id, name, "<入出しました>")

	for {
		i, err := conn.Read(buf[:]) // buf以上のものが書き込まれたときに、途切れてしまう、その場合は、なにか指定した文字列が入るまでは連続した文字として扱えばいい

		if err != nil {
			if err.Error() == "EOF" { // connect end

				delete(userMap, id)
				pushMessage(id, name, "<退出しました>")
				return
			}
		}

		chk.SE(err)
		go pushMessage(id, name, string(buf[:i]))
	}
}

// 全ユーザーにメッセージをプッシュする
func pushMessage(myID uint64, myName string, msg string) {

	for id, user := range userMap {

		// 自分には送信しない
		if myID == id {
			continue
		}

		// 書き込み
		user.conn.Write([]byte(fmt.Sprintf("%s:%s", myName, msg)))
	}

}
