package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	net.Conn
	r *bufio.Reader
}

func (c *tcpClient) sendGet(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s", klen, key)))
}

func (c *tcpClient) sendSet(key, value string) {
	klen := len(key)
	vlen := len(value)
	c.Write([]byte(fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)))
}

func (c *tcpClient) sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

//readLen 以空格为分隔符读取一个字符串并将之转换为一个整形
func readLen(r *bufio.Reader) int {
	tem, err := r.ReadString(' ')
	if err != nil {
		log.Println(err)
		return 0
	}
	l, err := strconv.Atoi(strings.TrimSpace(tem))
	if err != nil {
		log.Println(err)
		return 0
	}
	return l
}

func (c *tcpClient) recvResponse() (string, error) {
	vlen := readLen(c.r)
	if vlen == 0 {
		return "", nil
	}
	if vlen < 0 {
		err := make([]byte, -vlen)
		_, e := io.ReadFull(c.r, err)
		if e != nil {
			return "", e
		}
		return "", errors.New(string(err))
	}
	value := make([]byte, vlen)
	_, e := io.ReadFull(c.r, value)
	if e != nil {
		return "", e
	}
	return string(value), nil
}

func (c *tcpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		c.sendGet(cmd.Key)
		cmd.Value, cmd.Error = c.recvResponse()
		return
	} else if cmd.Name == "set" {
		c.sendSet(cmd.Key, cmd.Value)
		_, cmd.Error = c.recvResponse()
		return
	} else if cmd.Name == "del" {
		c.sendDel(cmd.Key)
		_, cmd.Error = c.recvResponse()
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

func newTCPClient(server string) *tcpClient {
	c, e := net.Dial("tcp", server+":12346")
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)
	return &tcpClient{c, r}
}
