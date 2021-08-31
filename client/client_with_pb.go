package main

import (
	"../pb"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Fatalln(err.Error())
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go handleClientConn(wg, conn)
	wg.Wait()
}

func handleClientConn(wg *sync.WaitGroup, conn net.Conn) {
	defer conn.Close()

	timeTicker := time.NewTicker(2 * time.Second)
	for i := 0; i < 100; i++ {
		str := "this is a go pb test " + time.Now().Format("2006-01-02 15:04:05")
		log.Println(str)
		msg := &pb.EchoMsg{
			Msg: str,
		}
		out, err := proto.Marshal(msg)
		if err != nil {
			log.Fatalln("marshal error: " + err.Error())
		}
		pkgBuf := make([]byte, len(out) + 4)
		binary.BigEndian.PutUint16(pkgBuf[0:], uint16(len(out) + 2))
		binary.BigEndian.PutUint16(pkgBuf[2:], uint16(pb.CmdId_ECHO))
		copy(pkgBuf[4:], out)
		conn.Write(pkgBuf)
		<-timeTicker.C
	}

	wg.Done()
}
