package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"log"
)

var one_data = `
DefaultTemplate: "rclone.conf.1box"
DefaultBox: "box1"
DefaultIpAddr: "192.168.216.197"
DefaultBrain: "box1"
AllBox:
  box1:
    Type: "webdav"
    Vendor: "other"
    Url: "http://192.168.216.197:800"
    User: "root"
    Password: "cGFzc3dvcmQ="
    VirtualDisk: "Z:"
DefaultID:
  BoxName: "box1"
  ID: "root"
  Token: ""
DefaultInterval: 30
DefaultLogFile: "sca_client.log"
DefaultAutoBoot: true
`

type VPNConfig struct {
	NetWorkId   string `yaml:NetWorkId`
	Connect     bool   `yaml:Connect`
	NetWorkName string `yaml:NetWorkName`
}
type IDConfig struct {
	BoxName string `yaml:"BoxName"`
	ID      string `yaml:"ID"`
	Token   string `yaml:"Token"`
}
type AllBox struct {
	Type        string `yaml:"Type"`
	Vendor      string `yaml:"Vendor"`
	Url         string `yaml:"Url"`
	User        string `yaml:"User"`
	Password    string `yaml:"Password"`
	VirtualDisk string `yaml:"VirtualDisk"`
}
type GlobalConfig struct {
	DefaultTemplate  string            `yaml:"DefaultTemplate"`
	DefaultBox       string            `yaml:"DefaultBox"`
	DefaultIpAddr    string            `yaml:"DefaultIpAddr"`
	DefaultBrain     string            `yaml:"DefaultBrain"`
	AllBox           map[string]AllBox `yaml:"AllBox"`
	DefaultID        IDConfig          `yaml:"DefaultID"`
	DefaultInterval  int               `yaml:"DefaultInterval"`
	DefaultLogFile   string            `yaml:"DefaultLogFile"`
	DefaultAutoBoot  bool              `yaml:"DefaultAutoBoot"`
	DefaultAutoMount bool              `yaml:"DefaultAutoMount"`
	VPN              VPNConfig         `yaml:"VPN"`
}

func InitKey() GlobalConfig {

	data, err := ioutil.ReadFile(Rundir + "\\config.yaml")
	if err != nil {
		data = []byte(one_data)
		ioutil.WriteFile(Rundir+"\\config.yaml", []byte(data), 0777)
	}
	//var config backend.RedisConfig
	var config GlobalConfig

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("unmarshal 1error: %v", err)
	}

	//log.Println(config)
	return config
}

func UpdateKey(config GlobalConfig) {
	//转换成yaml字符串类型
	d, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("marshal error: %v", err)
	}
	//log.Printf("--- config dump:\n%s\n\n", string(d))
	one_data = string(d)

	ioutil.WriteFile(Rundir+"\\config.yaml", []byte(d), 0777)

}
