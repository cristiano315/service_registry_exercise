package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Registry struct {
	mu sync.RWMutex
	services map[string][]string
}

type Args struct { Name, Addr string }
type Response struct { Addrs []string }

var port = ":1234"

//RPC method to register a service
func (r *Registry) Register(args *Args, reply *bool) error {
	// Use a mutex to ensure thread safety
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[args.Name] = append(r.services[args.Name], args.Addr)
	*reply = true
	return nil
}

//RPC method to discover the services
func (r *Registry) Discover(name string, reply *Response) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if addrs, ok := r.services[name]; ok {
		reply.Addrs = addrs
		return nil
	}
	return errors.New("Service not found")
}

//RPC method to deregister a service
func (r *Registry) Deregister(args *Args, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	addrs, ok := r.services[args.Name]
	if !ok {
		*reply = false
		return errors.New("service name not found")
	}

	// Search for the address to remove
	found := false
	for i, addr := range addrs {
		if addr == args.Addr {
			r.services[args.Name] = append(addrs[:i], addrs[i+1:]...)
			if len(r.services[args.Name]) == 0 {
				delete(r.services, args.Name)
				fmt.Printf("Service %s has no more instances and has been removed from registry.\n", args.Name)
			}
			found = true
			break
		}
	}

	// Se la slice è vuota dopo la rimozione, puliamo la mappa
	if len(r.services[args.Name]) == 0 {
		delete(r.services, args.Name)
	}

	*reply = found
	return nil
}

func main() {
	reg := &Registry{services: make(map[string][]string)}
	rpc.Register(reg)// Register the Registry struct for RPC

	listener, _ := net.Listen("tcp", port)
	log.Println("RPC service registry listening on", port)
	for{
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}