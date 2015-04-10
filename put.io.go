package main

import (
	"encoding/json"
	"fmt"
	"github.com/jawher/mow.cli"
	"io/ioutil"
	"log"
	"net/http"
    "net/url"
	"os"
    "strconv"
)

type FileList struct {
	Files  []*File `json:"files"`
	Parent File    `json:"parent"`
	Status string  `json:"status"`
}

type File struct {
	ContentType       string `json:"content_type"`
	CRC32             string `json:"crc32"`
	CreatedAt         string `json:"created_at"`
	FirstAccessedAt   string `json:"first_accessed_at"`
	Icon              string `json:"icon"`
	Id                int    `json:"id"`
	IsMp4Available    bool   `json:"is_mp4_available"`
	IsShared          bool   `json:"is_shared"`
	Name              string `json:"name"`
	OpensubtitlesHash string `json:"opensubtitles_hash"`
	ParentId          int    `json:"parent_id"`
	Screenshot        string `json:"screenshot"`
	Size              int    `json:"size"`
}

func main() {

	prog := cli.App("put.io", "access files from your put.io account")

	// global options
	token := prog.StringOpt("t token", "", "your oauth token from put.io")

	// commands
	prog.Command("list", "list all files in your put.io account", func(cmd *cli.Cmd) {

        parent := cmd.IntArg("DIR", 0, "the id of the directory to list")

        fmt.Println(*parent)

		c := new(Client)
		c.Token = *token

        // get the list of files and print out the info we care about for each
		c.ListFiles(*parent)
	})

	prog.Run(os.Args)

}

type Client struct {
	Token string
}

func (c *Client) ListFiles(parent int)  {

    v := url.Values{}
    v.Set("oauth_token", c.Token)
    v.Set("parent_id", strconv.Itoa(parent))

    fmt.Println(v.Encode())

    req := &http.Request{
        Method: "GET",
        Host:   "api.put.io",
        URL: &url.URL{
            Host: "api.put.io",
            Scheme: "https",
            Path: "/v2/files/list",
            RawQuery: v.Encode(),
        },
    }

    var client http.Client
    resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var fl FileList
	err = json.Unmarshal(body, &fl)

    for _, file := range fl.Files {
        fmt.Println(file.Id, file.Name)
    }
}
