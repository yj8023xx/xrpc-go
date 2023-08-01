package registry

import "net/url"

const (
	ZooKeeperName string = "zookeeper"
	NacosName     string = "nacos"
)

const (
	ZooKeeper = iota
	Nacos
)

type NewRegistryFunc func(string) Registry

var registryMap = make(map[string]NewRegistryFunc)
var nameMap = make(map[int]string)

func init() {
	registryMap[ZooKeeperName] = NewZookeeperRegistry
	nameMap[ZooKeeper] = ZooKeeperName
	registryMap[NacosName] = NewNacosRegistry
	nameMap[Nacos] = NacosName
}

type Registry interface {
	Connect(uri string) error
	RegisterService(serviceName, uri string) error
	GetServiceAddress(serviceName string) ([]string, error)
	Name() string
}

func NewRegistry(nameServiceUri string) (Registry, error) {
	u, err := url.Parse(nameServiceUri)
	if err != nil {
		return nil, err
	}
	return registryMap[u.Scheme](u.Host), nil
}
