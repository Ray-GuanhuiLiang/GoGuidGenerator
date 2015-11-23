package guid

import (
	"time"
	"fmt"
	"sync"
	"errors"
	"net"
	"bytes"
	"crypto/md5"
	"os"
	"strconv"
)

const MaxUint = ^uint32(0)

type Guid struct {
	sync.Mutex
	workId uint32
	tick   uint32
	lastTime uint32
	lastTick uint32
}

/**
 * 只会用到这个workId的前三个字节
 */
func NewGuid(workId uint32) (*Guid, error) {
	newWorkId := workId & 0x00FFFFFF
	if (newWorkId != workId) {
		return nil, errors.New("input workid is too large!!")
	}
	if workId == 0 {
		wid, err := defaultWorkId()
		if err != nil {
			return nil, err
		}
		workId = wid
	}
	return &Guid{workId: workId}, nil
}

func defaultWorkId() (uint32, error) {
	var buf bytes.Buffer
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	for _, inter := range interfaces {
		buf.Write(inter.HardwareAddr)
		buf.WriteByte(byte(0))
    }
	
	//fmt.Println("-------------------")
	
	inter2, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	for _, i2 := range inter2 {
		buf.WriteString(i2.String())
		buf.WriteByte(byte(0))
    }
	
	buf.WriteString(strconv.Itoa(os.Getpid()))
	
	//fmt.Println(buf.String())
	
	bs := md5.Sum(buf.Bytes())
	//fmt.Println(bs)
	
	ret := uint32(bs[0]) << 24 + uint32(bs[1]) << 16 + uint32(bs[2]) << 8 + uint32(bs[0])
	//fmt.Println(ret)
	
	return ret, nil
}

func (this *Guid) Generate() (uint64, error) {
	cur := (uint32)(time.Now().Unix())
	
	// 这个是带锁版本，有机会写个不带锁版本
	this.Lock()
	if cur > this.lastTime {
		this.lastTime = cur
	} else {
		if this.lastTick == MaxUint {
			if this.tick == 0 {
				this.Unlock()
				return 0, errors.New("meet max id count in 1 second")
			}
		} else {
			if this.tick + 1 == this.lastTick {
				this.Unlock()
				return 0, errors.New("meet max id count in 1 second")
			}
		}
	}
	thatTick := this.tick
	if this.tick == MaxUint {
		this.tick = 0
	} else {
		this.tick++
	}
	this.Unlock()
	
	return uint64(cur) << 32 + (uint64)(this.workId) << 16 + uint64(thatTick), nil
}
