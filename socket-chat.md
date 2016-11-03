### server.go
```go
package main

import (
	"fmt"
	"net"
	"os"
)

var cliQue []net.Conn

func main() {
	var (
		host   = "127.0.0.1"
		port   = "8000"
		remote = host + ":" + port
		data   = make([]byte, 1024)
	)
	fmt.Println("Initiating server... (Ctrl-C to stop)")

	lis, err := net.Listen("tcp", remote)
	defer lis.Close()

	if err != nil {
		fmt.Printf("Error when listen: %s, Err: %s\n", remote, err)
		os.Exit(-1)
	}

	for {
		var res string
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("Error accepting client: ", err.Error())
			os.Exit(0)
		}
		cliQue = append(cliQue, conn)

		go func(con net.Conn) {
			fmt.Println("New connection: ", con.RemoteAddr())

			// Get client's name
			length, err := con.Read(data)
			if err != nil {
				fmt.Printf("Client %v quit.\n", con.RemoteAddr())
				con.Close()
				disconnect(con, con.RemoteAddr().String())
				return
			}
			name := string(data[:length])
			comeStr := name + " entered the room."
			notify(con, comeStr)

			// Begin recieve message from client
			for {
				length, err := con.Read(data)
				if err != nil {
					fmt.Printf("Client %s quit.\n", name)
					con.Close()
					disconnect(con, name)
					return
				}
				res = string(data[:length])
				sprdMsg := name + " said: " + res
				fmt.Println(sprdMsg)
				res = "You said:" + res
				con.Write([]byte(res))
				notify(con, sprdMsg)
			}
		}(conn)
	}
}

// Notify other clients
func notify(conn net.Conn, msg string) {
	for _, con := range cliQue {
		if con.RemoteAddr() != conn.RemoteAddr() {
			con.Write([]byte(msg))
		}
	}
}

// Delete handle of the disconnected client and notify others
func disconnect(conn net.Conn, name string) {
	for index, con := range cliQue {
		if con.RemoteAddr() == conn.RemoteAddr() {
			disMsg := name + " has left the room."
			fmt.Println(disMsg)
			cliQue = append(cliQue[:index], cliQue[index+1:]...)
			notify(conn, disMsg)
		}
	}
}
```


### client.go

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var writeStr, readStr = make([]byte, 1024), make([]byte, 1024)

func main() {
	var (
		host   = "127.0.0.1"
		port   = "8000"
		remote = host + ":" + port
		reader = bufio.NewReader(os.Stdin)
	)

	con, err := net.Dial("tcp", remote)
	defer con.Close()

	if err != nil {
		fmt.Println("Server not found.")
		os.Exit(-1)
	}
	fmt.Println("Connection OK.")
	fmt.Printf("Enter your name: ")
	fmt.Scanf("%s", &writeStr)
	in, err := con.Write([]byte(writeStr))
	if err != nil {
		fmt.Printf("Error when send to server: %d\n", in)
		os.Exit(0)
	}
	fmt.Println("Now begin to talk!")
	go read(con)

	for {
		writeStr, _, _ = reader.ReadLine()
		if string(writeStr) == "quit" {
			fmt.Println("Communication terminated.")
			os.Exit(1)
		}

		in, err := con.Write([]byte(writeStr))
		if err != nil {
			fmt.Printf("Error when send to server: %d\n", in)
			os.Exit(0)
		}

	}
}

func read(conn net.Conn) {
	for {
		length, err := conn.Read(readStr)
		if err != nil {
			fmt.Printf("Error when read from server. Error:%s\n", err)
			os.Exit(0)
		}
		fmt.Println(string(readStr[:length]))
	}
}
```
