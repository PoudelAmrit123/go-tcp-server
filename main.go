package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type FileServer struct{}

func (fs FileServer) start() {

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {

		conn, err := lis.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		go fs.readLoop(conn)

	}

}

func (fs FileServer) readLoop(conn net.Conn) {

	// buf := make([]byte, 600)
	buf := new(bytes.Buffer)
	for {
		// n, err := conn.Read(buf)
		var size int64
		err := binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buf, conn, 1000)

		if err != nil {
			log.Fatal(err)

		}
		// file := buf[:n]
		fmt.Println("Reading the file ")
		fmt.Printf("Reading the file of size %d ", n)
		fmt.Println(buf.Bytes())
		// fmt.Printf("Reading the file of size %s", file)

	}

}

func sendFile(size int) error {

	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, size)
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))

	// n, err := conn.Write(file)
	if err != nil {
		return err
	}

	fmt.Printf("Written %d byted over the network", n)
	fmt.Println(file)

	return nil

}

func main() {

	go func() {
		sendFile(10)

	}()
	fmt.Println("Hello! File streaming ")
	Server := &FileServer{}
	Server.start()

}
