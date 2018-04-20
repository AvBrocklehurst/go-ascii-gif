package asciigif

import (
	"bytes"
	"image"
	"image/gif"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid new test",
			args: args{
				path: "example/alice.gif",
			},
			wantErr: false,
		},
		{
			name: "Invalid new test, file missing",
			args: args{
				path: "example/alicea.gif",
			},
			wantErr: true,
		},
		{
			name: "Invalid new test, not gif",
			args: args{
				path: "example/main.go",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAg, err := New(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: New() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			_ = gotAg
		})
	}
}

func TestNewFromURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid newFromURL test",
			args: args{
				url: "https://upload.wikimedia.org/wikipedia/commons/thumb/2/2c/Rotating_earth_%28large%29.gif/200px-Rotating_earth_%28large%29.gif",
			},
			wantErr: false,
		},
		{
			name: "Invalid newFromURL test, file missing",
			args: args{
				url: "https://not.found.io.pls",
			},
			wantErr: true,
		},
		{
			name: "Invalid newFromURL test, not gif",
			args: args{
				url: "https://google.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAg, err := NewFromURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: NewFromURL() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			_ = gotAg

		})
	}
}

func TestASCIIGif_splitGIF(t *testing.T) {
	bytes.NewReader([]byte{'1', '2'})
	type fields struct {
		height int
		width  int
		index  int
		images [][]byte
		runner chan bool
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Invalid Gif Split (random bytes)",
			fields: fields{},
			args: args{
				r: bytes.NewReader([]byte{'$', '@', 'B', '%', '8', '&', 'W', 'M', '#', '*', 'o', 'a', 'h', 'k', 'b', 'd', 'p', 'q', 'w', 'm', 'Z', 'O', '0', 'Q', 'L', 'C', 'J', 'U', 'Y', 'X', 'z', 'c', 'v', 'u', 'n', 'x', 'r', 'j', 'f', 't', '/', '\\', '|', '(', ')', '1', '{', '}', '[', ']', '?', '-', '_', '+', '~', '<', '>', 'i', '!', 'l', 'I', ';', ':', ',', '"', '^', '`', '\'', '.'}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &ASCIIGif{
				height: tt.fields.height,
				width:  tt.fields.width,
				index:  tt.fields.index,
				images: tt.fields.images,
				runner: tt.fields.runner,
			}
			if err := ag.splitGIF(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("ASCIIGif.splitGIF() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestASCIIGif_next(t *testing.T) {
	type fields struct {
		height int
		width  int
		index  int
		images [][]byte
		runner chan bool
	}
	tests := []struct {
		name      string
		fields    fields
		wantImage []byte
		wantIndex int
	}{
		{
			name: "Next image no loop test",
			fields: fields{
				index:  0,
				images: [][]byte{{'0', '0', '0'}, {'0', '0', '1'}, {'0', '1', '0'}, {'0', '1', '1'}},
			},
			wantImage: []byte{'0', '0', '0'},
			wantIndex: 1,
		},
		{
			name: "Next image loop test",
			fields: fields{
				index:  3,
				images: [][]byte{{'0', '0', '0'}, {'0', '0', '1'}, {'0', '1', '0'}, {'0', '1', '1'}},
			},
			wantImage: []byte{'0', '1', '1'},
			wantIndex: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &ASCIIGif{
				height: tt.fields.height,
				width:  tt.fields.width,
				index:  tt.fields.index,
				images: tt.fields.images,
				runner: tt.fields.runner,
			}
			gotImage := ag.next()
			if !reflect.DeepEqual(gotImage, tt.wantImage) {
				t.Errorf("ASCIIGif.next() = %v, want %v", gotImage, tt.wantImage)
			}
			if ag.index != tt.wantIndex {
				t.Errorf("ASCIIGif.next() index = %d, want %d", ag.index, tt.wantIndex)
			}
		})
	}
}

func TestASCIIGif_Start_Stop(t *testing.T) {
	gif, err := New("example/alice.gif")
	if err != nil {
		t.Fatalf("Error creating gif for start/stop test: %v", err)
	}
	tests := []struct {
		name string
		gif  *ASCIIGif
	}{
		{
			name: "Valid start/stop test",
			gif:  gif,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := tt.gif
			ag.Start()
			time.Sleep(time.Millisecond * 250)
			ag.Stop()
		})
	}
}

func TestASCIIGif_Stop(t *testing.T) {
	type fields struct {
		height int
		width  int
		index  int
		images [][]byte
		runner chan bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &ASCIIGif{
				height: tt.fields.height,
				width:  tt.fields.width,
				index:  tt.fields.index,
				images: tt.fields.images,
				runner: tt.fields.runner,
			}
			ag.Stop()
		})
	}
}

func TestASCIIGif_getGifHeight(t *testing.T) {
	type fields struct {
		height int
		width  int
		index  int
		images [][]byte
		runner chan bool
	}
	type args struct {
		gif *gif.GIF
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &ASCIIGif{
				height: tt.fields.height,
				width:  tt.fields.width,
				index:  tt.fields.index,
				images: tt.fields.images,
				runner: tt.fields.runner,
			}
			ag.getGifHeight(tt.args.gif)
		})
	}
}

func TestASCIIGif_getImageHeight(t *testing.T) {
	type fields struct {
		height int
		width  int
		index  int
		images [][]byte
		runner chan bool
	}
	type args struct {
		img image.Image
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &ASCIIGif{
				height: tt.fields.height,
				width:  tt.fields.width,
				index:  tt.fields.index,
				images: tt.fields.images,
				runner: tt.fields.runner,
			}
			ag.getImageHeight(tt.args.img)
		})
	}
}

func TestASCIIGif_resizeImage(t *testing.T) {
	type fields struct {
		height int
		width  int
		index  int
		images [][]byte
		runner chan bool
	}
	type args struct {
		old image.Image
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantImg image.Image
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &ASCIIGif{
				height: tt.fields.height,
				width:  tt.fields.width,
				index:  tt.fields.index,
				images: tt.fields.images,
				runner: tt.fields.runner,
			}
			if gotImg := ag.resizeImage(tt.args.old); !reflect.DeepEqual(gotImg, tt.wantImg) {
				t.Errorf("ASCIIGif.resizeImage() = %v, want %v", gotImg, tt.wantImg)
			}
		})
	}
}

func TestASCIIGif_asciifyFrame(t *testing.T) {
	type fields struct {
		height int
		width  int
		index  int
		images [][]byte
		runner chan bool
	}
	type args struct {
		img image.Image
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantASCII []byte
		wantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &ASCIIGif{
				height: tt.fields.height,
				width:  tt.fields.width,
				index:  tt.fields.index,
				images: tt.fields.images,
				runner: tt.fields.runner,
			}
			gotASCII, err := ag.asciifyFrame(tt.args.img)
			if (err != nil) != tt.wantErr {
				t.Errorf("ASCIIGif.asciifyFrame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotASCII, tt.wantASCII) {
				t.Errorf("ASCIIGif.asciifyFrame() = %v, want %v", gotASCII, tt.wantASCII)
			}
		})
	}
}
