#include <opencv2/core/core.hpp>
#include <opencv2/highgui/highgui.hpp>
#include <opencv2/imgproc/imgproc.hpp>

extern "C" {
IplImage* sharpen(IplImage* img, int sharpenAmount, double radius) {
    // Convert IplImage* -> Mat
    cv::Mat raw(img);
    
    // Now, sharpen the Mat.
    cv::Mat * sharpened = new cv::Mat;
    GaussianBlur(raw, *sharpened, cv::Size(0, 0), radius);
    addWeighted(raw, 1.0 + (sharpenAmount/100.0), *sharpened, -(sharpenAmount/100.0), 0, *sharpened);
    IplImage * image = new IplImage;
    *image = static_cast<IplImage>(*sharpened); 
    return image;
}
} // extern "C"