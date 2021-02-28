package main

import (
	"context"
	"log"
	"math"
	"net"

	"github.com/gkjoyes/emi-calculator/emicalculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) CalculateEMI(ctx context.Context, req *emicalculatorpb.EMICalculatorRequest) (*emicalculatorpb.EMICalculatorResponse, error) {
	rate := req.InterestRate / (12 * 100)                                                 // one month interest.
	time := req.YearsExpectedToLive * 12                                                  // one month period.
	principal := (req.TotalAmount - req.DownPayment) + float64(req.PropertyTransferTaxes) // loan amount i.e principal amount
	emi := (principal * float64(rate) * math.Pow(float64(rate+1), float64(time))) / (math.Pow(float64(rate+1), float64(time)) - 1)

	// Final response with monthly EMI value.
	response := emicalculatorpb.EMICalculatorResponse{
		MonthlyEmi: float32(emi) + (req.PropertyTaxes / 12), // Add yearly property tax with monthly EMI.
	}

	return &response, nil
}

func main() {

	// Listen on port 50051(default gRPC port).
	listen, err := net.Listen("tcp", "0.0.0.0:5300")
	if err != nil {
		log.Fatalf("failed to listen port 5300: %v\n", err)
	}

	// Create a new gRPC server and serve the requests.
	grpcServer := grpc.NewServer()
	emicalculatorpb.RegisterEMICalculatorServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
