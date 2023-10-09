package docker

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
)

const STAT_STRING string = `  "%v" : {
	"memory" : {
	  "used" : %v,
	  "limit" : %v
	},
	"cpu" : {
	  "usage_percent" : %v,
	  "limit" : %v
	},
	"net" : {
	  "up" : %v,
	  "down" : %v,
	  "delta_up": %v,
	  "delta_down": %v
	}
  }`

const LASTING_TIME int64 = 18000 // 5 hours (300 minutes)



type StatHolder map[string]map[int64]*Stat

type Stat struct {
	MemUsage 	float64
	MemLimit 	float64

	CpuPercent 	float64
	CpuUsage 	float64
	CpuQuota	float64

	NetRx 		float64
	NetTx 		float64
	NetDRx 		float64
	NetDTx 		float64
}

var stat_holder *sync.RWMutex = new(sync.RWMutex)
var Sh StatHolder = StatHolder{}

func (s *StatHolder) Add(timestamp int64, container_id string, new *Stat) {
	stat_holder.Lock()
	_,ok := (*s)[container_id]
	if !ok {
		(*s)[container_id] = map[int64]*Stat{}
	}
	(*s)[container_id][timestamp] = new
	stat_holder.Unlock()
	s.DestroyOlder(timestamp,container_id)	
}

func (s StatHolder) Export(container_id string,since int64) string {
	ret := "{\n"
	stat_holder.RLock()
	for key,val := range s[container_id] {
		if(key > since) {
			ret += fmt.Sprintf(
				STAT_STRING,
				key,
				strconv.FormatFloat(val.MemUsage,'f',6,64),
				strconv.FormatFloat(val.MemLimit,'f',6,64),
				strconv.FormatFloat(val.CpuPercent,'f',3,64),
				strconv.FormatFloat(val.CpuQuota,'f',3,64),
				strconv.FormatFloat(val.NetTx,'f',0,64),
				strconv.FormatFloat(val.NetRx,'f',0,64),
				strconv.FormatFloat(val.NetDTx,'f',0,64),
				strconv.FormatFloat(val.NetDRx,'f',0,64),
			) + ",\n"
		}
	}
	stat_holder.RUnlock()
	return ret[:len(ret)-2] + "\n}"
	
}

func (s *StatHolder) DestroyOlder(timestamp int64, container_id string) {
	stat_holder.Lock()
	_,ok := (*s)[container_id]
	if !ok {
		return
	}
	for key :=  range (*s)[container_id] {
			if key < timestamp - LASTING_TIME {
			delete((*s)[container_id],key)
		}
	}
	stat_holder.Unlock()
}

func (s *StatHolder) Wipe(container_id string) {
	stat_holder.Lock()
	delete((*s),container_id)
	stat_holder.Unlock()
}

func FetchStat() {
	dk,err := NewDockerHandler()
	if err != nil {
		log.Default().Panic(err)
	}

	lasts := map[string]*Stat{}
	lasts_t := map[string]int64{}
	var t1,t2 time.Time
	u := map[string]any{}
	Go := math.Pow(2,30)
	ko := math.Pow(2,10)
	for true {
		t1 = time.Now()
		containers,err := dk.Client.ContainerList(dk.Context,types.ContainerListOptions{})

		if err != nil {
			log.Default().Panic(err)
		}

		for _,container := range containers {
			stat,err := dk.Client.ContainerStats(dk.Context,container.ID,false)
			if err != nil {
				log.Default().Println(err.Error())
				continue
			}

			info,err := dk.Client.ContainerInspect(dk.Context,container.ID)
			if err != nil {
				log.Default().Println(err.Error())
				continue
			}


			
			data := json.NewDecoder(stat.Body)
			err = data.Decode(&u)
			if err != nil {
				log.Default().Panic(err)
			}
			stat.Body.Close()

			time,err :=  time.Parse(time.RFC3339Nano,u["read"].(string))
			if err != nil {
				log.Default().Panic(err)
			}

			//log.Default().Println(u["networks"])
			if(u == nil) {
				continue
			}

			use,ok := (u["memory_stats"].(map[string]any)["usage"].(float64))
			if !ok {
				continue
			}

			s := new(Stat)
			s.MemUsage = ( use )/Go
			s.MemLimit = u["memory_stats"].(map[string]any)["limit"].(float64)/Go
			s.CpuUsage = u["cpu_stats"].(map[string]any)["cpu_usage"].(map[string]any)["total_usage"].(float64)
			s.CpuQuota = (float64(info.HostConfig.CPUQuota)/float64(info.HostConfig.CPUPeriod))*100

			s.NetRx = u["networks"].(map[string]any)["eth0"].(map[string]any)["rx_bytes"].(float64)/ko
			s.NetTx = u["networks"].(map[string]any)["eth0"].(map[string]any)["tx_bytes"].(float64)/ko
			
			// Calcul du CpuPercentage
			l,ok := lasts[container.ID]
			t,ok2 := lasts_t[container.ID]
			if ok && ok2 {
				s.CpuPercent = ((s.CpuUsage - l.CpuUsage)/float64((time.UnixNano() - t)))*100
			}

			// Calculs des NetD
			if ok {
				s.NetDRx = s.NetRx - l.NetRx
				s.NetDTx = s.NetTx - l.NetTx
			}

			Sh.Add(time.Unix(),container.ID,s)

			lasts[container.ID] = s
			lasts_t[container.ID] = time.UnixNano()
		} 
		t2 = time.Now()

		//log.Default().Println(float64(t2.UnixMilli() - t1.UnixMilli())/1000)
		time.Sleep(10*time.Second - time.Duration(t2.UnixNano() - t1.UnixNano()))
	}
}