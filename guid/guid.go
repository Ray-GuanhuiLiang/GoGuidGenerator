package guid

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const MAX_UINT uint16 = 0xFFFF

//const MAX_UINT uint16 = 3

type Guid struct {
	lock     sync.Mutex
	workId   uint16
	tick     uint16
	lastTime uint32
	lastTick uint16
}

/**
 * 只会用到这个workId的前三个字节
 */
func NewGuid() (*Guid, error) {
	workId, err := defaultWorkId()
	if err != nil {
		return nil, err
	}
	return &Guid{workId: workId}, nil
}

func defaultWorkId() (uint16, error) {
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

	bs := md5.Sum(buf.Bytes())
	//fmt.Println(bs)

	//挑选16个字节的md5只的第1和第9个
	ret := uint16(bs[0])<<8 + uint16(bs[8])
	//fmt.Println(ret)

	return ret, nil
}

// GUID = TimeStamp(32bit) + workId(16bit) + IncNo(16bit)
func (this *Guid) Generate() (uint64, error) {
	cur := (uint32)(time.Now().Unix())

	this.lock.Lock()
	defer this.lock.Unlock()

	// TODO: 如果修改了系统时间，时间提前，会导致不可生成。
	if cur > this.lastTime {
		this.lastTime = cur
		this.lastTick = this.tick
	} else {
		if this.lastTick == 0 {
			if this.tick == MAX_UINT {
				return 0, errors.New("meet max id count in 1 second")
			}
		} else if this.tick+1 == this.lastTick {
			return 0, errors.New("meet max id count in 1 second")
		}
	}
	thatTick := this.tick
	//fmt.Printf("thatTick=%d, this.tick=%d, this.lastTick=%d\n", thatTick, this.tick, this.lastTick)
	if this.tick == MAX_UINT {
		this.tick = 0
	} else {
		this.tick++
	}
	//fmt.Printf("cur=%d, this.lastTime=%d, this.lastTick=%d, this.tick=%d\n", cur, this.lastTime, this.lastTick, this.tick)

	return uint64(cur)<<32 + (uint64)(this.workId)<<16 + uint64(thatTick), nil
}
