package tcp

import (
	"bufio"
	"fmt"
	"github.com/donghc/dcache/pkg/cache"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	listen, err := net.Listen("tcp", ":12346")
	if err != nil {
		panic(err)
	}
	for {
		accept, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go s.process(accept)
	}
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) process(conn net.Conn) {
	//延时关闭conn
	defer conn.Close()
	//对客户端链接进行一个缓冲读取。因为来自网络的数据不稳定，
	//在我们读取时，客户端的数据可能只传输了一半，我们希望可以阻塞等待，
	//直到需要的数据全部就位以后一次性返回给我们
	r := bufio.NewReader(conn)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error :", e)
				return
			}
			if op == 'S' {
				e = s.set(conn, r)
			} else if op == 'G' {
				e = s.get(conn, r)
			} else if op == 'D' {
				e = s.del(conn, r)
			} else {
				log.Println("close connection due to invalid operation :", op)
				return
			}
			if e != nil {
				log.Println("close connection due to error :", e)
				return
			}
		}
	}

}

//readKey 解析客户端发送过来的command，从中获取key
func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, err := readLen(r)
	if err != nil {
		return "", err
	}
	k := make([]byte, klen)
	_, err = io.ReadFull(r, k)
	if err != nil {
		return "", err
	}
	return string(k), nil
}

//readKeyAndValue 解析客户端发送过来的command，从中获取key和value
func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	klen, err := readLen(r)
	if err != nil {
		return "", nil, err
	}
	vlen, err := readLen(r)
	if err != nil {
		return "", nil, err
	}
	k := make([]byte, klen)
	_, err = io.ReadFull(r, k)
	if err != nil {
		return "", nil, err
	}
	v := make([]byte, vlen)
	_, err = io.ReadFull(r, v)
	if err != nil {
		return "", nil, err
	}
	return string(k), v, nil
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}
	v, e := s.Get(k)
	return sendResponse(v, e, conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	k, v, e := s.readKeyAndValue(r)
	if e != nil {
		return e
	}
	return sendResponse(nil, s.Set(k, v), conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	k, err := s.readKey(r)
	if err != nil {
		return err
	}
	return sendResponse(nil, s.Del(k), conn)
}

//sendResponse 根据参数将服务端的error或者value写入客户端链接
func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		temp := fmt.Sprintf("-%d", len(errString)) + errString
		_, e := conn.Write([]byte(temp))
		return e
	}
	vlen := fmt.Sprintf("%d", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	return e
}

//readLen 以空格为分隔符读取一个字符串并将之转换为一个整形
func readLen(r *bufio.Reader) (int, error) {
	tem, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}
	l, err := strconv.Atoi(strings.TrimSpace(tem))
	if err != nil {
		return 0, err
	}
	return l, nil
}
