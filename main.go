package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"poc_hybrid_grpc_rest/controller"
	"poc_hybrid_grpc_rest/pb"
	"poc_hybrid_grpc_rest/service"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := listenRest(wg)
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Add(1)
	go func() {
		err := listenGrpc(wg)
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
	os.Exit(0)
}

func listenRest(wg *sync.WaitGroup) error {
	defer wg.Done()
	userService := service.NewUserService()
	userController := controller.NewUserRestController(userService)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/user" {
			controller.RestAdapterFunc(userController.CreateUser)(w, r)
		} else if r.Method == http.MethodPost && r.URL.Path == "/validate" {
			controller.RestAdapterFunc(userController.Validate)(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("{\"customerror\":\"not found\"}"))
		}
	})
	return http.ListenAndServe(":8000", handler)
}

func listenGrpc(wg *sync.WaitGroup) error {
	defer wg.Done()
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		return err
	}
	userService := service.NewUserService()
	grpcUserController := controller.NewUserGrpcController(userService)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, grpcUserController)
	return s.Serve(listener)
}
