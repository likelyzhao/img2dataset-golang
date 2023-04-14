package loader

import (
	"log"

	"img2dataset/resizer"
	"img2dataset/utils"
	"bufio"
	"os"
	"fmt"
	"strings"
	"net/url"
	

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func LoaderTsvhea(filepath string, savedir string) []utils.ImageInfo {
	var res []utils.ImageInfo

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return res
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var total [][]string
	HostNames := make(map[string]string)
	//var strSlice []string
	for scanner.Scan() {
		str:= scanner.Text()
		parts := strings.Split(str, "\t")
		total = append(total, parts)

		u, err := url.Parse(parts[0])
		if err != nil {
		 panic(err)
		}
	   
		// Get the host from the URL
		host := u.Hostname()
		if _, ok := HostNames[host];!ok{
			HostNames[host] = ""
			//strSlice = append(strSlice,  host)
		}
	
	}
	fmt.Println(HostNames)

/*
	{
		file, err := os.Create("output.txt")
		if err != nil {
		panic(err)
		}
		defer file.Close()

		// Join the slice into a single string with newline separato
		str := strings.Join(strSlice, "\n")

		// Write the string to the file
		_, err = file.WriteString(str)
		if err != nil {
		panic(err)
		}

		fmt.Println("String slice written to file successfully!")

	}
 // Read the entire contents of a file into a byte slice

*/
	num := int(len(total))
	// num := 10
	ch := make(chan utils.ImageInfo, 20)

	for i := 0; i < 1000; i++ {
		go resizer.Processor(ch, savedir)
	}

	bar := utils.InitProgressBar(num)
	for i := 0; i < num; i++ {
		stus := make([]utils.ImageInfo, 1) //read 10 rows
		stus[0].Id = int64(i)
		stus[0].Url = total[i][0]
		stus[0].Text = total[i][1]

		for _, info := range stus {
			ch <- info
		}

		bar.Step()
	}

	return res
}



func LoaderParguet(filepath string, savedir string) []utils.ImageInfo {
	var res []utils.ImageInfo
	fr, err := local.NewLocalFileReader(filepath)
	if err != nil {
		log.Println("Can't open file")
		return res
	}

	pr, err := reader.NewParquetReader(fr, new(utils.ImageInfo), 4)
	if err != nil {
		log.Println("Can't create parquet reader", err)
		return res
	}
	num := int(pr.GetNumRows())

	ch := make(chan utils.ImageInfo, 20)

	for i := 0; i < 10000; i++ {
		go resizer.Processor(ch, savedir)
	}

	bar := utils.InitProgressBar(num / 100)
	for i := 0; i < num/100; i++ {
		stus := make([]utils.ImageInfo, 100) //read 10 rows
		if err = pr.Read(&stus); err != nil {
			log.Println("Read error", err)
		}
		for _, info := range stus {
			ch <- info
		}

		bar.Step()
	}

	pr.ReadStop()
	fr.Close()

	return res
}
