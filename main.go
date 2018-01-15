package main

import (
    "log"
    "fmt"
    "os"
    "strings"
    "bytes"
    "encoding/json"
    "net/http"
    "github.com/gorilla/schema"
)

type Payload struct {
    Text         string
    Channel_name string
    User_name    string
}

type Worklog struct {
    Time_spent    string `json:"timeSpent"`
    Comment       string `json:"comment"`
}

type Response struct {
    Response_type string `json:"response_type"`
    Text          string `json:"text"`
}

func main() {
    router := http.NewServeMux()
    router.HandleFunc("/", AddWorklog)
    log.Fatal(http.ListenAndServe(":4000", router))
}

func AddWorklog(writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm()
    if err != nil {
        fmt.Println("Error parsing form")
    }

    payload := new(Payload)
    decoder := schema.NewDecoder()
    decoder.IgnoreUnknownKeys(true)

    fmt.Println(request.PostForm)

    err = decoder.Decode(payload, request.Form)
    if err != nil {
        fmt.Println("Error decoding")
    }

    fmt.Println(payload.Text)

    words := strings.Fields(payload.Text)
    issue := words[0]
    time  := words[1]

    worklog := Worklog{Time_spent: time, Comment: ""} 
    buffer  := new(bytes.Buffer)
    json.NewEncoder(buffer).Encode(worklog)

    var url = os.Getenv("MATTERMOST_JIRA_HOST") + "/rest/api/2/issue/" + issue + "/worklog"

    client := &http.Client{}
    req, err := http.NewRequest("POST", url, buffer)
    req.SetBasicAuth(os.Getenv("MATTERMOST_JIRA_USERNAME"), os.Getenv("MATTERMOST_JIRA_PASSWORD"))
    req.Header.Add("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil{
        log.Fatal(err)
    }

    fmt.Println(resp.Body)

    writer.Header().Set("Content-Type", "application/json")

    response := Response{Response_type: "in_channel", Text: "@" + payload.User_name + " spent " + time + " on " + issue}
    json.NewEncoder(writer).Encode(response)
}
