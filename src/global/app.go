package global

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type app struct {
	BuildTime time.Time
	Version   string

	Host string
	Port string

	SocketHost string
	SocketPort string

	ProjectRoot string
	TemplateDir string

	locker sync.Mutex

	Con net.Conn
}

var App = &app{}
var buffServer = make([]byte, 1024)
var buffClient = make([]byte, 1024)
var clients = make(map[string]net.Conn)
var messages = make(chan string, 10)

// InitPath 初始化相关路径，包括项目根目录、模板目录
func (this *app) InitPath() {
	App.setProjectRoot()

	App.SetTemplateDir("default")
}

func (this *app) setProjectRoot() {
	curFilename := os.Args[0]

	binaryPath, err := exec.LookPath(curFilename)
	if err != nil {
		panic(err)
	}

	log.Printf("%v", binaryPath)

	binaryPath, err = filepath.Abs(binaryPath)

	if err != nil {
		panic(err)
	}

	basePath := filepath.Base(binaryPath)

	log.Printf("%v", basePath)

	projectRoot := filepath.Dir(filepath.Dir(binaryPath))

	log.Printf("%v", projectRoot)

	this.ProjectRoot = projectRoot + "/"
}

func (this *app) SetTemplateDir(theme string) {
	this.TemplateDir = this.ProjectRoot + "template/theme/" + theme + "/"
}

func (this *app) InitSocket(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			panic(err)
		}

		log.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		log.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			log.Println("Can't send")
			break
		}
	}

}

func (this *app) InitTcp() {
	port := this.SocketHost + ":" + this.SocketPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)

	if err != nil {
		log.Printf("error_2: %s\n", err)
	}

	l, err := net.ListenTCP("tcp", tcpAddr)

	defer l.Close()
	if err != nil {
		log.Printf("error_3: %s\n", err)
	}

	go BoradCast(clients, messages)

	for {
		con, err := l.Accept()
		if err != nil {
			log.Printf("error_1: %s\n", err)
		}
		clients[con.RemoteAddr().String()] = con
		go HandleServer(con, messages)
	}
}

func (this *app) LoginTcpServer(conType int) {
	localhost := this.Host + ":" + this.Port

	tcpAddr, err := net.ResolveTCPAddr("tcp4", localhost)

	if err != nil {
		log.Printf("connect tcp error:%v", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		defer conn.Close()
		log.Printf("connect tcp error:%v", err)
	}

	go TcpClientSend(conn)

	if conType == 1 {
		log.Printf("%v", 2)
		_, err := conn.Write([]byte(conn.RemoteAddr().String()))
		if err != nil {
			log.Printf("send message error: %v", err)
			return
		}
	} else {
		for {
			log.Printf("%v", localhost)
			_, err := conn.Read(buffClient)
			if err != nil {
				log.Printf("read message error:%v", err)
				return
			}
			log.Printf("收到消息:%s\n", string(buffClient))
		}
	}
}

func TcpClientSend(conn net.Conn) {
	var input string
	username := conn.RemoteAddr().String()
	for {
		fmt.Scanln(&input)

		if input == "/quit" {
			log.Printf("Info: %s%s\n", "ByeBye", username)
			conn.Close()
		}

		_, err := conn.Write([]byte(username + ":" + input))

		if err != nil {
			log.Printf("send message error: %v", err)
			return
		}
	}
}

/**
*
*广播信息
 */
func BoradCast(clients map[string]net.Conn, messages chan string) {
	for {
		msg := <-messages
		for index, client := range clients {
			_, err := client.Write([]byte(msg))
			if err != nil {
				log.Printf("error_a: %v\n", err)
				delete(clients, index)
			}
		}
	}
}

/**
*接收信息
 */
func HandleServer(conn net.Conn, messages chan string) {
	for {
		_, err := conn.Read(buffServer)
		if err != nil {
			log.Printf("error: %s\n", err)
			conn.Close()
			return
		}
		log.Printf("收到消息:%s\n", string(buffServer))
		messages <- string(buffServer)
	}
}
