package main

import (
    "net/http"
    "fmt"
    "log"
    "io/ioutil"
    "encoding/json"
)

type FileList struct {
    Files   []*File `json:"files"`
    Parent  File    `json:"parent"`
    Status  string  `json:"status"`
}

type File struct {
    ContentType string  `json:"content_type"`
    CRC32   string  `json:"crc32"`
    CreatedAt   string  `json:"created_at"`
    FirstAccessedAt string  `json:"first_accessed_at"`
    Icon    string  `json:"icon"`
    Id  int `json:"id"`
    IsMp4Available  bool    `json:"is_mp4_available"`
    IsShared    bool    `json:"is_shared"`
    Name    string  `json:"name"`
    OpensubtitlesHash   string  `json:"opensubtitles_hash"`
    ParentId    int `json:"parent_id"`
    Screenshot  string  `json:"screenshot"`
    Size    int `json:"size"`
}

type Client struct {
    Token   string
}

func (c *Client) Get(path string) (resp *http.Response, err error) {
    return http.Get(fmt.Sprint("https://api.put.io/v2/", path, "?oauth_token=", c.Token))
}

func (c *Client) ListFiles() {
    resp, err := c.Get("files/list")

    if err != nil {
        log.Fatal(err)
        return
    }

    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        log.Fatal(err)
        return
    }

    var f interface{}

    err = json.Unmarshal(body, &f)

    fmt.Print(f)

}

func main() {

    c := new(Client)
    c.Token = ""

    c.ListFiles()
}
