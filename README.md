# T-REZ

A super fast image resizer build on top of opencv and jpeg-turbo.

--

This package keeps a good quality of images however is built for _speed_.
It will strip out image informations like IPTC, EXIF, ColorSync profile, etc.

The benchmark below is just _illustrative_.
This code is derived form a production version that process (at peek) more than `1000 images/s` on a
single machine.

## Install

You need:

- opencv
- libjpeg (jpeg-turbo suggested)
- libpng (optional if you want to read png files)

### Macintosh

If you are on a mac and you have [brew](http://brew.sh) installed you can use [this
formula](https://gist.githubusercontent.com/DAddYE/bd6a4819ec0bbb2efb0a/raw/opencv.rb)

```
brew install https://gist.githubusercontent.com/DAddYE/bd6a4819ec0bbb2efb0a/raw/opencv.rb
```

### Linux

Install jpeg-turbo from sources:

```
libjpeg-turbo-1.3.1$ ./configure --with-jpeg8 --with-pic && make && sudo make install
```

and then opencv with:

```
opencv-2.4.9$ cmake -DBUILD_JASPER=OFF -DBUILD_JPEG=OFF \
-DJPEG_INCLUDE_DIR=/opt/libjpeg-turbo/include/ \
-DJPEG_LIBRARY=/opt/libjpeg-turbo/lib64/libturbojpeg.a -DBUILD_TESTS=OFF -DBUILD_PERF_TESTS=OFF \
-DBUILD_opencv_java=OFF -DWITH_OPENEXR=ON -DWITH_QT=OFF -DWITH_TBB=OFF -DWITH_GSTREAMER=OFF \
-DWITH_GTK=OFF -DWITH_V4L=OFF -DBUILD_DOCS=OFF -DBUILD_NEW_PYTHON_SUPPORT=OFF  -DWITH_1394=OFF \
-DWITH_OPENCL=OFF -DENABLE_SSSE3=ON -DENABLE_SSE41=ON -DENABLE_SSE42=ON -DENABLE_AVX=ON -Wno-dev \
-DCMAKE_INSTALL_PREFIX=/usr
```

Finally:

```
go get github.com/daddye/trez
```

## Features

It supports currently:

- `fit` resize algo
- `fill` resize algo
- `background` color
- `gravity` in case of `fit`
- `quality` of jpeg (default `95`)

## TODO

- `trim` borders
- `enlarge` option

## Benchmarks

On:

```
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                24
On-line CPU(s) list:   0-23
Thread(s) per core:    2
Core(s) per socket:    6
Socket(s):             2
NUMA node(s):          2
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 45
Stepping:              7
CPU MHz:               1895.270
BogoMIPS:              3790.86
Virtualization:        VT-x
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              15360K
NUMA node0 CPU(s):     0-5,12-17
NUMA node1 CPU(s):     6-11,18-23
```

running 24 threads:

```
$ GOMAXPROCS=24 go run bench/main.go -file testdata/American_Dad.jpg -size 200x200 -workers 24

## Resize speed of 14000 resizes
  mean: 34.377864ms
   min: 20.395457ms
   max: 61.757311ms
   %99: 55.475588ms
stdDev: 4.781191ms
  rate: 700.16 ops (images resized per second)
```

## LICENSE

Copyright (C) 2014 Davide D'Agostino

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the "Software"),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
