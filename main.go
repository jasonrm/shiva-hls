package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/kyoh86/xdg"
	"gitlab.com/jasonrm/shiva-hls/source"

	"github.com/quangngotan95/go-m3u8/m3u8"
)

var (
	outDir             = flag.String("out", "twitch", "output directory")
	TwitchClientId     = flag.String("twitch-client-id", "", "")
	TwitchClientSecret = flag.String("twitch-client-secret", "", "")
)

func downloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		os.Remove(filepath)
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(filepath)
		panic(err)
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
				fmt.Printf("download: %s\n", outFile)
				_ = downloadFile(outFile, remoteUri.String())
				continue
			}
			//fmt.Printf("exists: %s\n", file)
		}
	}(queue)

	return queue
}

func main() {
	flag.Parse()
	twitchUsername := flag.Arg(0)

	dbPath := path.Join(xdg.DataHome(), "shiva-hls.db")
	cacheDb, _ := sql.Open("sqlite3", dbPath)

	twitch := source.NewTwitch(*TwitchClientId, *TwitchClientSecret, cacheDb)

	for _, video := range twitch.Videos(twitchUsername) {
		fmt.Println(path.Base(video.URL))
		dlDir := path.Join(path.Clean(*outDir), strings.ToLower(twitchUsername), path.Base(video.URL))
		_ = os.MkdirAll(dlDir, os.ModePerm)

		ytOut := exec.Command("youtube-dl", "-g", video.URL)
		ytPlaylist, ytErr := ytOut.Output()
		if ytErr != nil {
			panic(ytErr)
		}
		playlistUrl := string(ytPlaylist)
		playlistUrl = strings.TrimSpace(playlistUrl)

		queue := downloadQueue(playlistUrl, dlDir)

		resp, err := http.Get(playlistUrl)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		playlist, err := m3u8.ReadString(string(contents))
		if err != nil {
			panic(err)
		}

		if playlist.IsValid() {
			outPlaylist := path.Join(dlDir, "playlist.m3u8")
			out, err := os.Create(outPlaylist)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(out, bytes.NewBuffer(contents))
			if err != nil {
				panic(err)
			}
		}

		fmt.Printf("playlist: %s\n", playlistUrl)

		c := 0
		for _, i := range playlist.Items {
			if i, ok := i.(*m3u8.SegmentItem); ok {
				c++
				queue <- i.Segment
			}
		}

		outMetadata := path.Join(dlDir, "metadata.json")
		out, err := os.Create(outMetadata)
		if err != nil {
			panic(err)
		}
		metadata, err := json.Marshal(video)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(out, bytes.NewBuffer(metadata))
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Minute)
	}
}
