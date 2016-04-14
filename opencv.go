package trez

//#cgo linux  pkg-config: opencv
//#cgo darwin pkg-config: opencv
//
//#include <opencv/highgui.h>
//#include <opencv/cv.h>
//
//#include "sharpen.h"
//
//uchar* ptr_from_mat(CvMat* mat){
//	return mat->data.ptr;
//}
//
//void set_data_mat(CvMat* mat, void* ptr) {
//	mat->data.ptr = ptr;
//}
import "C"
import (
	"errors"
	"math"
	"unsafe"
)

var (
	errNoData              = errors.New("image data length is zero")
	errInvalidSourceFormat = errors.New("invalid data source format")
	errEncoding            = errors.New("error during encoding")
)

func ResizeFromFile(file string, options Options) ([]byte, error) {
	filePath := C.CString(file)
	defer C.free(unsafe.Pointer(filePath))

	image := C.cvLoadImage(filePath, C.CV_LOAD_IMAGE_COLOR)
	if image == nil || image.width == 0 || image.height == 0 {
		return nil, errInvalidSourceFormat
	}
	return resize(image, options)
}

func Resize(data []byte, options Options) ([]byte, error) {
	if len(data) == 0 {
		return nil, errNoData
	}

	// enable optimizations
	C.cvUseOptimized(1)

	// create a mat
	mat := C.cvCreateMat(1, C.int(len(data)), C.CV_8UC1)
	C.set_data_mat(mat, unsafe.Pointer(&data[0]))

	// Decode the source image
	src := C.cvDecodeImage(mat, C.CV_LOAD_IMAGE_COLOR)
	C.cvReleaseMat(&mat)

	return resize(src, options)
}

func resize(src *C.IplImage, options Options) ([]byte, error) {
	// Validate the source
	if src == nil || src.width == 0 || src.height == 0 {
		return nil, errInvalidSourceFormat
	}
	// Ensure the source will be freed.
	defer C.cvReleaseImage(&src)

	// Ensure options has Width and Height set.
	if options.Width == 0 {
		options.Width = int(src.width)
	}
	if options.Height == 0 {
		options.Height = int(src.height)
	}

	// Get the size of the desired output image
	size := C.cvSize(C.int(options.Width), C.int(options.Height))

	// Get the x and y factors
	xf := float64(size.width) / float64(src.width)
	yf := float64(size.height) / float64(src.height)

	// Pointer to the final destination image.
	var dst *C.IplImage

	switch options.Algo {
	case FIT:
		ratio := math.Min(xf, yf)

		// Determine proper ROI rectangle placement
		rect := C.CvRect{}
		rect.width = C.int(math.Floor(float64(src.width) * ratio))
		rect.height = C.int(math.Floor(float64(src.height) * ratio))
		switch options.Gravity {
		case CENTER:
			rect.x = (size.width - rect.width) / 2
			rect.y = (size.height - rect.height) / 2
		case NORTH:
			rect.x = (size.width - rect.width) / 2
			rect.y = 0
		case NORTH_WEST:
			rect.x = 0
			rect.y = 0
		case NORTH_EAST:
			rect.x = (size.width - rect.width)
			rect.y = 0
		case SOUTH:
			rect.x = (size.width - rect.width) / 2
			rect.y = (size.height - rect.height)
		case SOUTH_WEST:
			rect.x = 0
			rect.y = (size.height - rect.height)
		case SOUTH_EAST:
			rect.x = (size.width - rect.width)
			rect.y = (size.height - rect.height)
		case WEST:
			rect.x = 0
			rect.y = (size.height - rect.height) / 2
		case EAST:
			rect.x = (size.width - rect.width)
			rect.y = (size.height - rect.height) / 2
		}

		// Initialize the output image
		dst = C.cvCreateImage(size, src.depth, src.nChannels)
		defer C.cvReleaseImage(&dst)

		b, g, r := options.Background[2], options.Background[1], options.Background[0]
		C.cvSet(unsafe.Pointer(dst), C.cvScalar(C.double(b), C.double(g), C.double(r), 0), nil)
		C.cvSetImageROI(dst, rect)
		C.cvResize(unsafe.Pointer(src), unsafe.Pointer(dst), C.CV_INTER_AREA)
		C.cvResetImageROI(dst)
	case FILL:
		// Algo: Scale image down keeping aspect ratio
		// constant, and then crop to requested size.
		ratio := math.Max(xf, yf)
		// Create an intermediate image
		intermediateSize := C.cvSize(
			C.int(math.Ceil(float64(src.width)*ratio)),
			C.int(math.Ceil(float64(src.height)*ratio)),
		)
		mid := C.cvCreateImage(intermediateSize, src.depth, src.nChannels)
		defer C.cvReleaseImage(&mid)

		C.cvResize(unsafe.Pointer(src), unsafe.Pointer(mid), C.CV_INTER_AREA)

		// Determine proper ROI rectangle placement
		rect := C.CvRect{}
		rect.width = size.width
		rect.height = size.height
		switch options.Gravity {
		case CENTER:
			rect.x = (mid.width - size.width) / 2
			rect.y = (mid.height - size.height) / 2
		case NORTH:
			rect.x = (mid.width - size.width) / 2
			rect.y = 0
		case NORTH_WEST:
			rect.x = 0
			rect.y = 0
		case NORTH_EAST:
			rect.x = (mid.width - size.width)
			rect.y = 0
		case SOUTH:
			rect.x = (mid.width - size.width) / 2
			rect.y = (mid.height - size.height)
		case SOUTH_WEST:
			rect.x = 0
			rect.y = (mid.height - size.height)
		case SOUTH_EAST:
			rect.x = (mid.width - size.width)
			rect.y = (mid.height - size.height)
		case WEST:
			rect.x = 0
			rect.y = (mid.height - size.height) / 2
		case EAST:
			rect.x = (mid.width - size.width)
			rect.y = (mid.height - size.height) / 2
		}

		C.cvSetImageROI(mid, rect)

		// Create the destination image
		dst = (*C.IplImage)(C.cvClone(unsafe.Pointer(mid)))
		defer C.cvReleaseImage(&dst)
		C.cvResetImageROI(mid)
	}

	// set default compression
	if options.Quality == 0 {
		options.Quality = 95
	}

	compression := [3]C.int{
		C.CV_IMWRITE_JPEG_QUALITY,
		C.int(options.Quality),
		0,
	}

	// Okay, we have our "final" image. Do we need to sharpen it?
	var final *C.IplImage
	if options.SharpenAmount > 0 && options.SharpenRadius > 0 {
		final = C.sharpen(dst, C.int(options.SharpenAmount), C.double(options.SharpenRadius))
	} else {
		final = dst
	}

	// encode
	ext := C.CString(".jpg")
	ret := C.cvEncodeImage(ext, unsafe.Pointer(final), &compression[0])
	C.free(unsafe.Pointer(ext))

	if ret == nil {
		return nil, errEncoding
	}

	ptr := C.ptr_from_mat(ret)
	data := C.GoBytes(unsafe.Pointer(ptr), ret.step)
	C.cvReleaseMat(&ret)

	return data, nil
}

type ratio struct {
	src float64
	max float64
}
