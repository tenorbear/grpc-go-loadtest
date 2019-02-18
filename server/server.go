/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/tenorbear/grpc-go-loadtest/helloworld"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "Which port to listen to.")
var hiccupChance = flag.Float64("hiccup_chance", 0.1, "How likely the server gives a slower/erroenous response.")

// server is used to implement helloworld.GreeterServer.
type server struct{}

type RandomHiccupError struct{}

func (err RandomHiccupError) Error() string {
	return "Just a random hiccup."
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloWithLatency(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	// There is a chance for increased latency.
	if chance := rand.Float64(); chance < *hiccupChance {
		time.Sleep(500 * time.Millisecond)
		return &pb.HelloReply{Message: "Sorry for the late hello " + in.Name}, nil
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloWithError(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	// There is a chance for error.
	if chance := rand.Float64(); chance < *hiccupChance {
		return &pb.HelloReply{Message: "No hello for you " + in.Name}, &RandomHiccupError{}
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
