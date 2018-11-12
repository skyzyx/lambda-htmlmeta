package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dyatlov/go-htmlinfo/htmlinfo"
	// "github.com/namsral/microdata"
	"willnorris.com/go/microformats"
)

var (
	args            []string
	body            []byte
	data            interface{}
	emptyAGW        *events.APIGatewayProxyResponse
	err             error
	isMock          *bool
	isOpenGraph     bool
	isOEmbed        bool
	isMicrodata     bool
	isMicroformats2 bool
	statusCode      int
	u               string
)

// The API Gateway handler
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	statusCode = int(200)
	emptyAGW = new(events.APIGatewayProxyResponse)
	cacheFrom := time.Now().Format(http.TimeFormat)
	cacheUntil := time.Now().AddDate(1, 0, 0).Format(http.TimeFormat)

	if request.QueryStringParameters["url"] == "" {
		return *emptyAGW, errors.New("The 'url' query string parameter is required.")
	}

	u = request.QueryStringParameters["url"]

	if !govalidator.IsURL(u) {
		return *emptyAGW, errors.New("This is not a valid URL.")
	}

	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if isOpenGraph, err = strconv.ParseBool(os.Getenv("META_OPENGRAPH")); err != nil {
		isOpenGraph = false
	}

	if isOEmbed, err = strconv.ParseBool(os.Getenv("META_OEMBED")); err != nil {
		isOEmbed = false
	}

	if isMicrodata, err = strconv.ParseBool(os.Getenv("META_MICRODATA")); err != nil {
		isMicrodata = false
	}

	if isMicroformats2, err = strconv.ParseBool(os.Getenv("META_MICROFORMATS2")); err != nil {
		isMicroformats2 = false
	}

	if isOpenGraph || isOEmbed || isMicrodata {
		info := htmlinfo.NewHTMLInfo()

		err = info.Parse(resp.Body, &u, nil)
		if err != nil {
			log.Fatal(err)
		}

		if isOpenGraph {
			data = info.OGInfo
		} else if isOEmbed {
			data = info.GenerateOembedFor(u)
		} else if isMicrodata {
			data = info
		}

		body, err = json.MarshalIndent(data, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
	} else if isMicroformats2 {
		URL, err := url.Parse(u)
		if err != nil {
			panic(err)
		}

		info := microformats.Parse(resp.Body, URL)

		body, err = json.MarshalIndent(info, "", "    ")
	}

	// HTTP response as JSON
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Last-Modified": cacheFrom,
			"Expires":       cacheUntil,
		},
		Body:       string(body),
		StatusCode: statusCode,
	}, nil
}

// The core function
func main() {
	isMock = flag.Bool("mock", false, "Read from the local `mock.json` file instead of an API Gateway request.")
	flag.Parse()

	if *isMock {
		// read json from file
		inputJSON, jsonErr := ioutil.ReadFile("./mock.json")
		if jsonErr != nil {
			fmt.Println(jsonErr.Error())
			os.Exit(1)
		}

		// de-serialize into Go object
		var inputEvent events.APIGatewayProxyRequest
		if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var response events.APIGatewayProxyResponse
		response, err = Handler(inputEvent)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(response)
	} else {
		lambda.Start(Handler)
	}
}
