package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	file, err := os.OpenFile("server-log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("can't open log file %e", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("can't close log file %e", err)
		}
	}()

	log.SetOutput(file)
	log.Print("start application\n")
	host := "0.0.0.0"
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9999"
	}

	err = startServer(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Printf("some error in server %e", err)
		fmt.Println("Произашло ошибка сервера")
	}
}

func startServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("can't listen: %v", err)
	}
	defer func() {
		err = listener.Close()
		if err != nil {
			err = fmt.Errorf("can't close Listener %e", err)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept client %e", err)
			fmt.Println("Ошибка подключения клиента")
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("can't close connection %e", err)
		}
	}()

	reader, write := bufio.NewReader(conn), bufio.NewWriter(conn)

	str, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("can't read command %e", err)
		fmt.Println("Не правилная команда")
		return
	}
	dir, err := ioutil.ReadDir("downloads")
	if err != nil {
		log.Printf("can't read directory <downloads> %e", err)
		fmt.Println("Не удалось прочесть директорию downloads")
		return
	}
	if str == "LIST\n" {
		handleList(write, dir)
		return
	}
	if str == "DOWNLOAD\n" {
		handleDownload(reader,write, dir)
		return
	}
	if str == "UPLOAD\n" {
		handleUpload(reader)
		return
	}
}

func handleList(write *bufio.Writer, dir []os.FileInfo) {
	for _, info := range dir {
		if !info.IsDir() {
			_, err := write.Write([]byte(info.Name() + "\n"))
			if err != nil {
				log.Printf("can't write to client %e", err)
				fmt.Println("Не удалось написать клиенту")
				return
			}
			err = write.Flush()
			if err != nil {
				log.Printf("can't write to client %e", err)
				fmt.Println("Не удалось написать клиенту")
				return
			}
		}
	}
	return
}

func handleDownload(reader *bufio.Reader, write *bufio.Writer, dir []os.FileInfo) {
	fileName, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("can't read command %e", err)
		fmt.Println("Не удалочь прочесть команду клиента")
		return
	}
	if fileName == "" {
		return
	}

	for _, info := range dir {
		if !info.IsDir() && info.Name()+"\n" == fileName {

			bytes, err := ioutil.ReadFile("downloads/" + info.Name())
			if err != nil {
				log.Printf("can't read file %e", err)
				fmt.Println("Не удалось прочесть файл")
				return
			}
			_, err = write.Write(bytes)
			if err != nil {
				log.Printf("can't write to client %e", err)
				fmt.Println("Не удалочь отправить файл клиенту")
				return
			}
			err = write.Flush()
			if err != nil {
				log.Printf("can't write to client %e", err)
				fmt.Println("Не удалочь отправить файл клиенту")
				return
			}
			return
		}
	}
}


func handleUpload(reader *bufio.Reader) {
	fileName, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("can't read file name %e", err)
		fmt.Println("Не удалось прочесть название файла")
		return
	}
	if fileName == "" {
		return
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("can't read bytes %e", err)
		fmt.Println("Не удалось прочесть файл")
		return
	}
	fileName = strings.TrimSpace(fileName)
	file, err := os.Create("downloads/" + fileName)

	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("can't close file %v", err)
		}
	}()
	_, err = file.Write(bytes)
	if err != nil {
		log.Printf("can't write to file %v", err)
		fmt.Println("Не удалось загрузить файл файл")
	}
	fmt.Println("done")
}