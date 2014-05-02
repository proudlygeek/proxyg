package proxyg

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
)

func print(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", fmt.Sprintf(s, a...))
}

func newconn(conn net.Conn) {
	connectRegexp := regexp.MustCompile("CONNECT (.*) HTTP/")
	reader := bufio.NewReader(conn)

	l, isPrefix, err := reader.ReadLine()

	if err != nil || isPrefix == true {
		print("Error: %v", err)
		conn.Close()
		return
	}

	dst := connectRegexp.FindStringSubmatch(string(l))

	// print("Buffer (content length is %d):\n %s\nHosted Match is %s", contentLength, buf[0:contentLength], dst)

	if dst == nil {
		io.WriteString(conn, "HTTP/1.0 502 Bad Gateway\r\n\r\n")
		conn.Close()
		return
	}

	for {
		l, _, _ := reader.ReadLine()

		if l == nil {
			return
		}
		if len(l) == 0 {
			break
		}
	}

	proxyConnect(conn, dst[1])

}

func Copy(a io.ReadWriteCloser, b io.ReadWriteCloser) {
	// Setup one-way forwarding
	io.Copy(a, b)
	a.Close()
	b.Close()
}

func proxyConnect(localConn net.Conn, host string) {
	print("Connecting Host %s...\n", host)

	// localAddress, _ := net.ResolveTCPAddr("", "127.0.0.1:8080")
	var localAddress *net.TCPAddr

	remoteAddress, err := net.ResolveTCPAddr("", host)
	if err != nil {
		fmt.Fprint(localConn, "HTTP/1.0 502 Bad Gateway. Address not resolved, baby.\r\n\r\n")
		print("Error: %v", err)
		localConn.Close()
		return
	}

	print("Remote Address: %v", remoteAddress)

	remoteConn, err := net.DialTCP("tcp", localAddress, remoteAddress)
	if remoteConn == nil {
		fmt.Fprint(localConn, "HTTP/1.0 502 Bad Gateway. Connection not established, honey.\r\n\r\n")
		print("Error: %v", err)
		localConn.Close()
		return
	}

	remoteConn.SetKeepAlive(true)
	fmt.Fprint(localConn, "HTTP/1.0 200 CONNECTION ESTABLISHED\r\n\r\n")

	go Copy(localConn, remoteConn)
	go Copy(remoteConn, localConn)
}

func ConnectListen() {
	host := flag.String("host", "127.0.0.1", "HTTP Proxy Host")
	port := flag.Int("port", 8080, "HTTP Proxy Port")
	flag.Parse()

	portspec := fmt.Sprintf("%s:%d", *host, *port)

	netlisten, err := net.Listen("tcp", portspec)
	if netlisten == nil {
		print("Error: %v", err)
		os.Exit(1)
	}
	defer netlisten.Close()
	fmt.Fprintf(os.Stderr, "Listening for HTTP CONNECT's on %s\n", portspec)
	fmt.Fprint(os.Stderr, "-------------------------------------\n")

	for {
		conn, err := netlisten.Accept()
		if conn != nil {
			go newconn(conn)
		} else {
			print("Error: %v", err)
		}
	}
}
