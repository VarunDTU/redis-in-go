package main

import (
	"log"
	"net"
	"os"
)
const (
	HOST = "Localhost"
	PORT = "5001"
	TYPE = "tcp"
)

func main() {
	tcpserver,err:=net.ResolveTCPAddr(TYPE,HOST+":"+PORT);
	if(err!=nil){log.Fatal("ResolveTCPAddr failed",err);os.Exit(1);}
	conn, err := net.DialTCP(TYPE,nil,tcpserver);
	if(err!=nil){log.Fatal("Dialup failed",err);os.Exit(1);}
	_,err=conn.Write([]byte("first message by client"));
	if(err!=nil){log.Fatal("write failure",err)}
	recieved:=make([]byte,1024);
	_,err=conn.Read(recieved);
	if(err!=nil){
		log.Fatal("read failure",err);
		os.Exit(1);
	}

	println("Recieved message:",string(recieved));

	conn.Close();








}
