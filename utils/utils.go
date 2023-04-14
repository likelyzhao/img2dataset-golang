package utils

import (
	"os"

	"github.com/schollz/progressbar/v3"
)

type ImageInfo struct {
	Id           int64   `json:"key" parquet:"name=SAMPLE_ID,type=INT64,repetitiontype=OPTIONAL"`
	Url          string  `json:"url" parquet:"name=URL,type=BYTE_ARRAY,convertedtype=UTF8,repetitiontype=OPTIONAL"`
	Text         string  `json:"caption" parquet:"name=TEXT,type=BYTE_ARRAY,convertedtype=UTF8, repetitiontype=OPTIONAL"`
	Height       int64   `parquet:"name=HEIGHT,type=INT64,repetitiontype=OPTIONAL"`
	Width        int64   `parquet:"name=WIDTH,type=INT64,repetitiontype=OPTIONAL"`
	License      string  `parquet:"name=LICENSE,type=BYTE_ARRAY,convertedtype=UTF8, repetitiontype=OPTIONAL"`
	Language     string  `parquet:"name=LANGUAGE,type=BYTE_ARRAY,convertedtype=UTF8, repetitiontype=OPTIONAL"`
	NSFW         string  `parquet:"name=NSFW,type=BYTE_ARRAY,convertedtype=UTF8,repetitiontype=OPTIONAL"`
	Similarity   float64 `parquet:"name=similarity,type=DOUBLE,repetitiontype=OPTIONAL"`
	Sha256       string  `json:"sha256"`
	ResizeHeight int64   `json:"NewHeight"`
	ResizeWidth  int64   `json:"NewWidth"`
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type ProgressBar struct {
	bar *progressbar.ProgressBar
}

func InitProgressBar(Length int) ProgressBar {
	var bar ProgressBar
	bar.bar = progressbar.Default(int64(Length))
	return bar
}

func (b ProgressBar) Step() {
	b.bar.Add(1)
}
