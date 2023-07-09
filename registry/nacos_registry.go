package registry

const (
	NacosRegistryPath = "/registry"
)

type NacosRegistry struct {
}

func (r *NacosRegistry) Connect(uri string) error {
	return nil
}

func (r *NacosRegistry) RegisterService(serviceName, uri string) error {
	return nil
}

func (r *NacosRegistry) GetServiceAddress(serviceName string) ([]string, error) {
	var uris []string
	return uris, nil
}

func (r *NacosRegistry) Name() string {
	return "nacos"
}

func NewNacosRegistry(uri string) Registry {
	registry := &NacosRegistry{}
	registry.Connect(uri)
	return registry
}
