package msglist

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"sync"
)

type pdfFilename struct {
	filename string
	datetime int
}

var (
	infoChannel           chan *MessageInfo
	generatedPdfChannel   chan *pdfFilename
	generatedPdfFilenames []*pdfFilename
	syncWaitGroup         sync.WaitGroup
)

func init() {
	if err := testCommand("wkhtmltopdf"); err != nil {
		log.Fatalf("wkhtmltopdf is not installed: %s\n", err)
	}
	if err := testCommand("gs"); err != nil {
		log.Fatalf("ghostscript is not installed: %s\n", err)
	}
	infoChannel = make(chan *MessageInfo, 100)
	generatedPdfChannel = make(chan *pdfFilename)
	generatedPdfFilenames = make([]*pdfFilename, 0, 500)
	for i := 0; i < 10; i++ {
		go pdfGenerator(infoChannel, generatedPdfChannel)
	}
	syncWaitGroup.Add(10)
	go generatedPdfCollector(generatedPdfChannel)
}

func testCommand(command string) error {
	cmd := exec.Command(command, "--version")
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	err := cmd.Run()
	return err
}

func (msgList *MessageList) ConvertToPDF() {
	var output chan<- *MessageInfo = infoChannel
	for _, info := range msgList.List {
		if len(info.AppMsgExtInfo.Title) > 0 && len(info.AppMsgExtInfo.ContentUrl) > 0 {
			output <- info
		}
	}
}

func (msgList *MessageList) ConcatPDFs(output string) {
	close(infoChannel)
	syncWaitGroup.Wait()
	syncWaitGroup.Add(1)
	close(generatedPdfChannel)
	syncWaitGroup.Wait()

	filenames := pdfFilenames(generatedPdfFilenames)
	filenames.Sort()

	args := []string{"-dBATCH", "-dNOPAUSE", "-q", "-sDEVICE=pdfwrite", "-sOutputFile=" + output}
	for _, pdfFilename := range filenames {
		args = append(args, pdfFilename.filename)
	}
	cmd := exec.Command("gs", args...)
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to execute `gs` command: %s\n", err)
	}
	for _, pdfFilename := range filenames {
		if err := os.Remove(pdfFilename.filename); err != nil {
			log.Printf("Failed to remove %s: %s\n", pdfFilename.filename, err)
		}
	}
}

func pdfGenerator(input <-chan *MessageInfo, output chan<- *pdfFilename) {
	for info := range input {
		var (
			datetime int    = info.CommMsgInfo.Datetime
			url      string = info.AppMsgExtInfo.ContentUrl
			title    string = info.AppMsgExtInfo.Title
			filename string = strconv.Itoa(info.CommMsgInfo.Id) + ".pdf"
		)

		cmd := exec.Command("wkhtmltopdf", "--title", title, url, filename)
		cmd.Stdout = ioutil.Discard
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to execute `wkhtmltopdf %s %s`: %s\n", url, filename, err)
		}
		output <- &pdfFilename{datetime: datetime, filename: filename}
	}
	syncWaitGroup.Done()
}

func generatedPdfCollector(input <-chan *pdfFilename) {
	for filename := range input {
		generatedPdfFilenames = append(generatedPdfFilenames, filename)
	}
	syncWaitGroup.Done()
}

type pdfFilenames []*pdfFilename

func (this pdfFilenames) Len() int {
	return len(this)
}

func (this pdfFilenames) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this pdfFilenames) Less(i, j int) bool {
	return this[i].datetime < this[j].datetime
}

func (this pdfFilenames) Sort() {
	sort.Sort(sort.Reverse(this))
}
