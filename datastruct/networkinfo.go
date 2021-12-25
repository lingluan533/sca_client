package datastruct


type AllNetworks []Network

type Network struct {
	NetworkID  string `yaml:"NetworkID"`
	NetworkName string `yaml:"NetworkName"`
	UserName   string `yaml:"UserName"`
	DeviceType string `yaml:"DeviceType"`
}
