package util

import (
	"fmt"
	"testing"
)

func TestSplitAndMatchName(t *testing.T) {
	f := SplitAndMatch("IMAX,杜比,imax,巨幕", " IMA激光厅")

	fmt.Println(f)
}
