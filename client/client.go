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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/tenorbear/grpc-go-loadtest/helloworld"
	"google.golang.org/grpc"
)

type ArgumentError string

func (a ArgumentError) Error() string {
	return fmt.Sprintf("Argument error: %s", string(a))
}

var mode = flag.String("mode", "normal", "normal|latency|error, to switch the type of message.")
var address = flag.String("address", "localhost:50051", "Server to connect to.")
var name = flag.String("name", "world", "Your name.")

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var r *pb.HelloReply
	switch *mode {
	case "normal":
		r, err = c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	case "latency":
		r, err = c.SayHelloWithLatency(ctx, &pb.HelloRequest{Name: *name})
	case "error":
		r, err = c.SayHelloWithError(ctx, &pb.HelloRequest{Name: *name})
	default:
		err = ArgumentError(*mode)
	}
	if err != nil {
		log.Fatalf("Error in making RPC: %s.", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
