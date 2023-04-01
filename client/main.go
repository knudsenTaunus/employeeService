package main

import (
	"context"
	"fmt"
	userpb "github.com/knudsenTaunus/employeeService/gen/go/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("localhost:9879", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	clientName := fmt.Sprintf("foo%s", time.Now().String())

	go subscribe(client, clientName)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	resp, err := client.Deregister(context.Background(), &userpb.UserServiceDeregisterRequest{ClientName: clientName})
	log.Printf("deregistered client %s", resp.ClientName)
}

func subscribe(client userpb.UserServiceClient, clientName string) {
	log.Print("subscribing to updates")
	registerResp, err := client.Register(context.Background(), &userpb.UserServiceRegisterRequest{ClientName: clientName})
	if err != nil {
		log.Fatalln("Couldn't request", err)
	}

	for {
		updateResp, err := registerResp.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatalln("Receiving", err)
			return
		}

		log.Print(updateResp.FirstName)
	}

}
