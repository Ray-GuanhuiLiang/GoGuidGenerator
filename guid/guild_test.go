package guid

import (
	"testing"
	"fmt"
	"time"
)

func TestTick(t *testing.T) {
	g, err := NewGuid()
	if err != nil {
		t.Error(err)
	}
	mockMaxGenerate(g, t)
	//1秒之后再来
	time.Sleep(time.Second)
	g.Generate()
	time.Sleep(time.Second)
	
	mockMaxGenerate(g, t)
}

func mockMaxGenerate(g *Guid, t *testing.T) {
	for i := uint16(0); i < MAX_UINT; i++ {
		_, err := g.Generate()
		if err != nil {
			t.Error(err)
		}
	}
	for j := 0; j < 3; j++ {
		_, err := g.Generate()
		if err == nil {
			t.Error("It must error here")
		} else {
			fmt.Printf("Expcet this error: %s\n", err)
		}
	}
}
