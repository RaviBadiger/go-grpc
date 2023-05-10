package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/go-grpc-assignment/protos"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8090"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewProductClient(conn)
	runGetProducts(client)
	runGetProduct(client, "1")
	runCreateProduct(client, "3","Vivo","4GB",23000)
	runGetProducts(client)
	runUpdateProduct(client,"3","Vivo","8GB",25000)
	runGetProducts(client)
	runDeleteProduct(client,"1")
	runGetProducts(client)
}

func runGetProducts(client pb.ProductClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetProducts(ctx, req)
	if err != nil {
		log.Fatalf(" %v.GetProducts(_) = _, %v ", client, err)
	}
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf(" %v.GetProducts = _ %v", client, err)
		}
		log.Printf("ProductInfo: %v", row)
	}
}

func runGetProduct(client pb.ProductClient, productId string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: productId}
	res, err := client.GetProduct(ctx, req)
	if err != nil {
		log.Fatalf(" %v.GetProducts(_) = _, %v ", client, err)
	}
	log.Printf("ProductInfo: %v", res)
}

func runCreateProduct(client pb.ProductClient, id string, name string, description string, price float32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.ProductInfo{Id: id, Name: name, Description: description, Price: price}
	res, err := client.CreateProduct(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreateProduct = _ %v", res)
	}
	if res.GetValue() != "" {
		log.Printf("CreateProduct ID : %v", res)
	} else {
		log.Printf("CreateProduct Failed")
	}
}

func runUpdateProduct(client pb.ProductClient, id string, name string, description string, price float32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.ProductInfo{Id: id, Name: name, Description: description, Price: price}
	res, err := client.UpdateProduct(ctx, req)
	if err != nil {
		log.Fatalf("%v.UpdateProduct = _ %v", res)
	}
	if int(res.GetValue()) != 1 {
		log.Printf("UpdateProduct Success")
	} else {
		log.Printf("UpdateProduct Failed")
	}
}

func runDeleteProduct(client pb.ProductClient, productId string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: productId}
	res, err := client.DeleteProduct(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeleteProduct = _ %v", res)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("DeleteProduct Success")
	} else {
		log.Printf("DeleteProduct Failed")
	}
}
