package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/jaismanish15/CurrencyConverter/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCurrencyConverterServiceServer
}

type Currency string

const (
	INR Currency = "INR"
	USD Currency = "USD"
)

var CurrencyMap = map[string]Currency{
	"INR": INR,
	"USD": USD,
}

var conversionRate = map[Currency]float32{
	INR: 1,
	USD: 83,
}

func (c Currency) convertToBase(amount float32) float32 {
	return amount * conversionRate[c]
}

func (c Currency) convertFromBase(amount float32) float32 {
	return amount / conversionRate[c]
}

func (s *server) Convert(ctx context.Context, req *pb.CurrencyConverterRequest) (*pb.CurrencyConverterResponse, error) {
	var fromCurrency Currency
	if currency, ok := CurrencyMap[req.InitialCurrency]; ok {
		fromCurrency = currency
	} else {
		return &pb.CurrencyConverterResponse{}, fmt.Errorf("unsupported currency: %s", req.InitialCurrency)
	}

	var toCurrency Currency
	if currency, ok := CurrencyMap[req.FinalCurrency]; ok {
		toCurrency = currency
	} else {
		return &pb.CurrencyConverterResponse{}, fmt.Errorf("unsupported currency: %s", req.FinalCurrency)
	}

	fmt.Println(fromCurrency.convertToBase(req.Amount))
	fmt.Println(toCurrency.convertFromBase(fromCurrency.convertToBase(req.Amount)))
	convertedAmount := toCurrency.convertFromBase(fromCurrency.convertToBase(req.Amount))
	return &pb.CurrencyConverterResponse{Amount: convertedAmount}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterCurrencyConverterServiceServer(srv, &server{})

	log.Println("Starting gRPC server on port 50051...")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
