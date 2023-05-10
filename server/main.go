package main

import (
	"context"
	pb "github.com/go-grpc-assignment/protos"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"strconv"
)

const (
	port = ":8090"
)

var products []*pb.ProductInfo

type productServer struct {
	pb.UnimplementedProductServer
}

func main() {
	initProducts()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterProductServer(s, &productServer{})

	log.Printf("Server listening at %v", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func initProducts() {
	var productItems = []pb.ProductInfo{
		{Id: "1", Name: "Samsung", Description: "8GB", Price: 45000},
		{Id: "2", Name: "iPhone", Description: "6GB", Price: 75000},
	}
	products = append(products, &productItems[0])
	products = append(products, &productItems[1])

}

func (s *productServer) GetProducts(in *pb.Empty, stream pb.Product_GetProductsServer) error {
	log.Printf("Received : %v", in)
	for _, product := range products {
		if err := stream.Send(product); err != nil {
			return err
		}
	}
	return nil
}

func (s *productServer) GetProduct(ctx context.Context, in *pb.Id) (*pb.ProductInfo, error) {
	log.Printf("Received: %v", in)

	res := &pb.ProductInfo{}

	for _, product := range products {
		if product.Id == in.GetValue() {
			res = product
			break
		}
	}
	return res, nil
}

func (s *productServer) CreateProduct(ctx context.Context, in *pb.ProductInfo) (*pb.Id, error) {
	log.Printf("Received: %v", in)
	res := pb.Id{}
	res.Value = strconv.Itoa(rand.Intn(10000000))
	in.Id = res.GetValue()
	products = append(products, in)
	return &res, nil
}

func (s *productServer) UpdateProduct(ctx context.Context, in *pb.ProductInfo) (*pb.Status, error) {
	log.Printf("Received: %v", in)

	res := pb.Status{}
	for index, product := range products {
		if product.GetId() == in.GetId() {
			products = append(products[:index], products[index+1:]...)
			in.Id = product.GetId()
			products = append(products, in)
			res.Value = 1
			break
		}
	}
	return &res, nil
}

func (s *productServer) DeleteProduct(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	log.Printf("Received : %v", in)

	res := pb.Status{}
	for index, product := range products {
		if product.GetId() == in.GetValue() {
			products = append(products[:index], products[index+1:]...)
			res.Value = 1
			break
		}
	}
	return &res, nil

}
