package util

import (
	"github.com/zhaojunlucky/golib/pkg/text"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func DetectFileEnc(txtFile string) encoding.Encoding {
	var encDetect = text.NewBytesEncodingDetect()
	enc, err := encDetect.DetectFileEncoding(txtFile)
	if err != nil {
		return encoding.Nop
	}
	// promote to GB18030
	if enc == simplifiedchinese.HZGB2312 {
		return simplifiedchinese.GB18030
	}
	return enc
}
