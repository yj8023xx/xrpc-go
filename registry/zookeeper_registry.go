package registry

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"time"
)

const (
	ZKRegistryPath = "/registry"
)

type ZookeeperRegistry struct {
	conn *zk.Conn
}

func (r *ZookeeperRegistry) Connect(uri string) error {
	conn, _, err := zk.Connect([]string{uri}, time.Second)
	if err != nil {
		return err
	}
	r.conn = conn
	return nil
}

func (r *ZookeeperRegistry) RegisterService(serviceName, uri string) error {
	var registryPath = ZKRegistryPath
	isExist, _, err := r.conn.Exists(registryPath)
	if err != nil {
		return err
	}
	// 创建永久节点
	if !isExist {
		r.conn.Create(registryPath, nil, 0, zk.WorldACL(zk.PermAll))
	}
	var servicePath = registryPath + "/" + serviceName
	isExist, _, err = r.conn.Exists(servicePath)
	if err != nil {
		return err
	}
	// 创建永久节点
	if !isExist {
		r.conn.Create(servicePath, nil, 0, zk.WorldACL(zk.PermAll))
	}
	// 创建临时顺序节点
	var addressPath = servicePath + "/"
	r.conn.Create(addressPath, []byte(uri), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	fmt.Println(addressPath)
	fmt.Println(uri)
	return nil
}

func (r *ZookeeperRegistry) GetServiceAddress(serviceName string) ([]string, error) {
	var servicePath = ZKRegistryPath + "/" + serviceName
	addrs, _, err := r.conn.Children(servicePath)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, nil
	}
	var uris []string
	for _, addr := range addrs {
		uri, _, _ := r.conn.Get(servicePath + "/" + addr)
		uris = append(uris, string(uri))
	}
	return uris, nil
}

func (r *ZookeeperRegistry) Name() string {
	return "zookeeper"
}

func NewZookeeperRegistry(uri string) Registry {
	registry := &ZookeeperRegistry{}
	registry.Connect(uri)
	return registry
}
