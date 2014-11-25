package trez

//#cgo CFLAGS: -Wall -Wextra -Os -Wno-unused-function -Wno-unused-parameter
//#cgo linux  pkg-config: opencv
//#cgo darwin pkg-config: opencv
//
//#include <opencv/highgui.h>
//#include <opencv/cv.h>
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

func Resize(data []byte, options Options) ([]byte, error) {
	if len(data) == 0 {
		return nil, errNoData
	}

	// enable optimizations
	C.cvUseOptimized(1)

	// create a mat
	mat := C.cvCreateMat(1, C.int(len(data)), C.CV_8UC1)
	C.set_data_mat(mat, unsafe.Pointer(&data[0]))

	// start decoding
	src := C.cvDecodeImage(mat, C.CV_LOAD_IMAGE_COLOR)
	C.cvReleaseMat(&mat)

	// check it's a valid source
	if src == nil || src.width == 0 || src.height == 0 {
		return nil, errInvalidSourceFormat
	}
	defer C.cvReleaseImage(&src)

	// set some defaults
	if options.Width == 0 {
		options.Width = int(src.width)
	}

	if options.Height == 0 {
		options.Height = int(src.height)
	}

	// prepare the destination image
	size := C.cvSize(C.int(options.Width), C.int(options.Height))
	dst := C.cvCreateImage(size, src.depth, src.nChannels)
	defer C.cvReleaseImage(&dst)

	// get the x,y factor
	xf := float64(dst.width) / float64(src.width)
	yf := float64(dst.height) / float64(src.height)

	rect := C.CvRect{}

	switch options.Algo {
	case FIT:
		ratio := math.Min(xf, yf)
		rect.width = C.int(math.Floor(float64(src.width) * ratio))
		rect.height = C.int(math.Floor(float64(src.height) * ratio))

		switch options.Gravity {
		case CENTER:
			rect.x = (dst.width - rect.width) / 2
			rect.y = (dst.height - rect.height) / 2
		case NORTH:
			rect.x = (dst.width - rect.width) / 2
			rect.y = 0
		case NORTH_WEST:
			rect.x = 0
			rect.y = 0
		case NORTH_EAST:
			rect.x = (dst.width - rect.width)
			rect.y = 0
		case SOUTH:
			rect.x = (dst.width - rect.width) / 2
			rect.y = (dst.height - rect.height)
		case SOUTH_WEST:
			rect.x = 0
			rect.y = (dst.height - rect.height)
		case SOUTH_EAST:
			rect.x = (dst.width - rect.width)
			rect.y = (dst.height - rect.height)
		case WEST:
			rect.x = 0
			rect.y = (dst.height - rect.height) / 2
		case EAST:
			rect.x = (dst.width - rect.width)
			rect.y = (dst.height - rect.height) / 2
		}

		b, g, r := options.Background[2], options.Background[1], options.Background[0]
		C.cvSet(unsafe.Pointer(dst), C.cvScalar(C.double(b), C.double(g), C.double(r), 0), nil)
		C.cvSetImageROI(dst, rect)
		C.cvResize(unsafe.Pointer(src), unsafe.Pointer(dst), C.CV_INTER_AREA)
		C.cvResetImageROI(dst)
	case FILL:
		ratio := math.Max(xf, yf)
		size := C.cvSize(
			C.int(math.Ceil(float64(src.width)*ratio)),
			C.int(math.Ceil(float64(src.height)*ratio)),
		)
		mid := C.cvCreateImage(size, src.depth, src.nChannels)
		defer C.cvReleaseImage(&mid)

		C.cvResize(unsafe.Pointer(src), unsafe.Pointer(mid), C.CV_INTER_AREA)

		if int(mid.width) > options.Width {
			rect.x = (mid.width - C.int(options.Width)) / 2
		}
		if int(mid.height) > options.Height {
			rect.y = (mid.height - C.int(options.Height)) / 2
		}

		rect.width = dst.width
		rect.height = dst.height

		C.cvSetImageROI(mid, rect)
		dst = (*C.IplImage)(C.cvClone(unsafe.Pointer(mid)))
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

	// encode
	ext := C.CString(".jpg")
	ret := C.cvEncodeImage(ext, unsafe.Pointer(dst), &compression[0])
	C.free(unsafe.Pointer(ext))

	if ret == nil {
		return nil, errEncoding
	}

	ptr := C.ptr_from_mat(ret)
	data = C.GoBytes(unsafe.Pointer(ptr), ret.step)
	C.cvReleaseMat(&ret)

	return data, nil
}

type ratio struct {
	src float64
	max float64
}
