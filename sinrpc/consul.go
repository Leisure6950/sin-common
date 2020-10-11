package sinrpc

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sin-z/sin-common/sinlog"
	"github.com/sin-z/sin-common/sinprocess"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	_consulAddr = "loacalhost:8500"
)

type ServiceInfo struct {
	ServiceID string
	IP        string
	Port      int
	Load      int
	Timestamp int //load updated ts
}

type KVData struct {
	Load      int `json:"load"`
	Timestamp int `json:"ts"`
}

func CheckErr(err error) {
	if err != nil {
		sinlog.Log().Debugf("error: %v", err)
		os.Exit(1)
	}
}

func doRegistService(serviceName string, ip string, port int) {
	_serviceName = serviceName + "-" + ip
	var tags []string

	service := &api.AgentServiceRegistration{
		ID:      _serviceName,
		Name:    serviceName,
		Port:    port,
		Address: ip,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			TCP:      ip + ":" + strconv.Itoa(port),
			Interval: "5s",
			Timeout:  "200ms",
			DeregisterCriticalServiceAfter: "60s",
		},
	}

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		sinlog.Log().Fatal(err)
	}
	_consulClient = client
	if err := _consulClient.Agent().ServiceRegister(service); err != nil {
		sinlog.Log().Fatal(err)
	}
	sinlog.Log().Debugf("Registered service %q in consul with tags %q", serviceName, strings.Join(tags, ","))
	autoUnRegistService(service.ID)
}

func autoUnRegistService(serviceID string) {
	sinprocess.ExitCallback.Register(func(err error) {
		if _consulClient == nil {
			return
		}
		if err := _consulClient.Agent().ServiceDeregister(serviceID); err != nil {
			sinlog.Log().Fatal(err)
		}
	})
}
func DoDiscover(consul_addr string, found_service string) {
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t.C:
			DiscoverServices(consul_addr, true, found_service)
		}
	}
}

func DiscoverServices(addr string, healthyOnly bool, service_name string) {
	consulConf := api.DefaultConfig()
	consulConf.Address = addr
	client, err := api.NewClient(consulConf)
	CheckErr(err)

	services, _, err := client.Catalog().Services(&api.QueryOptions{})
	CheckErr(err)

	fmt.Println("--do discover ---:", addr)

	var sers []ServiceInfo
	for name := range services {
		servicesData, _, err := client.Health().Service(name, "", healthyOnly,
			&api.QueryOptions{})
		CheckErr(err)
		for _, entry := range servicesData {
			if service_name != entry.Service.Service {
				continue
			}
			for _, health := range entry.Checks {
				if health.ServiceName != service_name {
					continue
				}
				fmt.Println("  health nodeid:", health.Node, " service_name:", health.ServiceName, " service_id:", health.ServiceID, " status:", health.Status, " ip:", entry.Service.Address, " port:", entry.Service.Port)

				var node ServiceInfo
				node.IP = entry.Service.Address
				node.Port = entry.Service.Port
				node.ServiceID = health.ServiceID

				//get data from kv store
				s := GetKeyValue(service_name, node.IP, node.Port)
				if len(s) > 0 {
					var data KVData
					err = json.Unmarshal([]byte(s), &data)
					if err == nil {
						node.Load = data.Load
						node.Timestamp = data.Timestamp
					}
				}
				fmt.Println("service node updated ip:", node.IP, " port:", node.Port, " serviceid:", node.ServiceID, " load:", node.Load, " ts:", node.Timestamp)
				sers = append(sers, node)
			}
		}
	}

	service_locker.Lock()
	servics_map[service_name] = sers
	service_locker.Unlock()
}

func DoUpdateKeyValue(consul_addr string, service_name string, ip string, port int) {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			StoreKeyValue(consul_addr, service_name, ip, port)
		}
	}
}

func StoreKeyValue(consul_addr string, service_name string, ip string, port int) {

	my_kv_key = my_service_name + "/" + ip + ":" + strconv.Itoa(port)

	var data KVData
	data.Load = rand.Intn(100)
	data.Timestamp = int(time.Now().Unix())
	bys, _ := json.Marshal(&data)

	kv := &api.KVPair{
		Key:   my_kv_key,
		Flags: 0,
		Value: bys,
	}

	_, err := consul_client.KV().Put(kv, nil)
	CheckErr(err)
	fmt.Println(" store data key:", kv.Key, " value:", string(bys))
}

func GetKeyValue(service_name string, ip string, port int) string {
	key := service_name + "/" + ip + ":" + strconv.Itoa(port)

	kv, _, err := consul_client.KV().Get(key, nil)
	if kv == nil {
		return ""
	}
	CheckErr(err)

	return string(kv.Value)
}
