package helper

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"sync"
)

var BytesBufferPool = &sync.Pool{New: func() any {
	logrus.Info("create buffer from Pool 📦")
	buf := &bytes.Buffer{}
	return buf
}}
