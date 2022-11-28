package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Unix int64  `json:"unix"` // unix time stamp in milliseconds
	Utc  string `json:"utc"`  // UTC value as a string
}

type ErrorResponse struct {
	ErrorValue string `json:"error"`
}

type HealthCheckResponse struct {
	Version string `json:"version"`
	BuildAt string `json:"build-date"`
}

const (
	ResponseLayout = "Mon, 02 Jan 2006 15:04:05 GMT"
	Port           = 8080
)

// will be populated
// through ldflags
// at compile time
var CommitID string
var BuildDate string

func main() {
	err := Run()

	if err != nil {
		log.Fatalf("Exiting due to error:%v\n", err)
		os.Exit(1)
	}

}

func Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", handleHealthCheck).Methods("GET")
	router.HandleFunc("/healthz/", handleHealthCheck).Methods("GET")
	router.HandleFunc("/api", handleEmptyTimeStamp).Methods("GET")
	router.HandleFunc("/api/", handleEmptyTimeStamp).Methods("GET")
	router.HandleFunc("/api/{timestamp}", handleTimeStamp).Methods("GET")
	router.HandleFunc("/api/{timestamp}/", handleTimeStamp).Methods("GET")

	// start the server
	webport := fmt.Sprintf(":%d", Port)
	err := http.ListenAndServe(webport, router)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Server started on port %d\n", Port)
	return nil
}

func handleHealthCheck(writer http.ResponseWriter, req *http.Request) {
	healthCheck := &HealthCheckResponse{
		Version: CommitID,
		BuildAt: BuildDate,
	}
	bytes, err := json.Marshal(healthCheck)

	if err != nil {
		writer.WriteHeader(http.StatusOK)
		return
	}

	fmt.Fprint(writer, string(bytes))

}

func handleTimeStamp(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	rawTimeStringInMilis, err := strconv.ParseInt(vars["timestamp"], 10, 64)

	if err != nil {
		// the timestamp cannot be parsed into a int64
		// attempt to parse it in "YYYY-MM-DD" format
		result, stringDateParseError := parseStringDate(vars["timestamp"])

		if stringDateParseError != nil {
			SendErrorResponse(writer)
			return
		}

		response := &Response{
			Unix: result.UnixMicro(),
			Utc:  result.Format(ResponseLayout),
		}

		SendSuccessResponse(writer, response)
		return
	}

	// extract seconds and nanoseconds from this
	seconds := rawTimeStringInMilis / 1000
	nanoSeconds := (rawTimeStringInMilis - 1000*seconds) * 1000
	timeStamp := time.Unix(seconds, nanoSeconds)

	// similar to time.RFC1123, but ends with GMT not UTC
	response := &Response{
		Unix: rawTimeStringInMilis,
		Utc:  timeStamp.UTC().Format(ResponseLayout),
	}

	SendSuccessResponse(writer, response)
}

func handleEmptyTimeStamp(writer http.ResponseWriter, req *http.Request) {
	now := time.Now()
	// similar to time.RFC1123, but ends with GMT not UTC
	res := &Response{Unix: now.UTC().UnixMicro(), Utc: now.UTC().Format(ResponseLayout)}
	SendSuccessResponse(writer, res)
}

/**
* Parse strings like "1978-09-04"
 */
func parseStringDate(dateString string) (time.Time, error) {
	parseLayout := "2006-01-02"
	result, err := time.Parse(parseLayout, dateString)

	if err != nil {
		// could not parse in that format as specified by parseLayout
		return time.Now(), err
	}

	return result, nil
}

func SendSuccessResponse(writer http.ResponseWriter, response *Response) {
	responseJson, err := json.Marshal(response)

	if err != nil {
		SendErrorResponse(writer)
		return
	}
	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, string(responseJson))
}

func SendErrorResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusInternalServerError)
	errorResponse, _ := json.Marshal(&ErrorResponse{
		ErrorValue: "Invalid Date",
	})
	fmt.Fprint(writer, string(errorResponse))
}
