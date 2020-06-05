package util

import (
	"fmt"
	"proxy/config"
	"sort"
	"time"
)

type HttpServerSort []*HttpServer
func (p HttpServerSort) Len() int           { return len(p) }
func (p HttpServerSort) Less(i, j int) bool { return p[i].CWeight > p[j].CWeight } // 从大到小排序
func (p HttpServerSort) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
// 目标server
type HttpServer struct {
	Host string
	Weight int
	CWeight int // 当前权重
	FailWeight int //一旦失败后，降低的权重
	Status string // 状态, 默认UP，宕机DOWN
	FailCount int // 计数器，默认0
	SuccessCount int // 成功计数器
}
func NewHttpServer(host string, weight int) *HttpServer {
	return &HttpServer{Host:host,Weight:weight, CWeight:0, Status:"UP"} // 默认CWeight为0
}

// 均衡负载
type LoadBalance struct {
	Servers HttpServerSort
	CurIndex int // 指向当前访问的服务器index
	DownCount int // Down的节点数
}
func NewLoadBalance() *LoadBalance  {
	return &LoadBalance{Servers:make([]*HttpServer, 0)}
}
func (this *LoadBalance) AddServer(server *HttpServer)  {
	this.Servers=append(this.Servers, server)
}
//即时计算总权重
func(this *LoadBalance) getSumWeight() int  {
	sum:=0
	for _,server:=range this.Servers{
		realWeight:=server.Weight-server.FailWeight
		if realWeight>0 {
			sum=sum+realWeight
		}
	}
	return sum
}

var LoadBalanceInit *LoadBalance
func init() {
	LoadBalanceInit=NewLoadBalance()
	fmt.Println(config.ProxyConfigs)
	// 循环配置文件
	for _,value:=range config.ProxyConfigs {
		LoadBalanceInit.AddServer(NewHttpServer(value.Url, value.Weight))
	}
	go(func() {
		checkServers(LoadBalanceInit.Servers)
	})()
}

// 平滑加权轮询
func (this *LoadBalance) RoundRobinByWeight() *HttpServer {
	for _,s:=range this.Servers {
		s.CWeight=s.CWeight+(s.Weight-s.FailWeight)
	}
	sort.Sort(this.Servers) // 根据权重排序
	maxServers:=this.Servers[0] // 返回最大 作为命中服务
	maxServers.CWeight=maxServers.CWeight-this.getSumWeight() // 当前权重-权重总数
	return maxServers
}

// 定时检查服务器健康
func checkServers(servers HttpServerSort)  {
	// 每三秒执行一次
	t:=time.NewTicker(time.Second*3)
	check:=NewHttpChecker(servers)
	for {
		select {
		case <- t.C:
			check.Check(time.Second*2)
			for _, s:=range servers {
				fmt.Println("host-", s.Host, "Status-",s.Status, "FailCount-",s.FailCount, "SuccessCount-",s.SuccessCount, "Weight-", s.Weight, "CWeight-", s.CWeight,"FailWeight-", s.FailWeight)
			}
			fmt.Println("----------------------健康检查分割线--------------------------")
		}
	}
}