package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"
)

func Test_DownloadInServerOk(t *testing.T) {

	host := "localhost"
	port := rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err := startServer(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(time.Second * 2)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	_, err = writer.Write([]byte("DOWNLOAD\n"))
	if err != nil {
		t.Fatalf("can't comand %s to server: %v", "DOWNLOAD", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("can't comand %s to server: %v", "DOWNLOAD", err)
		return
	}
	_, err = writer.Write([]byte("example\n"))
	if err != nil {
		t.Fatalf("can't comand %s to server: %v", "example", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("can't comand %s to server: %v", "example", err)
		return
	}
	file, err := os.Create("example")
	if err != nil {
		t.Fatalf("Can't open file: %v", err)
	}
	defer file.Close()
	bytees, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("Can't download file: %v", err)
	}

	_, err = file.Write(bytees)
	if err != nil {
		t.Fatalf("Can't write to file: %v", err)
	}
	src, err := ioutil.ReadFile("downloads/example")
	if err != nil {
		t.Fatalf("Can't read file src: %v", err)
	}
	src1, err := ioutil.ReadFile("example")

	if err != nil {
		t.Fatalf("Can't read file src1: %v", err)
	}

	if !bytes.Equal(src, src1) {
		t.Fatalf("files are not equal: %v", err)
	}
}

func Test_UploadToServerOk(t *testing.T) {

	host := "localhost"
	port := rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err := startServer(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(time.Second * 2)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	_ = bufio.NewReader(conn)
	_, err = writer.Write([]byte("UPLOAD\n"))
	if err != nil {
		t.Fatalf("can't comand %s to server: %v", "UPLOAD", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("can't comand %s to server: %v", "UPLOAD", err)
		return
	}
	_, err = writer.Write([]byte("example2\n"))
	if err != nil {
		t.Fatalf("can't comand %s to server: %v", "example2", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("can't comand %s to server: %v", "example2", err)
		return
	}
	file, err := os.Open("example2")
	if err != nil {
		t.Fatalf("Can't open file: %v", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			t.Fatalf("Can't open file: %v", err)
		}
	}()

	bytees, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Can't download file: %v", err)
	}

	_, err = writer.Write(bytees)
	if err != nil {
		t.Fatalf("Can't write to server: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("Can't write to server: %v", err)
	}

	src, err := ioutil.ReadFile("downloads/example2")
	if err != nil {
		t.Fatalf("Can't read file src: %v", err)
	}
	src1, err := ioutil.ReadFile("example2")

	if err != nil {
		t.Fatalf("Can't read file src1: %v", err)
	}

	if !bytes.Equal(src, src1) {
		t.Fatalf("files are not equal: %v", err)
	}
}

func Test_ListInServerOk(t *testing.T) {
	host := "localhost"
	port := rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err := startServer(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(time.Second * 2)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	_, err = writer.Write([]byte("LIST\n"))
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", "LIST", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", "LIST", err)
	}
	reader := bufio.NewReader(conn)
	counter := ""
	for {
		readString, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("can't read frocm loalhost %v", err)
		}
		counter += readString
	}
	fmt.Println(counter)
	result := `example
example2
`
	if result != counter{
		t.Fatalf("is not equal")
	}
}

