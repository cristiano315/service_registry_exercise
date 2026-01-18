package main

import (
	"fmt"
	"net/rpc"
	"sync/atomic"
	"time"
)

type LoadBalancer struct { counter uint64 }
type Services struct {
	Names []string
	Calls []string
}

func (lb *LoadBalancer) Next(addrs []string) string {
	idx := atomic.AddUint64(&lb.counter, 1) -1
	return addrs[idx % uint64(len(addrs))]
}

var port = ":1234"
var services = Services{ Names: []string{ "weather", "counter" }, Calls: []string{ "WeatherService.GetWeather", "CounterService.Increment" } }

func main() {
	//Create Load Balancers
	var lbs []LoadBalancer
	for i:=0 ; i < len(services.Names); i++ {
		lbs = append(lbs, LoadBalancer{})
	}

	addresses := make(map[string][]string)

	//Get addresses for all services every minute
	go func(){
		//RPC discovery
		regClient, err := rpc.Dial("tcp", "registry" + port)
		for err != nil {
			time.Sleep(1 * time.Second)
			fmt.Println("[Client] Waiting for registry to be available...")
		}
		for i := 0; i < len(services.Names); i++{
			var res struct { Addrs []string  }
			regClient.Call("Registry.Discover", services.Names[i], &res)
			addresses[services.Names[i]] = res.Addrs
		}
		regClient.Close()
		time.Sleep(1 * time.Minute)
	}()
	time.Sleep(2 * time.Second) //Wait for initial discovery

	for{
		//Call each service
		for i, svcName := range services.Names{
			if(len(addresses[svcName]) > 0){
			//Load balancing
			selectedAddr := lbs[i].Next(addresses[svcName])

			//Service call
			svcClient, err := rpc.Dial("tcp", selectedAddr)
			if(err == nil){
					var reply string
					svcClient.Call(services.Calls[i], "", &reply)
					fmt.Printf("[Client] Response: %s\n", reply)
					svcClient.Close()
				}
			} else {
				fmt.Printf("[Client] No instances available for service: %s\n", svcName)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
