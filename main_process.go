package main

import (
    "fmt"
    "net/http"
)

var urls = []string{
    "http://www.baidu.com/",
    "http://golang.org/",
    "http://blog.golang.org/",
}

func main() {
    // Execute an HTTP HEAD request for all url's
    // and returns the HTTP status string or an error string.
    for _, url := range urls {
        resp, err := http.Head(url)
        if err != nil {
            fmt.Println("Error:", url, err)
        }
        fmt.Print(url, ": ", resp.Status)
    }
}

//Server app
package main

import (
    "fmt"
    "net"
)

func main() {
    fmt.Println("Starting the server ...")
    // 创建 listener
    listener, err := net.Listen("tcp", "localhost:50000")
    if err != nil {
        fmt.Println("Error listening", err.Error())
        return //终止程序
    }
    // 监听并接受来自客户端的连接
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting", err.Error())
            return // 终止程序
        }
        go doServerStuff(conn)
    }
}

func doServerStuff(conn net.Conn) {
    for {
        buf := make([]byte, 512)
        _, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading", err.Error())
            return //终止程序
        }
        fmt.Printf("Received data: %v", string(buf))
    }
}

//Client
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    //打开连接:
    conn, err := net.Dial("tcp", "localhost:50000")
    if err != nil {
        //由于目标计算机积极拒绝而无法创建连接
        fmt.Println("Error dialing", err.Error())
        return // 终止程序
    }

    inputReader := bufio.NewReader(os.Stdin)
    fmt.Println("First, what is your name?")
    clientName, _ := inputReader.ReadString('\n')
    // fmt.Printf("CLIENTNAME %s", clientName)
    trimmedClient := strings.Trim(clientName, "\r\n") // Windows 平台下用 "\r\n"，Linux平台下使用 "\n"
    // 给服务器发送信息直到程序退出：
    for {
        fmt.Println("What to send to the server? Type Q to quit.")
        input, _ := inputReader.ReadString('\n')
        trimmedInput := strings.Trim(input, "\r\n")
        // fmt.Printf("input:--s%--", input)
        // fmt.Printf("trimmedInput:--s%--", trimmedInput)
        if trimmedInput == "Q" {
            return
        }
        _, err = conn.Write([]byte(trimmedClient + " says: " + trimmedInput))
    }
}

//dial
// make a connection with www.example.org:
package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    conn, err := net.Dial("tcp", "192.0.32.10:80") // tcp ipv4
    checkConnection(conn, err)
    conn, err = net.Dial("udp", "192.0.32.10:80") // udp
    checkConnection(conn, err)
    conn, err = net.Dial("tcp", "[2620:0:2d0:200::10]:80") // tcp ipv6
    checkConnection(conn, err)
}
func checkConnection(conn net.Conn, err error) {
    if err != nil {
        fmt.Printf("error %v connecting!")
        os.Exit(1)
    }
    fmt.Println("Connection is made with %v", conn)
}

//socket
package main

import (
    "fmt"
    "io"
    "net"
)

func main() {
    var (
        host          = "www.apache.org"
        port          = "80"
        remote        = host + ":" + port
        msg    string = "GET / \n"
        data          = make([]uint8, 4096)
        read          = true
        count         = 0
    )
    // 创建一个socket
    con, err := net.Dial("tcp", remote)
    // 发送我们的消息，一个http GET请求
    io.WriteString(con, msg)
    // 读取服务器的响应
    for read {
        count, err = con.Read(data)
        read = (err == nil)
        fmt.Printf(string(data[0:count]))
    }
    con.Close()
}

//Simple TCP server
/ Simple multi-thread/multi-core TCP server.
package main

import (
    "flag"
    "fmt"
    "net"
    "os"
)

const maxRead = 25

func main() {
    flag.Parse()
    if flag.NArg() != 2 {
        panic("usage: host port")
    }
    hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))
    listener := initServer(hostAndPort)
    for {
        conn, err := listener.Accept()
        checkError(err, "Accept: ")
        go connectionHandler(conn)
    }
}

func initServer(hostAndPort string) *net.TCPListener {
    serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
    checkError(err, "Resolving address:port failed: '"+hostAndPort+"'")
    listener, err := net.ListenTCP("tcp", serverAddr)
    checkError(err, "ListenTCP: ")
    println("Listening to: ", listener.Addr().String())
    return listener
}

func connectionHandler(conn net.Conn) {
    connFrom := conn.RemoteAddr().String()
    println("Connection from: ", connFrom)
    sayHello(conn)
    for {
        var ibuf []byte = make([]byte, maxRead+1)
        length, err := conn.Read(ibuf[0:maxRead])
        ibuf[maxRead] = 0 // to prevent overflow
        switch err {
        case nil:
            handleMsg(length, err, ibuf)
        case os.EAGAIN: // try again
            continue
        default:
            goto DISCONNECT
        }
    }
DISCONNECT:
    err := conn.Close()
    println("Closed connection: ", connFrom)
    checkError(err, "Close: ")
}

func sayHello(to net.Conn) {
    obuf := []byte{'L', 'e', 't', '\'', 's', ' ', 'G', 'O', '!', '\n'}
    wrote, err := to.Write(obuf)
    checkError(err, "Write: wrote "+string(wrote)+" bytes.")
}

func handleMsg(length int, err error, msg []byte) {
    if length > 0 {
        print("<", length, ":")
        for i := 0; ; i++ {
            if msg[i] == 0 {
                break
            }
            fmt.Printf("%c", msg[i])
        }
        print(">")
    }
}

func checkError(error error, info string) {
    if error != nil {
        panic("ERROR: " + info + " " + error.Error()) // terminate program
    }
}
