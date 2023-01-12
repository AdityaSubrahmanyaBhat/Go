package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-ping/ping"
)

type Server struct {
	address string
	proxy   httputil.ReverseProxy
}

type loadBalancer struct{
	port string
	servers []Server
	roundRobinCount int
}

func newLoadBalancer(port string, servers []Server) loadBalancer{
	return loadBalancer{
		port: port,
		servers: servers,
		roundRobinCount: 0,
	}
}

func newServer(address string) Server{
	simpleUrl, err := url.Parse(address)
	if err != nil{
		log.Fatal("error : %v",err)
	}
	return Server{
		address: address,
		proxy: *httputil.NewSingleHostReverseProxy(simpleUrl),
	}
}

func(server Server) isAlive() bool{

	pinger, _ := ping.NewPinger(server.address)
	pinger.Count = 3
	err := pinger.Run()
	fmt.Println(pinger.Statistics())
	if err != nil{
		fmt.Println(pinger.Statistics())
		return false
	}else{
		fmt.Println(pinger.Statistics())
		return true
	}
}

func(server Server) Address() string{
	return server.address
}

func(server Server) Serve(res http.ResponseWriter, req *http.Request){
	server.proxy.ServeHTTP(res,req)
}

func(lb *loadBalancer) getNextAvailableServer() Server{
	server := lb.servers[lb.roundRobinCount % len(lb.servers)]
	if !server.isAlive(){
		lb.roundRobinCount++
		fmt.Println(lb.roundRobinCount)
		server = lb.servers[lb.roundRobinCount % len(lb.servers)]
	}
	lb.roundRobinCount++ 
	return server
}

func(lb loadBalancer) serveProxy(res http.ResponseWriter, req *http.Request){
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding the request to : %v",targetServer.Address())
	targetServer.Serve(res, req)
}

func main() {
	PORT := ":8000"
	servers := []Server{
		newServer("http://www.duckduckgo.com"),
		newServer("http://www.bing.com"),
		newServer("http://www.facebook.com"),
	}
	lb := newLoadBalancer(PORT,servers)
	fmt.Printf("Server started at port %v\n",PORT)
	http.HandleFunc("/",func(res http.ResponseWriter, req *http.Request) {
		lb.serveProxy(res,req)
	})
	http.ListenAndServe(lb.port,nil)
}