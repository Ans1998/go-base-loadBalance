package config

import (
	"gopkg.in/ini.v1"
	"log"
	"strconv"
)

/**
  main 程序入口
  init  package入口(可多个)
 */
type webConfig struct {
	Url string
	Weight int
}
var ProxyConfigs map[string]webConfig // 定义一个数组map

func init()  {
	ProxyConfigs=make(map[string]webConfig) // 开辟数组map空间（所以每次更新env文件就重启服务器）

	cfg,err:=ini.Load("config/env")
	if err!=nil {
		log.Panicln(err)
	}
	section, err:=cfg.GetSection("proxy")
	if err!=nil {
		log.Panicln(err)
	}
	childSections:=section.ChildSections() // 获取子分区
	for _,val:=range childSections {
		name,_:=val.GetKey("name")
		path,_:=val.GetKey("path")
		weight,_:=val.GetKey("weight")
		weightInt,_:=strconv.Atoi(weight.Value())
		if path!=nil && weight!=nil {
			ProxyConfigs[name.Value()]=webConfig{Url:path.Value(), Weight: weightInt}
		}
	}
}
