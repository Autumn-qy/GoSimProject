package server

type APIClient struct {
	AccessKeySecret string
	AccessKeyId     string
	InterfaceAddr   string
}

func (api APIClient) FetchData() []map[string]interface{} {
	return nil
}

func (api APIClient) WriteBack(resp string) {
	//TODO implement me
	panic("implement me")
}
