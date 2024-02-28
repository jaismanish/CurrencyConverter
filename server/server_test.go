package main

import (
	"context"
	"net"
	"strings"
	"testing"

	"github.com/jaismanish15/CurrencyConverter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	proto.RegisterCurrencyConverterServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func runTest(t *testing.T, initialCurrency, finalCurrency string, amount, expectedResult float32) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewCurrencyConverterServiceClient(conn)

	req := &proto.CurrencyConverterRequest{
		InitialCurrency: initialCurrency,
		FinalCurrency:   finalCurrency,
		Amount:          amount,
	}

	resp, err := client.Convert(context.Background(), req)
	if err != nil {
		if strings.Contains(err.Error(), "unsupported currency") {
			t.Logf("Unsupported currency encountered: %v", err)
			return
		}
		t.Fatalf("Convert failed: %v", err)
	}

	if resp.Amount != expectedResult {
		t.Errorf("Expected result: %v, got: %v", expectedResult, resp.Amount)
	}
}

// --------------------------------------------------

func TestConvert_USDToINR(t *testing.T) {
	runTest(t, "USD", "INR", 100, 8300)
}

func TestConvert_INRToUSD(t *testing.T) {
	runTest(t, "INR", "USD", 83, 1)
}

func TestConvert_INRToEUR(t *testing.T) {
	runTest(t, "INR", "EUR", 93, 1)
}

func TestConvert_EURToINR(t *testing.T) {
	runTest(t, "EUR", "INR", 1, 93)
}

func TestConvert_SameCurrencyINR(t *testing.T) {
	runTest(t, "INR", "INR", 10, 10)
}

func TestConvert_SameCurrencyUSD(t *testing.T) {
	runTest(t, "USD", "USD", 10, 10)
}

func TestConvert_UnsupportedCurrency(t *testing.T) {
	runTest(t, "XYZ", "INR", 100, 0)
}
