package main

import (
	"context"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
	"github.com/redis/go-redis/v9"
)

type CounterService struct{
	ID string
	Value int
	RDB *redis.Client
}
type Args struct { Name, Addr string }

var servicePort = ":1236"
var registryPort = ":1234"
var redisPort = ":6379"
var ctx = context.Background()

func (s *CounterService) Increment(args string, reply *string) error {
	//Increment counter in Redis
	val, err := s.RDB.Incr(ctx, "shared_counter").Result()
	if err != nil {
		return err
	}

	//Set reply value
	*reply = "Counter value: " + strconv.Itoa(int(val)) + " (response from " + s.ID + ")"
	return nil
}

func main() {
	hostname, _ := os.Hostname()
	addr := hostname + servicePort
	//Setup RPC server and redis client
	counterSvc := &CounterService{ID: hostname}
	counterSvc.RDB = redis.NewClient(&redis.Options{Addr: "redis" + redisPort})
	defer counterSvc.RDB.Close()
	rpc.Register(counterSvc)
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
	client.Call("Registry.Register", &Args{Name: "counter", Addr: addr}, &success)
	if(success){
		log.Printf("Service %s registered with address %s and ready.\n", hostname, addr)
	}
	
	//Quit after 5 minutes
	<-time.After(5 * time.Minute)

	log.Println("Closing service")
	client.Call("Registry.Deregister", &Args{Name: "counter", Addr: addr}, &success)
	client.Close()
	listener.Close()
}