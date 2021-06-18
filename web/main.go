package main

import (
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-plugins/registry/etcdv3"
	"net/http"

	"a.a/mp-server/web/handler"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	etcReg := etcdv3.NewRegistry()
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.web"),
		web.Version("latest"),
		web.Registry(etcReg),
		web.Address("127.0.0.1:11225"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	service.HandleFunc("/web/call", handler.WebCall)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
