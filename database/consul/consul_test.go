package consul

var consul *Consul

func init() {
	var err error
	if consul, err = New("192.168.1.6:8500"); err != nil {
		panic(err)
	}
}
