package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	. "gosmic/internal/config"
)

func main() {
	var err error

	configFile := flag.String("config", "config.yaml", "path to the config file. Example: api -config /full/path/to/config.yaml")
	flag.Parse()

	config, err := GetConfigService(*configFile)
	if err != nil {
		panic(err)
	}

	storageConfig := config.GetStorage()
	osmConfig := config.GetOsmConfig()

	pbfSource := osmConfig.Sources.PBF.Url
	pbfFileName := osmConfig.Sources.PBF.FileName
	pbfStorage := storageConfig.PBFs

	err = ensureFolderExists(pbfStorage)
	if err != nil {
		panic(err)
	}

	err = downloadPBF(pbfSource, pbfStorage, pbfFileName)
	if err != nil {
		panic(err)
	}
}

func ensureFolderExists(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadPBF(url string, folderPath string, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	outPath := filepath.Join(folderPath, fileName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	progress := &progressReader{reader: resp.Body, total: resp.ContentLength}
	go progress.start()

	size, err := io.Copy(outFile, progress)
	if err != nil {
		return err
	}
	progress.done()

	fmt.Printf("\nDownload completed. Total size: %d bytes\n", size)
	return nil
}

type progressReader struct {
	reader       io.Reader
	total        int64
	downloaded   int64
	lastReported int64
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	p.downloaded += int64(n)
	return n, err
}

func (p *progressReader) start() {
	for {
		select {
		case <-time.After(time.Second):
			p.report()
		}
	}
}

func (p *progressReader) report() {
	percentage := float64(p.downloaded) / float64(p.total) * 100

	if p.downloaded-p.lastReported > p.total/10000 || p.downloaded == p.total {
		fmt.Printf("\rDownloading... %.2f%% complete", percentage)
		p.lastReported = p.downloaded
	}
}

func (p *progressReader) done() {
	fmt.Printf("\rDownloading... 100.00%% complete\n")
}
