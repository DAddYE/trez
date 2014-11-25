# T-REZ

A super fast image resizer build on top of opencv

## Install

You need:

- opencv
- libjpeg-turbo (suggested)

If you are on a mac and you have [brew](http://brew.sh) installed you can use [this
formula](https://gist.githubusercontent.com/DAddYE/bd6a4819ec0bbb2efb0a/raw/opencv.rb)

```
brew install https://gist.githubusercontent.com/DAddYE/bd6a4819ec0bbb2efb0a/raw/opencv.rb
```

## Benchmarks

On my MacBook Pro mid 2012 (Yosemite):

```
$ GOMAXPROCS=8 go run bench/main.go -file testdata/American_Dad.jpg -size 200x200 -workers 8

## Resize speed of 10000

  mean: 22.15215ms
   max: 29.909601ms
   %99: 27.517023ms
stdDev: 1.707429ms
  rate: 361.44 images resized per second
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
