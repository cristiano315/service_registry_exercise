package main

import (
	"net"
	"net/rpc"
	"os"
	"log"
	"time"
)

type WeatherService struct{ ID string }
type Args struct { Name, Addr string }

var servicePort = ":1235"
var registryPort = ":1234"

func (s *WeatherService) GetWeather(args string, reply *string) error {
	*reply = "Sunny (response from " + s.ID + ")"
	return nil
}

func main() {
	hostname, _ := os.Hostname()
	addr := hostname + servicePort

	//Setup RPC server
	weatherSvc := &WeatherService{ID: hostname}
	rpc.Register(weatherSvc)
	listener, _ := net.Listen("tcp", servicePort)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go rpc.ServeConn(conn)
		}
	} ()

	//Register service with registry
	client, _ := rpc.Dial("tcp", "registry" + registryPort)
	var success bool
	client.Call("Registry.Register", &Args{Name: "weather", Addr: addr}, &success)

	if(success){
		log.Printf("Service %s registered with address %s and ready.\n", hostname, addr)
	}
	
	//Quit after 5 minutes
	<-time.After(5 * time.Minute)

	log.Println("Closing service")
	client.Call("Registry.Deregister", &Args{Name: "weather", Addr: addr}, &success)
	client.Close()
	listener.Close()
}