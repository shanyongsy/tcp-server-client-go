package main

import (
	"../pb"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8001")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer l.Close()
	log.Println("server listen on 8001")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error: " + err.Error())
			continue
		}
		log.Println("accept client: " + conn.RemoteAddr().String())
		go handleConnectionWithPb(conn)
	}
}

func handleConnectionWithPb(conn net.Conn) {
	defer conn.Close()

	for {
		lenBuf := make([]byte, 2)
		_, err := io.ReadFull(conn, lenBuf)
		if err != nil {
			if err == io.EOF {
				log.Println("client quit...")
			} else {
				log.Println("read length error: " + err.Error())
			}
			break
		}
		pkgLen := binary.BigEndian.Uint16(lenBuf[0:])
		buf := make([]byte, pkgLen)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			log.Println("read package error: " + err.Error())
			break
		}
		handlePackage(conn, buf)
	}
}

func handlePackage(conn net.Conn, buf []byte) {
	cmd := binary.BigEndian.Uint16(buf[0:])
	if cmd == uint16(pb.CmdId_ECHO) {
		msg := &pb.EchoMsg{}
		if err := proto.Unmarshal(buf[2:], msg); err != nil {
			log.Println("unmarshal error: " + err.Error())
		}
		log.Println(msg.GetMsg())
		//conn.Write([]byte(fmt.Sprintf("%s: %s", "server response", time.Now().Format("2006-01-02 15:04:05"))))
		conn.Write([]byte("server response\n"))
	}
}
