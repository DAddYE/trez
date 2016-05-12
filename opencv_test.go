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

// Setup the temp directory for a new test run
func cleanup() {
	os.RemoveAll(tmpdir)
	os.MkdirAll(tmpdir, 0775)
	fmt.Printf("images available at %s\n", tmpdir)
}

// Write jpeg encoded byte array to disk
func dumpImage(filename string, data []byte) {
	err := ioutil.WriteFile(path.Join(tmpdir, filename+".jpg"), data, 0665)
	if err != nil {
		panic(err)
	}
}

func imageTestHelper(t *testing.T, src_path string, output_prefix string, opts Options) {
	assert := assert.New(t)

	// Read the input image and make sure there are no errors
	f, err := os.Open(src_path)
	require.NoError(t, err)
	src, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	// Perform the resize operation and make sure there are no errors
	// Data will be jpeg encoded
	data, err := Resize(src, opts)

	dest_name := fmt.Sprintf("%s%dx%d_%s", output_prefix, opts.Width, opts.Height, opts.Algo)

	// Write the resized image back to disk
	dumpImage(dest_name, data)

	// Read back the resized image
	im, _, err := image.DecodeConfig(bytes.NewReader(data))

	// Make sure the resized image looks good
	assert.NoError(err)
	assert.Equal(opts.Width, im.Width)
	assert.Equal(opts.Height, im.Height)
}

func TestGeneric(t *testing.T) {
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
			dumpImage(fmt.Sprintf("American_Dad_to_%dx%d_%s", opt.Width, opt.Height, opt.Algo), data)
			_, _, err = image.DecodeConfig(bytes.NewReader(data))
			assert.NoError(err)
		}
	}
}

func Test100x100(t *testing.T) {

	src_path := "testdata/100x100_square.png"
	output_prefix := "100x100_square_to_"

	opts := Options{
		Algo:    FILL,
		Gravity: WEST,
		Width:   50,
		Height:  50,
	}

	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a tall portion at the center of the image
func Test100x200TallCenter(t *testing.T) {

	src_path := "testdata/100x200_tall_center.png"
	output_prefix := "100x200_tall_center_to_"

	// Just a crop, take the center
	opts := Options{
		Algo:    FILL,
		Gravity: CENTER,
		Width:   50,
		Height:  200,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the center
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH,
		Width:   25,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the center
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH,
		Width:   100,
		Height:  400,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a tall portion at the left of the image
func Test100x200TallLeft(t *testing.T) {

	src_path := "testdata/100x200_tall_left.png"
	output_prefix := "100x200_tall_left_to_"

	// Just a crop, take the left
	opts := Options{
		Algo:    FILL,
		Gravity: WEST,
		Width:   50,
		Height:  200,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the left
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH_WEST,
		Width:   25,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the left
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH_WEST,
		Width:   100,
		Height:  400,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a tall portion at the right of the image
func Test100x200TallRight(t *testing.T) {

	src_path := "testdata/100x200_tall_right.png"
	output_prefix := "100x200_tall_right_to_"

	// Just a crop, take the right
	opts := Options{
		Algo:    FILL,
		Gravity: EAST,
		Width:   50,
		Height:  200,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the right
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH_EAST,
		Width:   25,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the right
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH_EAST,
		Width:   100,
		Height:  400,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a wide portion at the bottom of the image
func Test100x200WideBottom(t *testing.T) {

	src_path := "testdata/100x200_wide_bottom.png"
	output_prefix := "100x200_wide_bottom_to_"

	// Just a crop, take the bottom
	opts := Options{
		Algo:    FILL,
		Gravity: SOUTH,
		Width:   100,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the bottom
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH_EAST,
		Width:   50,
		Height:  25,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the bottom
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH_WEST,
		Width:   200,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a wide portion at the bottom of the image
func Test100x200WideCenter(t *testing.T) {

	src_path := "testdata/100x200_wide_center.png"
	output_prefix := "100x200_wide_center_to_"

	// Just a crop, take the bottom
	opts := Options{
		Algo:    FILL,
		Gravity: CENTER,
		Width:   100,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the bottom
	opts = Options{
		Algo:    FILL,
		Gravity: EAST,
		Width:   50,
		Height:  25,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the bottom
	opts = Options{
		Algo:    FILL,
		Gravity: WEST,
		Width:   200,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a wide portion at the top of the image
func Test100x200WideTop(t *testing.T) {

	src_path := "testdata/100x200_wide_top.png"
	output_prefix := "100x200_wide_top_to_"

	// Just a crop, take the top
	opts := Options{
		Algo:    FILL,
		Gravity: NORTH,
		Width:   100,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the top
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH_EAST,
		Width:   50,
		Height:  25,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the top
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH_WEST,
		Width:   200,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a wide source image and we are always cropping a tall portion at the center of the image
func Test200x100TallCenter(t *testing.T) {

	src_path := "testdata/200x100_tall_center.png"
	output_prefix := "200x100_tall_center_to_"

	// Just a crop, take the center
	opts := Options{
		Algo:    FILL,
		Gravity: CENTER,
		Width:   50,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the center
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH,
		Width:   25,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the center
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH,
		Width:   100,
		Height:  200,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a tall portion at the left of the image
func Test200x100TallLeft(t *testing.T) {

	src_path := "testdata/200x100_tall_left.png"
	output_prefix := "200x100_tall_left_to_"

	// Just a crop, take the left
	opts := Options{
		Algo:    FILL,
		Gravity: WEST,
		Width:   50,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the left
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH_WEST,
		Width:   25,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the left
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH_WEST,
		Width:   100,
		Height:  200,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a tall portion at the right of the image
func Test200x100TallRight(t *testing.T) {

	src_path := "testdata/200x100_tall_right.png"
	output_prefix := "200x100_tall_right_to_"

	// Just a crop, take the right
	opts := Options{
		Algo:    FILL,
		Gravity: EAST,
		Width:   50,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the right
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH_EAST,
		Width:   25,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the right
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH_EAST,
		Width:   100,
		Height:  200,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a wide portion at the bottom of the image
func Test200x100WideBottom(t *testing.T) {

	src_path := "testdata/200x100_wide_bottom.png"
	output_prefix := "200x100_wide_bottom_to_"

	// Just a crop, take the bottom
	opts := Options{
		Algo:    FILL,
		Gravity: SOUTH,
		Width:   200,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the bottom
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH_EAST,
		Width:   100,
		Height:  25,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the bottom
	opts = Options{
		Algo:    FILL,
		Gravity: SOUTH_WEST,
		Width:   400,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a wide portion at the bottom of the image
func Test200x100WideCenter(t *testing.T) {

	src_path := "testdata/200x100_wide_center.png"
	output_prefix := "200x100_wide_center_to_"

	// Just a crop, take the bottom
	opts := Options{
		Algo:    FILL,
		Gravity: CENTER,
		Width:   200,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the bottom
	opts = Options{
		Algo:    FILL,
		Gravity: EAST,
		Width:   100,
		Height:  25,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the bottom
	opts = Options{
		Algo:    FILL,
		Gravity: WEST,
		Width:   400,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

// Here we have a tall source image and we are always cropping a wide portion at the top of the image
func Test200x100WideTop(t *testing.T) {

	src_path := "testdata/200x100_wide_top.png"
	output_prefix := "200x100_wide_top_to_"

	// Just a crop, take the top
	opts := Options{
		Algo:    FILL,
		Gravity: NORTH,
		Width:   200,
		Height:  50,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Shrink, take the top
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH_EAST,
		Width:   100,
		Height:  25,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

	// Enlarge, take the top
	opts = Options{
		Algo:    FILL,
		Gravity: NORTH_WEST,
		Width:   400,
		Height:  100,
	}
	imageTestHelper(t, src_path, output_prefix, opts)

}

func TestResizeFromFile(t *testing.T) {
	opts := Options{
		Algo:          FILL,
		Gravity:       NORTH_WEST,
		Width:         300,
		Height:        200,
		SharpenAmount: 100,
		SharpenRadius: 0.5,
	}

	srcPath := "testdata/200x100_wide_top.png"

	_, err := ResizeFromFile(srcPath, opts)
	assert.NoError(t, err)

	_, err = ResizeFromFile("NotAFile", opts)
	assert.Error(t, err)
}
