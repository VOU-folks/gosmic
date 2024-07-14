package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func DownloadPBF(url string, folderPath string, fileName string) error {
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

	fmt.Printf("Download completed. Total size: %d bytes\n\n", size)
	return nil
}

type progressReader struct {
	reader       io.Reader
	total        int64
	downloaded   int64
	lastReported int64
	finished     bool
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	p.downloaded += int64(n)
	return n, err
}

func (p *progressReader) start() {
	p.finished = false
	for {
		if p.finished {
			break
		}

		select {
		case <-time.After(100 * time.Millisecond):
			p.report()
		}
	}
}

func (p *progressReader) report() {
	percentage := float64(p.downloaded) / float64(p.total) * 100

	if p.downloaded-p.lastReported > p.total/10000 || p.downloaded == p.total {
		fmt.Printf("\rDownloading... %.2f%% complete (%v / %v)", percentage, p.downloaded, p.total)
		p.lastReported = p.downloaded
	}
}

func (p *progressReader) done() {
	p.finished = true
	fmt.Printf("\rDownloading... 100.00%% complete (%v / %v) \n", p.downloaded, p.total)
}
