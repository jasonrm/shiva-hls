package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/quangngotan95/go-m3u8/m3u8"
)

var (
	playlistUrl = flag.String("url", "", "playlist location")
	outDir      = flag.String("out", ".", "output directory")
)

func downloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloadQueue(playlistUrl string, outDir string) chan string {
	queue := make(chan string)

	u, _ := url.Parse(playlistUrl)

	go func(chan string) {
		for file := range queue {
			outFile := path.Join(outDir, file)
			remoteUri, _ := url.Parse(u.String())
			remoteUri.Path = path.Join(path.Dir(u.Path), file)
			if _, err := os.Stat(outFile); os.IsNotExist(err) {
				fmt.Printf("download: %s -> %s\n", remoteUri, outFile)
				_ = downloadFile(outFile, remoteUri.String())
				continue
			}
			fmt.Printf("exists: %s\n", file)
		}
	}(queue)

	return queue
}

func main() {
	flag.Parse()

	queue := downloadQueue(*playlistUrl, path.Base(*outDir))

	_ = os.MkdirAll(*outDir, os.ModePerm)

	resp, err := http.Get(*playlistUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// show the HTML code as a string %s
	//fmt.Printf("%s\n", contents)

	playlist, err := m3u8.ReadString(string(contents))
	if err != nil {
		panic(err)
	}

	if playlist.IsValid() {
		outPlaylist := path.Join(*outDir, "playlist.m3u8")
		out, err := os.Create(outPlaylist)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(out, bytes.NewBuffer(contents))
		if err != nil {
			panic(err)
		}
	}

	for _, i := range playlist.Items {
		if i, ok := i.(*m3u8.SegmentItem); ok {
			queue <- i.Segment
		}
	}
}
