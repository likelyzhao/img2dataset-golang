package resizer

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"img2dataset/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/h2non/bimg"
)

func getUrl(url string) (data []byte, err error) {
	ret, err := http.Get(url)
	if err != nil {
		log.Println(url)
		return data, err
	}
	body := ret.Body

	data, _ = io.ReadAll(body)

	return data, nil
}

func Processor(in <-chan utils.ImageInfo, savedir string) {

	for {
		info := <-in

		// log.Println("starting image with Index " + strconv.FormatInt(info.Id, 10))
		// open "test.jpg"
		//savedir := "out2/"

		os.Mkdir(savedir, 0777)

		saveimagepath := savedir + strconv.FormatInt(info.Id, 10) + ".jpg"
		flag, _ := utils.PathExists(saveimagepath)
		if flag {
			//log.Println("image with Index " + strconv.FormatInt(info.Id, 10) + " finished")
			continue
		}
		savejsonpath := savedir + strconv.FormatInt(info.Id, 10) + ".json"
		flag, _ = utils.PathExists(savejsonpath)
		if flag {
			//log.Println("image with Index " + strconv.FormatInt(info.Id, 10) + " finished")
			continue
		}

		buffer, err := getUrl(info.Url)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			log.Println("image with Index " + strconv.FormatInt(info.Id, 10) + " get url error")
			continue
		}

		// check the buffer
		newImage_img := bimg.NewImage(buffer)
		if newImage_img.Type() == "unknown" {
			log.Println("image with Index " + strconv.FormatInt(info.Id, 10) + " format error")
			continue
		}

		if newImage_img.Type() != "jpeg" {
			buffer, err = bimg.NewImage(buffer).Convert(bimg.JPEG)
			if err != nil {
				log.Println("image with Index " + strconv.FormatInt(info.Id, 10) + " change format error")
				continue
			}
		}

		newImage, err := Resizer(buffer)
		h := sha256.New()
		h.Write(buffer)
		bs := h.Sum(nil)
		info.Sha256 = fmt.Sprintf("%x", bs)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		file, _ := os.OpenFile(savejsonpath, os.O_CREATE|os.O_WRONLY, 0)

		enc := json.NewEncoder(file)
		err = enc.Encode(info)
		if err != nil {
			log.Println("Error in encoding json")
			continue
		}

		size, _ := bimg.NewImage(newImage).Size()
		info.ResizeHeight = int64(size.Height)
		info.ResizeWidth = int64(size.Width)
		bimg.Write(saveimagepath, newImage)

		//log.Println("saveing image in " + saveimagepath)
		//log.Println("image with Index " + strconv.FormatInt(info.Id, 10) + " finished")
	}
}

func Resizer(buffer []byte) ([]byte, error) {
	// open "test.jpg"
	size, _ := bimg.NewImage(buffer).Size()
	if size.Width <= 256 || size.Height <= 256 {
		return buffer, nil
	} else {
		newImage, err := bimg.NewImage(buffer).Resize(256, 256)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return newImage, err
	}
}
