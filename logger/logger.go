/*
 * Copyright (c) 2021 X-Net Services GmbH
 * Info: https://x-net.at
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"x-net.at/idp/helpers"
)

var Log *logrus.Logger
var file os.File
var fileEnabled bool = false

func Init() {
	fmt.Println("logger: ", helpers.Config.Logger)
	Log = logrus.New()
	if helpers.Config.Logger == "file" {
		file, err := os.OpenFile("/var/log/idp/x-idp.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic("Fatal: Could not open Logfile: " + err.Error())
		}
		Log.SetOutput(file)
		fileEnabled = true
	} else {
		Log.SetOutput(os.Stdout)
	}

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	Log.SetFormatter(formatter)

	if helpers.Config.Debug {
		Log.SetLevel(logrus.WarnLevel)
	} else {
		Log.SetLevel(logrus.DebugLevel)
	}
}

func Destroy() {
	if fileEnabled {
		file.Close()
	}
}
