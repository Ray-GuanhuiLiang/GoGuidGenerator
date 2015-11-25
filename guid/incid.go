package guid

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	_ "fmt"
	"strconv"
	"sync"
	"log"
)

var SECTION uint64 = 10
var ZKPATH = "/guid/incid"

type IncId struct {
	conn *zk.Conn
	startIndex uint64
	currentIndex uint64
	mu sync.Mutex
}

func NewIncId() (*IncId, error) {
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		return nil, err
	}
	data, _, err := conn.Get(ZKPATH)
	if err != nil {
		conn.Close()
		return nil, err
	}
	str := string(data)
	startIndex, err := strconv.ParseUint(str, 10, 0)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &IncId{conn, startIndex, startIndex, sync.Mutex{}}, nil
}

func (this *IncId) Generate() (uint64, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.currentIndex++
	if this.currentIndex - this.startIndex >= SECTION {
		s := strconv.FormatUint(this.currentIndex, 10)
		_, err := this.conn.Set(ZKPATH, []byte(s), -1)
		if err != nil {
			this.currentIndex--
			return 0, err
		}
		this.startIndex = this.currentIndex
	}
	return this.currentIndex, nil
}
