package util

import (
	"net"
	"fmt"
)

func CreateSocketListen(ipAddr string, port int) (tcpListen *net.TCPListener, err error) {
	serAddr := fmt.Sprintf("%s:%d", ipAddr, port)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serAddr)
	if err != nil {
		return nil, err
	}

	tcpListen, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	return tcpListen, nil
}

func CreateSocketConnect(ipAddr string, port int) (conn net.Conn, err error){
	serHost := fmt.Sprintf("%s:%d", ipAddr, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serHost)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

