// Usage: gobclient serviceadress
package msgsClient

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"lab2/messages"
	"net"
	"os"
)

var (
	msg    interface{}
	reader *bufio.Reader
)

func MsgsClient(host string) {
	fmt.Println("Connecting to server")
	service := host
	conn, err := net.Dial("tcp", service)
	checkError(err)
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	reader = bufio.NewReader(os.Stdin)
	for msg, err := getinput(); err == nil; msg, err = getinput() {
		fmt.Println("Sending")
		encoder = gob.NewEncoder(conn)
		encoder.Encode(&msg)
	}
	return
}

func getinput() (msg interface{}, err error) {
	var msgtype string
	choice := "Type Msg Enter  or Err Enter to register a normal or an error message. End Enter to finish."

	line := readfrominput(choice)
	fmt.Sscanf(line, "%s", &msgtype)
	switch {
	case msgtype == "Msg":
		var sender, content string
		line := readfrominput("Write: Sender Message")
		fmt.Sscanf(line, "%s %s", &sender, &content)
		msg = messages.StrMsg{sender, content}
		err = nil
		fmt.Println("Message from " + sender + " saying " + content)

	case msgtype == "Err":
		var sender, content string
		line := readfrominput("Write: Sender Error")
		fmt.Sscanf(line, "%s %s", &sender, &content)
		msg = messages.ErrMsg{sender, content}
		err = nil
		fmt.Println("Error from " + sender + " saying " + content)

	case msgtype == "End":
		fmt.Println("Ending input.")
		msg = nil
		err = errors.New("Ending input")
	}
	return msg, err
}

func readfrominput(instruction string) (line string) {
	var err error
	fmt.Println(instruction)
	for line, err = reader.ReadString('\n'); err != nil; line, err = reader.ReadString('\n') {
		fmt.Fprintln(os.Stderr, "invalid input")
		fmt.Println(instruction)
	}
	return line
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error", err.Error())
		os.Exit(1)
	}
}
