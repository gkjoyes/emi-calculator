package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gkjoyes/emi-calculator/emicalculatorpb"
	"google.golang.org/grpc"
)

var client emicalculatorpb.EMICalculatorServiceClient

// Establish a connection with the gRPC server first.
func init() {
	conn, err := grpc.Dial("server:5300", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to establish a connection with grpc server: %v\n", err)
	}

	client = emicalculatorpb.NewEMICalculatorServiceClient(conn)
}

// EMIRequest request includes all required parameters to calculate monthly EMI.
type EMIRequest struct {
	DownPayment          float64 `json:"down_payment"`
	TotalAmount          float64 `json:"total_amount"`
	InterestRate         float32 `json:"interest_rate"`
	PropertyTax          float32 `json:"property_tax"`
	PropertyTransferTaxe float32 `json:"property_transfer_tax"`
	YearsExpectedToLive  int32   `json:"years_expected_to_live"`
}

// Response represent final response format for emi-calculator.
type Response struct {
	Status string         `json:"status"`
	Code   int64          `json:"-"`
	Error  string         `json:"error,omitempty"`
	Result *FinalResponse `json:"result,omitempty"`
}

// FinalResponse represents the final result.
type FinalResponse struct {
	MonthlyEMI float32 `json:"monthly_emi,omitempty"`
}

// fail sent failed response to the client.
func fail(w http.ResponseWriter, res *Response) {

	j, err := json.Marshal(res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	w.WriteHeader(int(res.Code))
	w.Write(j)
}

// succeed sent an ok response to the client.
func succeed(w http.ResponseWriter, res *Response) {

	j, err := json.Marshal(res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	w.WriteHeader(int(res.Code))
	w.Write(j)
}

// calculateEMI will call gRPC service for calculating monthly EMI.
func calculateEMI(req *EMIRequest) (float32, error) {

	request := emicalculatorpb.EMICalculatorRequest{
		DownPayment:           req.DownPayment,
		TotalAmount:           req.TotalAmount,
		InterestRate:          req.InterestRate,
		PropertyTaxes:         req.PropertyTax,
		PropertyTransferTaxes: req.PropertyTransferTaxe,
		YearsExpectedToLive:   req.YearsExpectedToLive,
	}

	// Call CalculateEMI gRPC service for getting monthly EMI.
	res, err := client.CalculateEMI(context.Background(), &request)
	if err != nil {
		return 0, fmt.Errorf("error while calling gRPC CalculateEMI service: %v", err)
	}

	return res.MonthlyEmi, nil
}

func validateRequest(req *EMIRequest) error {
	if req.TotalAmount == 0 {
		return fmt.Errorf("the total amount cannot be zero")
	}

	if req.YearsExpectedToLive <= 0 {
		return fmt.Errorf("please enter a valid year number greater than zero")
	}

	if req.InterestRate < 0 {
		return fmt.Errorf("please provide valid interest rate >= 0")
	}

	if req.PropertyTransferTaxe < 0 {
		return fmt.Errorf("please provide valid property transfer tax >= 0")
	}

	if req.DownPayment > req.TotalAmount {
		return fmt.Errorf("downpayment cannot greater than the total amount")
	}
	return nil
}

func emiCalculator(w http.ResponseWriter, r *http.Request) {

	// Currently, we cannot allow other methods except post.
	if r.Method != http.MethodPost {
		fail(w, &Response{Status: "nok", Code: http.StatusNotFound, Error: "not found"})
		return
	}

	// Read request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fail(w, &Response{Status: "nok", Code: http.StatusBadRequest, Error: "invalid request body"})
		return
	}
	defer r.Body.Close()

	var req EMIRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		fail(w, &Response{Status: "nok", Code: http.StatusBadRequest, Error: err.Error()})
		return
	}

	// Check all request parameters are valid.
	if err := validateRequest(&req); err != nil {
		fail(w, &Response{Status: "nok", Code: http.StatusBadRequest, Error: err.Error()})
		return
	}

	// Calculate monthly EMI using gRPC service.
	emi, err := calculateEMI(&req)
	if err != nil {
		fail(w, &Response{Status: "nok", Code: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	succeed(w, &Response{Status: "ok", Code: http.StatusOK, Result: &FinalResponse{MonthlyEMI: emi}})
}

func handleRequests() {
	http.HandleFunc("/emi-calculator", emiCalculator)
	log.Fatal(http.ListenAndServe(":5200", nil))
}

func main() {
	handleRequests()
}
