package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/knudsenTaunus/employeeService/generated/go/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		if err != nil {
			if status.Code(err) == codes.Unavailable {
				log.Fatalln("server is not available anymore", err)
				return
			}

			log.Fatalln("received error: ", err)
			return
		}

		log.Print(updateResp.String())
	}

}
