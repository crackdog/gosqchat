//ts3sq provides a library for the ts3 server query interface.
package ts3sq

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Ts3sqs struct {
	serverconn net.Conn
}

func New(address string) (*Ts3sqs, error) {
	c, err := net.Dial("tcp", address)
	if err == nil {
		t := new(Ts3sqs)
		t.serverconn = c
		return t, nil
	} else {
		return nil, err
	}
}

func (s *Ts3sqs) send(msg string) error {
	length, err := s.serverconn.Write([]byte(msg))
	if err == nil && length < len(msg) {
		return fmt.Errorf("only %d of %d bytes were sended.", length, len(msg))
	} else {
		return err
	}
}

func escape(s string) string {
	//escape replaces all spaces with \s
	return s
}

func (s *Ts3sqs) WaitForMessageLine() (string, error) {
	//WaitForMessageLine reads a line from the server connection.
	return bufio.NewReader(s.serverconn).ReadString('\n')
}

func (s *Ts3sqs) getError() error {
	//getError returns the error message from the server.
	msg, err := s.WaitForMessageLine()
	if err != nil {
		return err
	} else {
		if strings.Contains(msg, "error id=0 msg=ok") {
			return nil
		} else {
			return fmt.Errorf("msg error: '%s'", strings.TrimSpace(msg))
		}
	}
}

func (s *Ts3sqs) sendWithGettingError(msg string) error {
	err := s.send(msg)
	if err != nil {
		return err
	} else {
		return s.getError()
	}
}

func (s *Ts3sqs) Login(username, password string) error {
	//logging in...
	username = escape(username)
	password = escape(password)
	msg := fmt.Sprintf("login client_login_name=%s client_login_password=%s\n")
	return s.sendWithGettingError(msg)
}

func (s *Ts3sqs) Logout() error {
	//logging out
	return s.sendWithGettingError("logout\n")
}

func (s *Ts3sqs) Clientlist() (string, error) {
	//Clientlist sends a clientlist request to the ts3 server.
	panic("Clientlist() is not implemented yet")
}

func (s *Ts3sqs) Use(server_id int) error {
	//Use sends a request to use a server.
	msg := fmt.Sprintf("use sid=%d", server_id)
	return s.sendWithGettingError(msg)
}

func (s *Ts3sqs) Servernotifyregister(event string) error {
	//Servernotifyregister sends a notify request for a given event.
	msg := fmt.Sprintf("servernotifyregister event=%s", escape(event))
	return s.sendWithGettingError(msg)
}

func (s *Ts3sqs) Servernotifyunregister(event string) error {
	//Servernotifyunregister sends a unnotify request for a given event.
	msg := fmt.Sprintf("servernotitfyunregister event=%s", escape(event))
	return s.sendWithGettingError(msg)
}

func (s *Ts3sqs) Sendtextmessage(targetmode, target int, raw_msg string) error {
	msg := fmt.Sprintf("sendtextmessage targetmode=%d target=%d msg=%s",
		targetmode, target, escape(raw_msg))
	return s.sendWithGettingError(msg)
}
