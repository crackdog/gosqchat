//Connection represents the connection to the server and
//contains Connection and Encryption objects.
package connection

import (
	//"bufio"
	"fmt"
	"net"
	"sync"
)

/*const (
	Buflen = 4096 //Max length of each message
)*/

type Connection struct {
	//not finished
	conn      net.Conn
	readMutex sync.Mutex
	closed    bool
}

func NewConnection(adress string) *Connection {
	c := new(Connection)
	tmpconn, err := net.Dial("tcp", adress)
	c.closed = false
	if err != nil {
		//handle error
		fmt.Println(err)
		return nil
	} else {
		c.conn = tmpconn
		return c
	}
}

//Read function...
func (c *Connection) Read(b []byte) (n int, err error) {
	c.readMutex.Lock()
	n, err = c.conn.Read(b)
	c.readMutex.Unlock()
	return
}

//Write function...
func (c *Connection) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

//Send sends a string to the server and is blocking.
func (c *Connection) Send(msg string) {
	//fmt.Printf("%s\n", msg)
	fmt.Fprintf(c.conn, "%s\n", msg)
}

//Close closes the connection
func (c *Connection) Close() {
	c.closed = true
	c.Close()
}

//IsClosed returns true if the connection is closed.
func (c *Connection) IsClosed() bool {
	return c.closed
}
