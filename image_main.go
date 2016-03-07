package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func NewFileNameFromURL(url string) string {
	hash := md5.New()
	io.WriteString(hash, url)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func download_image(rawURL string) (string, string, error) {
	var filename string
	//fmt.Println("Downloading file...")
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := check.Get(rawURL) // add a filter to check redirect
	if err != nil {
		return "", "", err
	}

	content_type := resp.Header.Get("Content-Type")
	var ext string = ""
	//fmt.Printf("type(content_type): %s\n", reflect.TypeOf(content_type))
	switch content_type {
	case "image/jpeg":
		fallthrough
	case "image/jpg":
		ext = ".jpg"
		break
	case "image/png":
		ext = ".png"
		break
	case "image/gif":
		ext = ".gif"
		break
	default:
		return "", "", fmt.Errorf("Unknown content-type:%s", content_type)
	}

	content_disposition := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(content_disposition)
	if err == nil {
		filename = params["filename"]
	}

	//fmt.Printf("resp: %s\n", resp)
	//fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	//fmt.Printf("Content-Disposition: %s\n", resp.Header.Get("Content-Disposition"))
	//fmt.Printf("Disposition: %s\n", disposition)
	//fmt.Printf("Filename: %s\n", filename)
	//fmt.Printf("Ext: %s\n", ext)

	defer resp.Body.Close()

	//fmt.Println(resp.Status)

	if len(filename) == 0 {
		fileURL, err := url.Parse(rawURL)
		if err != nil {
			return "", "", fmt.Errorf("Parsing Url: %s\n", rawURL)
		}

		path := fileURL.Path
		segments := strings.Split(path, "/")
		//fmt.Printf("path: %s\n", path)
		//fmt.Printf("path segments: %s\n", segments)
		//for i, val := range segments {
		//	fmt.Printf("path[%d]=%s\n", i, val)
		//}
		if len(segments) > 3 {
			filename = segments[len(segments)-1] // change the number to accommodate changes to the url.Path position
		} else {
			filename = NewFileNameFromURL(rawURL) + ext
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return "", "", err
	}

	//fmt.Printf("[%s] with %v bytes downloaded\n", filename, size)

	return content_type, filename, nil
}

func main() {
	//rawURL := "https://d1ohg4ss876yi2.cloudfront.net/golang-resize-image/big.jpg"
	//rawURL := "http://www.ban8.co.kr/shopimages/ban8/020000000506.jpg"
	//rawURL := "http://cfile240.uf.daum.net/image/1761CF384FF2F0441BE32D"
	rawURL := os.Args[1]

	_, filename, err := download_image(rawURL)
	if err != nil {
		fmt.Printf("download err: %s\n", err)
		return
	}

	fmt.Printf("downloaded %s\n", filename)
}
