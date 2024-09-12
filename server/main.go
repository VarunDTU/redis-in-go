package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/VarunDTU/redis-in-go/resp"
)
const (HOST="localhost";PORT="5001";TYPE="tcp")
func main() {
    listen,err:=net.Listen(TYPE,HOST+":"+PORT);
	if(err!=nil){
		log.Fatal(err);
		os.Exit(0);
	}
	defer listen.Close();
	for {
		conn,err:=listen.Accept()
		if(err!=nil){
			log.Fatal(err);
			os.Exit(0);
		}
		go handleRequest(conn);
	}

}

func handleRequest(conn net.Conn){
	
}

func read(conn net.Conn){
	buffer:=make([]byte,1024);
	_,err:=conn.Read(buffer);
	if(err!=nil){
		log.Fatal(err);
		os.Exit(0);
	}
	time :=time.Now().Format(time.ANSIC);
	responeStr:=fmt.Sprintf("time: %v ,message: %v",time,string(buffer[:]));

	responseResp,_:=resp.RespFromBytes([]byte(responeStr));




	conn.Write([]byte(responseResp));

	defer conn.Close();
}

