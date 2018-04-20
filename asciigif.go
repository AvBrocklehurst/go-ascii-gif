package asciigif

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nfnt/resize"
)

const (
	defaultWidth int = 80
)

//http://paulbourke.net/dataformats/asciiart/
var asciiChars = []byte{'$', '@', 'B', '%', '8', '&', 'W', 'M', '#', '*', 'o', 'a', 'h', 'k', 'b', 'd', 'p', 'q', 'w', 'm', 'Z', 'O', '0', 'Q', 'L', 'C', 'J', 'U', 'Y', 'X', 'z', 'c', 'v', 'u', 'n', 'x', 'r', 'j', 'f', 't', '/', '\\', '|', '(', ')', '1', '{', '}', '[', ']', '?', '-', '_', '+', '~', '<', '>', 'i', '!', 'l', 'I', ';', ':', ',', '"', '^', '`', '\'', '.'}

//ASCIIGif represents an ascii gif
type ASCIIGif struct {
	height int
	width  int

	index  int
	images [][]byte
	runner chan bool
}

//New creates a new ASCIIGif. It takes a string containing the path to the gif
func New(path string) (ag *ASCIIGif, err error) {
	ag = new(ASCIIGif)
	ag.width = defaultWidth
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	err = ag.splitGIF(f)
	if err != nil {
		return
	}
	ag.runner = make(chan bool, 1)
	return
}

//NewFromURL creates a new ASCIIGif. It takes a string containing the URL of the gif
func NewFromURL(url string) (ag *ASCIIGif, err error) {
	ag = new(ASCIIGif)
	ag.width = defaultWidth
	res, err := http.Get(url)
	if err != nil {
		return
	}
	err = ag.splitGIF(res.Body)
	if err != nil {
		return
	}
	ag.runner = make(chan bool, 1)
	return
}

func (ag *ASCIIGif) splitGIF(r io.Reader) (err error) {
	gif, err := gif.DecodeAll(r)
	if err != nil {
		err = ErrInvalidGif
		return
	}
	ag.getGifHeight(gif)
	frame := image.NewGray(image.Rect(0, 0, gif.Config.Width, gif.Config.Height))
	draw.Draw(frame, frame.Bounds(), gif.Image[0], image.ZP, draw.Src)
	var ascii []byte
	for _, img := range gif.Image {
		draw.Draw(frame, frame.Bounds(), img, image.ZP, draw.Over)
		ascii, err = ag.asciifyFrame(frame)
		if err != nil {
			err = fmt.Errorf("Error asciifying frame: %v", err)
			return
		}
		ag.images = append(ag.images, ascii)
	}
	return
}

func (ag *ASCIIGif) next() (image []byte) {
	image = ag.images[ag.index]
	if ag.index == len(ag.images)-1 {
		ag.index = 0
		return
	}
	ag.index++
	return
}

//Start starts the ascii gif and prints it to stdout
func (ag *ASCIIGif) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ag.runner:
				ticker.Stop()
				break
			case <-ticker.C:
				//Go to top left of terminal
				fmt.Print("\033[H\033[2J")
				fmt.Println(string(ag.next()))
			}
		}
	}()
}

//Stop stops the asciigif after the next complete print
func (ag *ASCIIGif) Stop() {
	ag.runner <- true
}

func (ag *ASCIIGif) getGifHeight(gif *gif.GIF) {
	ag.getImageHeight(gif.Image[0])
}

func (ag *ASCIIGif) getImageHeight(img image.Image) {
	bounds := img.Bounds()
	ag.height = (10 * bounds.Max.Y * ag.width) / (bounds.Max.X * 16)
	return
}

func (ag *ASCIIGif) resizeImage(old image.Image) (img image.Image) {
	img = resize.Resize(uint(ag.width), uint(ag.height), old, resize.Lanczos3)
	return
}

//asciifyFrame is a function to turn an individual frame of a gif into a slice
//of bytes containing the ascii representation.
func (ag *ASCIIGif) asciifyFrame(img image.Image) (ascii []byte, err error) {
	img = ag.resizeImage(img)
	buf := new(bytes.Buffer)
	for y := 0; y < ag.height; y++ {
		for x := 0; x < ag.width; x++ {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()
			//see https://en.wikipedia.org/wiki/Grayscale#Luma_coding_in_video_systems
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			index := int((lum / 255) * 68 / 255)
			err = buf.WriteByte(asciiChars[index])
			if err != nil {
				return
			}
		}
		err = buf.WriteByte('\n')
		if err != nil {
			return
		}
	}
	ascii = buf.Bytes()
	return
}
