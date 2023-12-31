package docker

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/f7ed0/HebergOS-dockerlink/logger"
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

const LASTING_TIME_DEFAULT int64 = 14400 // 4 hours (240 minutes)
const LASTING_TIME_DAY int64 = 86400 // 1 day
const LASTING_TIME_WEEK int64 = 604800 // 1 week

const INSTANT_ONLY = 0
const DAY_ONLY int = 1
const WEEK_ONLY int = 2

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

	KeepForDay	bool
	KeepForWeek	bool
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

func (s StatHolder) Export(container_id string,since int64,scale int) string {
	ret := "{\n"
	stat_holder.RLock()
	for key,val := range s[container_id] {
		if(key > since) {
			t := time.Now().Unix()
			if(scale == INSTANT_ONLY && key < t - LASTING_TIME_DEFAULT) {
				continue
			}
			if((scale == DAY_ONLY) && (!val.KeepForDay || key < t - LASTING_TIME_DAY)) {
				continue
			}
			if((scale == WEEK_ONLY) && (!val.KeepForWeek || key < t - LASTING_TIME_WEEK)) {
				continue
			}
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
	if(len(ret) > 2){
		ret = ret[:len(ret)-2]
	}
	return ret + "\n}"
	
}

func (s *StatHolder) DestroyOlder(timestamp int64, container_id string) {
	stat_holder.Lock()
	_,ok := (*s)[container_id]
	if !ok {
		return
	}
	for key,val :=  range (*s)[container_id] {
		if !val.KeepForDay && key < timestamp - LASTING_TIME_DEFAULT {
			delete((*s)[container_id],key)
		} else if !val.KeepForWeek && key < timestamp - LASTING_TIME_DAY{
			delete((*s)[container_id],key)
		} else if key < timestamp - LASTING_TIME_WEEK {
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
		logger.Default.LogPanic(err.Error())
	}

	lasts := map[string]*Stat{}
	lasts_t := map[string]int64{}
	var t1,t2 time.Time
	u := map[string]any{}
	Go := math.Pow(2,30)
	ko := math.Pow(2,10)
	rounds := 0
	for true {
		rounds = (rounds+1)%80
		t1 = time.Now()
		containers,err := dk.Client.ContainerList(dk.Context,types.ContainerListOptions{
			All : true,
		})

		if err != nil {
			logger.Default.LogPanic(err.Error())
		}

		for _,container := range containers {
			s := new(Stat)

			s.KeepForDay = (rounds%8 == 0)
			s.KeepForWeek = (rounds == 0)

			var now time.Time

			info,err := dk.Client.ContainerInspect(dk.Context,container.ID)
			if err != nil {
				logger.Default.Log("ERROR",err.Error())
				continue
			}

			now = time.Now()

			s.MemUsage = 0
			s.MemLimit = 0
			s.CpuUsage = 0
			s.CpuQuota = 0

			s.NetRx = 0
			s.NetTx = 0

			s.CpuPercent = 0

			s.NetDRx = 0
			s.NetDTx = 0

			if info.State.Running {
				stat,err := dk.Client.ContainerStats(dk.Context,container.ID,true)
				if err == nil {
					data := json.NewDecoder(stat.Body)
					err = data.Decode(&u)
					if err != nil {
						logger.Default.LogPanic(err.Error())
					}
					stat.Body.Close()

					if(u == nil) {
						logger.Default.LogPanic("u is nil")
					}

					now,err =  time.Parse(time.RFC3339Nano,u["read"].(string))
					if err != nil {
						logger.Default.LogPanic(err.Error())
					}

					use,ok := (u["memory_stats"].(map[string]any)["usage"].(float64))
					if !ok {
						continue
					}

					
					s.MemUsage = ( use )/Go
					s.MemLimit = u["memory_stats"].(map[string]any)["limit"].(float64)/Go

					s.CpuUsage = u["cpu_stats"].(map[string]any)["cpu_usage"].(map[string]any)["total_usage"].(float64)
					s.CpuQuota = (float64(info.HostConfig.CPUQuota)/float64(info.HostConfig.CPUPeriod))*100
					if  math.IsNaN(s.CpuQuota) {
						s.CpuQuota = 0
					}

					s.NetRx = u["networks"].(map[string]any)["eth0"].(map[string]any)["rx_bytes"].(float64)/ko
					s.NetTx = u["networks"].(map[string]any)["eth0"].(map[string]any)["tx_bytes"].(float64)/ko
					
					// Calcul du CpuPercentage
					l,ok := lasts[container.ID]
					t,ok2 := lasts_t[container.ID]
					if ok && ok2 {
						s.CpuPercent = ((s.CpuUsage - l.CpuUsage)/float64((now.UnixNano() - t)))*100
						if s.CpuPercent < 0 {
							s.CpuPercent = 0
						}
					}

					// Calculs des NetD
					if ok {
						s.NetDRx = s.NetRx - l.NetRx
						s.NetDTx = s.NetTx - l.NetTx
						if s.NetDRx < 0 {
							s.NetDRx = 0
						}
						if s.NetDTx < 0 {
							s.NetDTx = 0
						}
					}
				}

				
			}			

			Sh.Add(now.Unix(),container.ID,s)

			lasts[container.ID] = s
			lasts_t[container.ID] = now.UnixNano()
		} 
		t2 = time.Now()

		logger.Default.Log("INFO","Statloop took %v seconds.",float64(t2.UnixMilli() - t1.UnixMilli())/1000)
		time.Sleep(15*time.Second - time.Duration(t2.UnixNano() - t1.UnixNano()))
	}
}