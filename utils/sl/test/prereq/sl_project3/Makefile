# Copyright (c) 2016 Intel Corporation
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
all: libsl_project3.so

libsl_project3.so: sl_project3.h sl_project3.c
	$(CC) -shared sl_project3.c -fPIC -o libsl_project3.so \
	-Wl,-rpath=$(SL_PROJECT2)/lib -L$(SL_PROJECT2)/lib \
	-I$(SL_PROJECT2)/include -lsl_project2

install: libsl_project3.so
	install -D sl_project3.h $(PREFIX)/include/sl_project3.h
	install -D libsl_project3.so $(PREFIX)/lib/libsl_project3.so

clean:
	rm -f libsl_project3.so
