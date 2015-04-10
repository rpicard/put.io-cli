package main

import (
	"encoding/json"
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/pivotal-golang/bytefmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
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

		// DIR is optional
		cmd.Spec = "[DIR]"

		parent := cmd.IntArg("DIR", 0, "the id of the directory to list")

		cmd.Action = func() {

			c := new(Client)
			c.Token = *token

			// get the list of files and print out the info we care about for each
			c.ListFiles(*parent)
		}
	})

	prog.Command("download", "download a specific file (no folders)", func(cmd *cli.Cmd) {

		id := cmd.IntArg("FILE", 0, "the id of the file to download")

		cmd.Action = func() {

			c := new(Client)
			c.Token = *token

			// download the file to the current working directory
			c.DownloadFile(*id)
		}
	})

	prog.Run(os.Args)

}

type Client struct {
	Token string
}

func (c *Client) DownloadFile(id int) {

	v := url.Values{}
	v.Set("oauth_token", url.QueryEscape(c.Token))

	req := &http.Request{
		Method: "GET",
		Host:   "api.put.io",
		URL: &url.URL{
			Host:     "api.put.io",
			Scheme:   "https",
			Path:     fmt.Sprintf("/v2/files/%d/download", id),
			RawQuery: v.Encode(),
		},
	}

	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("bad status code from server: %s", resp.Status))
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))

	if err != nil {
		log.Fatal(err)
	}

	// avoid directory traversal
	filename := filepath.Base(params["filename"])

	// O_EXCL will not open the file if it already exists
	outfile, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	written, err := io.Copy(outfile, resp.Body)

	if err != nil {
		// clean up the file
		os.Remove(outfile.Name())
		log.Fatal(err)
	}

	fmt.Printf("%s saved to %s\n", bytefmt.ByteSize(uint64(written)), outfile.Name())

}

func (c *Client) ListFiles(parent int) {

	v := url.Values{}
	v.Set("oauth_token", url.QueryEscape(c.Token))
	v.Set("parent_id", url.QueryEscape(strconv.Itoa(parent)))

	req := &http.Request{
		Method: "GET",
		Host:   "api.put.io",
		URL: &url.URL{
			Host:     "api.put.io",
			Scheme:   "https",
			Path:     "/v2/files/list",
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

	// use a text/tabwriter to align things when they are printed
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)

	for _, file := range fl.Files {
		fmt.Fprintf(w, "%s\t%d\t%s\n", file.ContentType, file.Id, file.Name)
	}

	w.Flush()
}
