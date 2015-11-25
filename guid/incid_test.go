package guid

import (
	"testing"
	"fmt"
	"time"
)

func TestGenerate(t *testing.T) {
	incid, err := NewIncId()
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 19; i++ {
		fmt.Println(incid.Generate())
	}
	time.Sleep(time.Second)
}