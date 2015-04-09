package main

import (
	"encoding/json"
	"fmt"
	"github.com/jawher/mow.cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

type Client struct {
	Token string
}

func main() {

	prog := cli.App("put.io", "Access files from your put.io account")

	// global options
	token := prog.StringOpt("t token", "", "your oauth token from put.io")

	// commands
	prog.Command("list", "list all files in your put.io account", func(cmd *cli.Cmd) {

		c := new(Client)
		c.Token = *token

		files := c.ListFiles().Files

		for _, file := range files {
			fmt.Println(file.Id, file.Name)
		}
	})

	prog.Run(os.Args)

}

func (c *Client) Do(method string, path string) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, fmt.Sprint("https://api.put.io/v2/", path, "?oauth_token=", c.Token), nil)

	if err != nil {
		log.Fatal(err)
	}

	var client http.Client

	return client.Do(req)
}

func (c *Client) ListFiles() FileList {
	resp, err := c.Do("GET", "files/list")

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var fl FileList
	err = json.Unmarshal(body, &fl)

	return fl
}
