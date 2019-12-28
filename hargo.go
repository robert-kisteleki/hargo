package main

import (
    "encoding/json"
    "encoding/base64"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

var getversion bool
var verbose bool
var listonly bool
var filter_url_inc string
var filter_url_exc string
var filter_method string
var prefix string
var suffix string
var sequences bool
var current_sequence int

const (
    version = "0.0.1"
    filename_max_len = 128
)

func main() {
    current_sequence = 1

    parse_flags()

    if getversion {
        fmt.Println("hargo version", version)
        return
    }

    byteValue, _ := ioutil.ReadAll(os.Stdin)

    var har map[string]interface{}
    json.Unmarshal([]byte(byteValue), &har)

    process_json_log(har["log"].(map[string]interface{}))
}

func parse_flags() {
    flag.BoolVar(&getversion, "V", false, "Display version information")
    flag.BoolVar(&verbose, "v", false, "Be verbose")
    flag.BoolVar(&listonly, "l", false, "Don't save anything, just list the URLs")
    flag.StringVar(&filter_url_inc, "i", "", "Match only URLs with this string in them")
    flag.StringVar(&filter_url_exc, "e", "", "Match only URLs without this string in them")
    flag.StringVar(&filter_method, "m", "", "Match only these methods (GET, POST, ...)")
    flag.StringVar(&prefix, "p", "", "Use this prefix for saved files. Can include local directory name")
    flag.StringVar(&suffix, "s", "", "Use this suffix for saved files")
    flag.BoolVar(&sequences, "n", false, "Use sequence number instead of URL for saved files")
    flag.Parse()
}

func process_json_log(log map[string]interface{}) {
    process_json_entries(log["entries"].([]interface{}))
}

func process_json_entries(entries []interface{}) {
    for _, v := range entries {
        process_json_entry(v.(map[string]interface{}))
    }
}

func process_json_entry(entry map[string]interface{}) {
    var url string
    var method string
    var content []byte
    for k, v := range entry {
        if k == "request" {
            url, method = process_json_request(v.(map[string]interface{}))
        }
        if k == "response" {
            content = process_json_response(v.(map[string]interface{}))
        }
    }

    if listonly {
        fmt.Printf("URL (%s): %s\n", method, url)
    } else {
        process_candidate(url, method, &content)
    }
}

func process_candidate(url string, method string, content *[]byte) {
    if (filter_method=="" || filter_method==method) &&
       (filter_url_inc=="" ||  strings.Contains(url, filter_url_inc)) &&
       (filter_url_exc=="" || !strings.Contains(url, filter_url_exc)) {
        if verbose {
            fmt.Println("Filter matched for URL:", url)
        }
        out_name := prefix
        if sequences {
            out_name += fmt.Sprintf("%06d", current_sequence)
            current_sequence++
        } else {
            url = strings.TrimPrefix(url, "http://")
            url = strings.TrimPrefix(url, "https://")
            out_name += strings.SplitN(url, "/", 2)[1][:filename_max_len]
        }
        out_name += suffix
        if verbose {
            fmt.Println("Writing to:", out_name)
        }

        f, err := os.Create(out_name)
        if err != nil {
           panic(err)
        }
        defer f.Close()
        f.Write(*content)
    }
}

func process_json_request(request map[string]interface{}) (string,string) {
    url := "URL"
    method := "METHOD"
    for k, v := range request {
        if k == "url" {
            url = v.(string)
        }
        if k == "method" {
            method = v.(string)
        }
    }
    return url, method
}


func process_json_response(response map[string]interface{}) []byte {
    for k, v := range response {
        if k == "content" {
            return process_json_content(v.(map[string]interface{}))
        }
    }
    return nil
}

func process_json_content(content map[string]interface{}) []byte {
    text := []byte("CONTENT")
    temp_text := ""
    encoding := ""

    for k, v := range content {
        if k == "text" {
            temp_text = v.(string)
        }
        if k == "encoding" {
            encoding = v.(string)
        }
    }

    if encoding == "base64" {
        temp_decoded, err := base64.StdEncoding.DecodeString(temp_text)
        if err != nil {
            fmt.Println("ERROR: decode error:", err)
            return []byte("")
        }
        text = []byte(temp_decoded)
    } else {
        text = []byte(temp_text)
    }

    return text
}
