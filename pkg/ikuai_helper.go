package pkg

import (
	"github.com/githgf/ikuai"
	"github.com/githgf/ikuai/action"
	"log"
	"strconv"
	"time"
)

var vlanMap map[string]*action.VlanData

func LoadAll(ik *ikuai.IKuai) {
	if vlanMap == nil {
		vlanMap = make(map[string]*action.VlanData)
	}

	log.Println("开始获取ikuai adsl")

	re, err := ik.ShowLanList()
	if err != nil {
		log.Println("ShowLanList-err ", err)
		return
	}

	for _, wan := range re.Data.SnapshootWan {

		num := 0

		limit := "0,100"
		for true {
			rep, err := ik.ShowWanVlanInterface(wan.Interface, limit)
			if err != nil {
				continue
			}
			num = num + len(rep.Data.VlanData)

			for _, vlan := range rep.Data.VlanData {

				vlanMap[vlan.VlanName] = &vlan
			}

			if num >= rep.Data.VlanTotal {
				break
			}
			limit = strconv.Itoa(num) + ",100"

		}

	}

	log.Println("ikuai adsl 获取完毕")
}

func GetByInfName(name string) *action.VlanData {
	return vlanMap[name]
}

func GetAllInf() map[string]*action.VlanData {
	return vlanMap
}

func StartLoadIkuaiAsync(ik *ikuai.IKuai) {
	for true {
		defer func() {
			if err := recover(); err != nil {
				log.Println("load-ikuai 异常:", err)
				time.Sleep(1 * time.Minute)

				StartLoadIkuaiAsync(ik)
			}
		}()
		LoadAll(ik)
		time.Sleep(1 * time.Minute)

	}
}
