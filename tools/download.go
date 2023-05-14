package tools

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/cavaliergopher/grab/v3"
	"github.com/dustin/go-humanize"
)

var MaxRetry = 6

var Get = make(chan string, 4)

type Activity struct {
	Address, ResponseCode, Status string
	Progress                      float32
	Update                        func()
}

var Activites []*Activity

func Listen(Update func()) {
	for {
		x := <-Get
		y := new(Activity)
		y.Address = x
		y.Update = Update
		go y.Download(x)
		Activites = append(Activites, y)
	}
}

func (A *Activity) Download(Address string) {

	var tries = 1

	var resp = new(grab.Response)

	// var err error
	client := grab.NewClient()

	req, _ := grab.NewRequest(filepath.Join(xdg.UserDirs.Download, "gogetit")+"/", Address)

	req.HTTPRequest.Referer()

	A.Status = fmt.Sprintf("Requesting :%v ", Address)
	A.Update()

	for resp = client.Do(req); resp.HTTPResponse == nil; {
		tries++
		if tries > MaxRetry {
			A.Status = fmt.Sprintf("Download Failed With %d Tries: %v", tries-1, resp.Err())
			A.Update()

			log.Printf("Download Failed With %d Tries: %v", tries-1, resp.Err())

			return
		}
		time.Sleep(time.Second * time.Duration(tries))
		A.Status = fmt.Sprintf("Retrying #%d because:  %v ", tries, resp.Err())
		A.Update()
		log.Println("Retrying because: ", resp.Err())

	}

	if resp.HTTPResponse != nil {
		log.Println(resp.HTTPResponse.Status)
		A.Address = filepath.Base(resp.Filename)
		A.ResponseCode = resp.HTTPResponse.Status
	}

	fmt.Printf("Downloading %v...\n", req.URL())

	t := time.NewTicker(500 * time.Millisecond)

	defer t.Stop()

Loop:
	for {
		select {

		case <-t.C:
			A.Progress = float32(resp.Progress())

			A.Status = fmt.Sprintf("Transferred %v / %v bytes (%.2f%%)\n", humanize.Bytes(uint64(resp.BytesComplete())), humanize.Bytes(uint64(resp.Size())), 100*resp.Progress())
			A.Update()

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		A.Status = fmt.Sprintf("Download failed: %v\n", err)
		A.Update()
		return
	}

	A.Status = fmt.Sprintf("Download saved to %v \n", resp.Filename)
	A.Update()

}
