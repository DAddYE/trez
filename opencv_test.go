package trez

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tmpdir = "/tmp/trez"

func cleanup() {
	os.RemoveAll(tmpdir)
	os.MkdirAll(tmpdir, 0775)
	fmt.Printf("images available at %s\n", tmpdir)
}

func dumpImage(filename string, data []byte) {
	err := ioutil.WriteFile(path.Join(tmpdir, filename+".jpg"), data, 0665)
	if err != nil {
		panic(err)
	}
}

func TestResize(t *testing.T) {
	cleanup()
	assert := assert.New(t)

	f, err := os.Open("testdata/American_Dad.jpg")
	require.NoError(t, err)

	src, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	// just encode
	data, err := Resize(src, Options{})
	assert.NoError(err)
	dumpImage("just_image", data)
	im, _, err := image.DecodeConfig(bytes.NewReader(data))
	assert.NoError(err)
	assert.Equal(1024, im.Width)
	assert.Equal(768, im.Height)
	size := len(data)

	// just quality
	data, err = Resize(src, Options{Quality: 50})
	assert.NoError(err)
	dumpImage("just_quality", data)
	im, _, err = image.DecodeConfig(bytes.NewReader(data))
	assert.NoError(err)
	assert.Equal(1024, im.Width)
	assert.Equal(768, im.Height)
	assert.True(size > len(data))

	opts := []Options{
		{Width: 200},
		{Height: 200},
		{Width: 200, Height: 200},
		{Width: 150, Height: 120},
		{Width: 2000, Height: 1020},
		{Width: 2000, Height: 3000},
	}

	for _, opt := range opts {
		for _, algo := range []Algo{FIT, FILL} {
			opt.Algo = algo
			opt.Gravity = WEST
			opt.Background = [3]int{255, 255, 255}
			data, err := Resize(src, opt)
			assert.NoError(err)
			dumpImage(fmt.Sprintf("%dx%d_%s", opt.Width, opt.Height, opt.Algo), data)
			_, _, err = image.DecodeConfig(bytes.NewReader(data))
			assert.NoError(err)
		}
	}
}
