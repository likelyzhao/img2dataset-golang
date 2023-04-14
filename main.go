package main

import (
	"img2dataset/loader"
)

func main() {

	//loader.LoaderParguet("/data/zhaozhijian/chinese/part-00002-fc82da14-99c9-4ff6-ab6a-ac853ac82819-c000.snappy.parquet")
	loader.LoaderTsvhea("/data/zhaozhijian/img2dataset-golang/cc12m.tsv", "out3/")
}
