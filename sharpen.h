#ifndef CPPLIB_H
    #define CPPLIB_H
    #ifdef __cplusplus
    extern "C" {
    #endif
    #include <opencv/cv.h>
    IplImage sharpen(IplImage* img, int sharpenAmount, double radius);
    #ifdef __cplusplus
    }
    #endif
#endif