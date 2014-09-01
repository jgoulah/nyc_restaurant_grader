package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "net/http"
    "net/url"
    "encoding/json"
    "flag"
    "strings"
)

// command flags
var verbose bool
var restaurant string

func main() {
    flag.BoolVar(&verbose, "verbose", false, "display violations")
    flag.Parse()
    extra_args := flag.Args()


    if len(extra_args) == 0 {
        fmt.Printf("Please supply a restaurant to search for!\n")
        os.Exit(0)
    }

    restaurant = strings.Join(extra_args, " ")
    fmt.Printf("Search for: '%s'\n\n", restaurant)
    get_content()
}

type GradeType struct {
    Restaurant    string `json:"dba"`
    Street        string `json:"street"`
    Grade         string `json:"grade"`
    GradeDate     string `json:"grade_date"`
    ViolationDesc string `json:"violation_description"`
    CuisineDesc   string `json:"cuisine_description"`
    ZipCode       string `json:"zipcode"`
}

func http_get(Url *url.URL) ([]GradeType) {

    res, err := http.Get(Url.String())

    if err != nil {
        panic(err.Error())
    }

    body, err := ioutil.ReadAll(res.Body)

    if err != nil {
        panic(err.Error())
    }

    var data []GradeType

    json.Unmarshal(body, &data)
    return data
}

func get_base_url() (*url.URL) {
    var Url *url.URL
    Url, err := url.Parse("http://data.cityofnewyork.us")
    if err != nil {
        panic("boom")
    }
    Url.Path += "/resource/xx67-kt59.json"
    return Url
}

func get_content() {

    // construct the exact match url
    Url := get_base_url()
    parameters := url.Values{}
    parameters.Add("dba", restaurant)
    Url.RawQuery = parameters.Encode()

    data := http_get(Url)
    if len(data) == 0 {
        fmt.Printf("Couldn't find exact match, trying fuzzy search...\n\n")
        // construct the fuzzy url
        Url := get_base_url()
        parameters := url.Values{}
        parameters.Add("$q", restaurant)
        Url.RawQuery = parameters.Encode()

        data = http_get(Url)
        if len(data) == 0 {
            fmt.Printf("Couldn't find match, giving up!\n")
            os.Exit(0)
        }
    }

    var displayedInfo = false
    violations := make([]string, 0)
    for _,v := range data {
        if len(v.Grade) == 0 {
            continue;
        }

        if !displayedInfo {
            displayedInfo = true
            fmt.Printf("%s\n", v.Restaurant)
            fmt.Printf("%s %s\n\n", v.Street, v.ZipCode)
            fmt.Printf("Grade: %s\n\n", v.Grade)
        }

        violations = append(violations, v.ViolationDesc)
    }

    if verbose && len(violations) > 0  {
        fmt.Printf("Violations:\n\n")
        for _,v := range violations {
            fmt.Printf(" - %s\n", v)
        }
    }


    os.Exit(0)
}
