package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)
const (
	HOST = "Localhost"
	PORT = "5001"
	TYPE = "tcp"
)

func main() {
	if(len(os.Args)!=2){
		fmt.Println("require one argument");
		os.Exit(1);
	}
	data:=[]byte(getClrf(os.Args[1]));
	fmt.Println(getClrf(os.Args[1]))
	tcpserver,err:=net.ResolveTCPAddr(TYPE,HOST+":"+PORT);
	if(err!=nil){log.Fatal("ResolveTCPAddr failed",err);os.Exit(1);}
	fmt.Println("client listening at PORT:"+PORT);
	conn, err := net.DialTCP(TYPE,nil,tcpserver);
	if(err!=nil){log.Fatal("Dialup failed",err);os.Exit(1);}
	_,err=conn.Write([]byte(data));
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

func getClrf(input string)(string){
	input=strings.ReplaceAll(input,"\n","\\n")
	input=strings.ReplaceAll(input,"\r","\\r")
	return input;
}
