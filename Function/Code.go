package Function

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

type Code struct {
	Type      int // 类型
	values    int // 数值
	deviation int // 误差%
}

func Create_Code(PIN string) (barcode.Barcode, error) {
	cs, _ := code128.Encode(PIN)
	return barcode.Scale(cs, 350, 70)
}

func Verify_Code(PIN string) bool {

	return false
}
