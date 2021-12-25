package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ying32/govcl/vcl"
	"golang.org/x/sys/windows"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sca_client/datastruct"
	"strings"
	"syscall"
	"time"
)

func Ips() string {

	ips := make(map[string]string)

	interfaces, _ := net.Interfaces()

	for _, i := range interfaces {
		byName, _ := net.InterfaceByName(i.Name)

		addresses, _ := byName.Addrs()

		for _, v := range addresses {
			ips[byName.Name] = v.String()

			fmt.Println(byName.Name, v.String(), v.Network())
			if strings.HasPrefix(v.String(), "192.168.216.") {
				return strings.TrimSuffix(v.String(), "/24")
			}

		}
	}
	return "127.0.0.1"

}

//TODO: 解决由于权限问题 导致C:\ProgramData\ZeroTier\One\authtoken.secret文件无法读取的问题    已解决！
func getVpnToken() string {
	verb := "runas"
	//exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	//args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString("cmd")
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	//fmt.Println("当前路径：" + cwd)
	var flag = "/c copy C:\\ProgramData\\ZeroTier\\One\\authtoken.secret " + cwd
	argPtr, _ := syscall.UTF16PtrFromString(flag)

	var showCmd int32 = 0 //0表示窗口不显示  1表示执行时窗口正常显示

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
	arg := "type " + cwd + "\\authtoken.secret"
	cmd0 := exec.Command("PowerShell", arg)
	stdout0, err := cmd0.StdoutPipe() // 获取命令输出内容
	if err != nil {
		fmt.Println(err)
	}
	if err := cmd0.Start(); err != nil { //开始执行命令
		fmt.Println(err)
	}
	useBufferIO := false
	var outputBuf0 bytes.Buffer
	var token string
	if !useBufferIO {
		for {
			tempoutput := make([]byte, 256)
			n, err := stdout0.Read(tempoutput)
			if err != nil {
				if err == io.EOF { //读取到内容的最后位置
					break
				} else {
					fmt.Println(err)
				}
			}
			if n > 0 {
				outputBuf0.Write(tempoutput[:n])
			}
		}
		token = strings.Replace(outputBuf0.String(), "\r\n", "", -1)
		//fmt.Println("获取到token:" + token)
	}
	return token
}

//获取盒子的节点ID

func getHeziNodeId() string {
	var vpnstatus datastruct.StatusRespone

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				//本地地址  ipaddr是本地外网IP
				lAddr, err := net.ResolveTCPAddr(netw, Ips()+":0")
				if err != nil {
					return nil, err
				}
				//被请求的地址
				rAddr, err := net.ResolveTCPAddr(netw, addr)
				if err != nil {
					return nil, err
				}
				conn, err := net.DialTCP(netw, lAddr, rAddr)
				if err != nil {
					return nil, err
				}
				deadline := time.Now().Add(35 * time.Second)
				conn.SetDeadline(deadline)
				return conn, nil
			},
		},
	}
	req, err := http.NewRequest("GET", "http://"+Tgconf.DefaultIpAddr+":9993/status", nil)
	if err != nil {
		return ""
	}
	req.Header.Add("X-ZT1-Auth", getHeziVpnToken(Tgconf.DefaultIpAddr))
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(result))
	err = json.Unmarshal(result, &vpnstatus)
	if err != nil {
		fmt.Print("Unmarshalerr=", err)
		return ""
	}

	return vpnstatus.Address
}

//获取本机的节点ID
func getNodeId() string {
	var vpnstatus datastruct.StatusRespone

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:9993/status", nil)
	if err != nil {
		return ""
	}
	req.Header.Add("X-ZT1-Auth", getVpnToken())
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(result))
	err = json.Unmarshal(result, &vpnstatus)
	if err != nil {
		fmt.Print("Unmarshalerr=", err)
		return ""
	}
	fmt.Println("本机address:", vpnstatus.Address)
	return vpnstatus.Address

}

//获取当前vpn连接状态 返回 ：状态信息， networkid  -1表示连接失败 -3表示盒子未连接vpn网络 0表示连接中 1 表示连接成功
func getHeziVPNConnectStatus() (int64, *datastruct.GetJoinedNetworkConfigResponse) {

	var getJoinedNetworkConfigResponse *datastruct.GetJoinedNetworkConfigResponse
	fmt.Println("Tgconf.VPN.NetWorkId", Tgconf.VPN.NetWorkId)
	if Tgconf.VPN.NetWorkId != "" {
		fmt.Println("here****")
		client := &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					//本地地址  ipaddr是本地外网IP
					lAddr, err := net.ResolveTCPAddr(netw, "192.168.216.211:0")
					if err != nil {
						return nil, err
					}
					//被请求的地址
					rAddr, err := net.ResolveTCPAddr(netw, addr)
					if err != nil {
						return nil, err
					}
					conn, err := net.DialTCP(netw, lAddr, rAddr)
					if err != nil {
						return nil, err
					}
					deadline := time.Now().Add(35 * time.Second)
					conn.SetDeadline(deadline)
					return conn, nil
				},
			},
		}
		req, err := http.NewRequest("GET", "http://"+Tgconf.DefaultIpAddr+":9993/network/"+Tgconf.VPN.NetWorkId, nil)
		if err != nil {
			fmt.Print("err=", err)
			return -1, nil
		}

		req.Header.Add("X-ZT1-Auth", getHeziVpnToken(Tgconf.DefaultIpAddr))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print("err=", err)
			return -1, nil
		}
		result, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(result))
		err = json.Unmarshal(result, &getJoinedNetworkConfigResponse)
		//如果盒子没有连过该VPN则转化为结构体时会报错
		if err != nil {
			return -1, nil

		}
		if getJoinedNetworkConfigResponse.Status == "OK" {

			return 1, getJoinedNetworkConfigResponse
		} else if getJoinedNetworkConfigResponse.Status == "REQUESTING_CONFIGURATION" {
			return 0, getJoinedNetworkConfigResponse
		} else {
			return -3, nil
		}
	} else {
		return -2, nil
	}

}

//盒子断开VPN -1:退出失败  1：退出成功
func heziDisconnectVPN(networkid string, heziip string) int {
	localip := Ips()
	fmt.Println(localip)
	//如下构建的http.client是绑定了固定ip来发送数据的
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				//本地地址  ipaddr是本地外网IP
				lAddr, err := net.ResolveTCPAddr(netw, localip+":0")
				if err != nil {
					return nil, err
				}
				//被请求的地址
				rAddr, err := net.ResolveTCPAddr(netw, addr)
				if err != nil {
					return nil, err
				}
				conn, err := net.DialTCP(netw, lAddr, rAddr)
				if err != nil {
					return nil, err
				}
				deadline := time.Now().Add(35 * time.Second)
				conn.SetDeadline(deadline)
				return conn, nil
			},
		},
	}
	req, err := http.NewRequest("DELETE", "http://"+heziip+":9993/network/"+networkid, nil)
	if err != nil {
		return -1
	}
	req.Header.Add("X-ZT1-Auth", getHeziVpnToken(Tgconf.DefaultIpAddr))
	resp, err := client.Do(req)
	if err != nil {
		return -1
	}
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(result))
	return 1
}

//TODO  从server中获取盒子的token
func getHeziVpnToken( heziip string) string {
	var heziVPNTokenResponse *datastruct.HeziVPNTokenResponse
	localip := Ips()
	//如下构建的http.client是绑定了固定ip来发送数据的
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				//本地地址  ipaddr是本地外网IP
				lAddr, err := net.ResolveTCPAddr(netw, localip+":0")
				if err != nil {
					return nil, err
				}
				//被请求的地址
				rAddr, err := net.ResolveTCPAddr(netw, addr)
				if err != nil {
					return nil, err
				}
				conn, err := net.DialTCP(netw, lAddr, rAddr)
				if err != nil {
					return nil, err
				}
				deadline := time.Now().Add(35 * time.Second)
				conn.SetDeadline(deadline)
				return conn, nil
			},
		},
	}
	req, err := http.NewRequest("DELETE", "http://"+heziip+":8000/getVPNToken", nil)
	if err != nil {
		fmt.Print("err=", err)

	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("err=", err)

	}
	result, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(result))
	err = json.Unmarshal(result, &heziVPNTokenResponse)

	return heziVPNTokenResponse.Auth_token
}

//本机断开VPN
func localDisConnectVPN(networkid string) int {
	var flag int

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://127.0.0.1:9993/network/"+networkid, nil)
	if err != nil {
		fmt.Print("err=", err)
	}
	req.Header.Add("X-ZT1-Auth", getVpnToken())

	_, err = client.Do(req)

	if err != nil {
		fmt.Print("err=", err)
	}
	//result, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("本机断开与子网：%s的断开连接成功", networkid)
	flag = 1

	return flag

	////连接成功后：
	////1.修改配置文件VPN部分
	////Tgconf.VPN.NetWorkId = joinNetworkResponse.ID
	//Tgconf.VPN.Connect = false
	////Tgconf.VPN.NetWorkName = joinNetworkResponse.Name
	//fmt.Println(Tgconf)
	//UpdateKey(Tgconf)
	//time.Sleep(1 * time.Second)
	//Tgconf = InitKey()
}

//本机连接VPN
func localConnectVPN(networkid string) (int, datastruct.JoinNetworkResponse) {
	var status int
	var joinNetworkResponse datastruct.JoinNetworkResponse

	var flag bool

	var joinNetworkRequest datastruct.JoinNetworkRequest
	fmt.Println("localConnectVPN本机连接vpn：" + networkid)
	joinNetworkRequest.AllowManaged = true
	joinNetworkRequest.AllowDNS = true
	joinNetworkRequest.AllowGlobal = true
	joinNetworkRequest.AllowDefault = true
	request, err := json.Marshal(joinNetworkRequest)
	if err != nil {
		flag = false
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:9993/network/"+networkid, bytes.NewBuffer(request))
	if err != nil {
		flag = false
	}

	req.Header.Add("X-ZT1-Auth", getVpnToken())

	resp, err := client.Do(req)
	if err != nil {
		flag = false
	}
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(result))
	err = json.Unmarshal(result, &joinNetworkResponse)
	if err != nil {
		flag = false
	}
	//先临时检查一下本机的vpn状态  之后交由自动任务进行实时检查
	if joinNetworkResponse.Status == "OK" || joinNetworkResponse.Status == "REQUESTING_CONFIGURATION" {
		fmt.Println("本地连接vpn返回的状态：", joinNetworkResponse.Status)
		flag = true

		//Tgconf.VPN.NetWorkId = joinNetworkResponse.ID
		//Tgconf.VPN.Connect = true
		//Tgconf.VPN.NetWorkName = joinNetworkResponse.Name
		//UpdateKey(Tgconf)
		//连接成功 刷新子网列表
		//fillNetworks(mainForm)
		//连接成功 刷新成员列表
		//showAllMembers(joinNetworkResponse.ID, lv1)
		//连接成功 刷新分配的ip
		//vcl.ThreadSync(func() {
		//	//mainForm.VPNIPEdit.SetText(connectresponse.AssignedAddresses[0])
		//	mainForm.VPNTest.SetCaption("断开子网")
		//})
	}

	if flag {
		status = 1

	} else {
		status = -1

	}
	return status, joinNetworkResponse
}

//盒子连接VPN
func heziJoinVPN(networkid string, heziip string) int {
	var flag bool
	var joinNetworkResponse datastruct.JoinNetworkResponse
	var joinNetworkRequest datastruct.JoinNetworkRequest

	fmt.Println("***************" + networkid)
	joinNetworkRequest.AllowManaged = true
	joinNetworkRequest.AllowDNS = true
	joinNetworkRequest.AllowGlobal = true
	joinNetworkRequest.AllowDefault = true
	request, err := json.Marshal(joinNetworkRequest)
	if err != nil {
		flag = false
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				//本地地址  ipaddr是本地外网IP
				lAddr, err := net.ResolveTCPAddr(netw, Ips()+":0")
				if err != nil {
					return nil, err
				}
				//被请求的地址
				rAddr, err := net.ResolveTCPAddr(netw, addr)
				if err != nil {
					return nil, err
				}
				conn, err := net.DialTCP(netw, lAddr, rAddr)
				if err != nil {
					return nil, err
				}
				deadline := time.Now().Add(35 * time.Second)
				conn.SetDeadline(deadline)
				return conn, nil
			},
		},
	}
	req, err := http.NewRequest("POST", "http://"+Tgconf.DefaultIpAddr+":9993/network/"+networkid, bytes.NewBuffer(request))
	if err != nil {
		flag = false
	}

	req.Header.Add("X-ZT1-Auth", getHeziVpnToken(Tgconf.DefaultIpAddr))

	resp, err := client.Do(req)
	if err != nil {

		flag = false
	}
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(result))
	err = json.Unmarshal(result, &joinNetworkResponse)
	if err != nil {
		flag = false
	}
	//先临时检查一下盒子的vpn状态  之后交由自动任务进行实时检查
	if joinNetworkResponse.Status == "OK" || joinNetworkResponse.Status == "REQUESTING_CONFIGURATION" {
		flag = true
	}

	if flag {
		return 1
	} else {
		return -1
	}
}

//获取当前vpn连接状态 返回 ：状态信息， networkid
func getVPNConnectStatus() (int64, *datastruct.GetJoinedNetworkConfigResponse) {
	var getJoinedNetworkConfigResponse *datastruct.GetJoinedNetworkConfigResponse
	if Tgconf.VPN.Connect == true && Tgconf.VPN.NetWorkId != "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://127.0.0.1:9993/network/"+Tgconf.VPN.NetWorkId, nil)
		if err != nil {
			fmt.Print("err=", err)
		}
		req.Header.Add("X-ZT1-Auth", getVpnToken())
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print("err=", err)
		}
		result, _ := ioutil.ReadAll(resp.Body)
		//fmt.Print(string(result))
		err = json.Unmarshal(result, &getJoinedNetworkConfigResponse)
		if err != nil {
			fmt.Print("反序列化getJoinedNetworkConfigRequest出错=", err)
		}
		if getJoinedNetworkConfigResponse.Status == "OK" {
			Tgconf.VPN.NetWorkId = getJoinedNetworkConfigResponse.ID
			Tgconf.VPN.Connect = true
			Tgconf.VPN.NetWorkName = getJoinedNetworkConfigResponse.Name
			UpdateKey(Tgconf)
			return 1, getJoinedNetworkConfigResponse
		} else if getJoinedNetworkConfigResponse.Status == "REQUESTING_CONFIGURATION" {
			Tgconf.VPN.NetWorkId = getJoinedNetworkConfigResponse.ID
			Tgconf.VPN.Connect = true
			Tgconf.VPN.NetWorkName = getJoinedNetworkConfigResponse.Name
			UpdateKey(Tgconf)
			return 0, getJoinedNetworkConfigResponse
		} else {
			return -2, getJoinedNetworkConfigResponse
		}

	} else if Tgconf.VPN.Connect == false && Tgconf.VPN.NetWorkId != "" {
		//如果上次关闭软件时候，vpn是关着的，那么用户默认设置的为不连接vpn
		return 2, nil
	} else if Tgconf.VPN.NetWorkId == "" {
		//如果从未连接过VPN
		return -1, nil
	}
	return 3, nil
}

func getAllVPNUsersById(vpnID string) (datastruct.MemberListResponse, error) {
	//TOKEN: "kocvbWu3ByI3SnVZ3LnZGNtMgzTKXr1p"
	var memberList datastruct.MemberListResponse
	client := http.Client{}

	req, err := http.NewRequest("GET", "https://my.zerotier.com/api/v1/network/"+vpnID+"/member", nil)
	if err != nil {
		fmt.Print("err=", err)

	}
	//TODO bearer认证TOKEN是在网页端依账号生成的一个账号对应多个子网 一个账号对应一个TOKEN
	req.Header.Add("Authorization", "Bearer "+"kocvbWu3ByI3SnVZ3LnZGNtMgzTKXr1p")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("err=", err)
		return nil, err
	}
	result, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(result, &memberList)
	if err != nil {
		fmt.Print("Unmarshalerr=", err)
		return nil, err
	}
	fmt.Printf("获取到%s的全部成员信息：", vpnID)
	return memberList, nil
}

func changeLocalAndHeziStatus(connected bool) {
	if connected {

		//设置vpn页面的本机与盒子的连接状态图
		hezinodeId := getHeziNodeId()
		fmt.Println("主机已连接盒子！")
		vcl.ThreadSync(func() { //非主线程访问ui
			pictuer := vcl.NewPicture()
			pictuer.LoadFromFile("./images/linked.jpg")
			mainForm.LocalAndHezi.SetPicture(pictuer)
			if hezinodeId != "" {
				mainForm.HeziNodeIdLabel.SetCaption(hezinodeId)
			} else {
				mainForm.HeziNodeIdLabel.SetCaption("未获取到ID")
			}

		})
	} else {
		//设置vpn页面的本机与盒子的连接状态图

		vcl.ThreadSync(func() { //非主线程访问ui
			pictuer := vcl.NewPicture()
			pictuer.LoadFromFile("./images/unlinked.jpg")
			mainForm.LocalAndHezi.SetPicture(pictuer)
			mainForm.HeziNodeIdLabel.SetCaption("尚未连接到盒子！")
		})
	}
}
func changeHeziAndVPNStatus(connected int) {
	if connected == 1 {
		//设置vpn页面的盒子与VPN的连接状态图

		vcl.ThreadSync(func() { //非主线程访问ui
			pictuer := vcl.NewPicture()
			pictuer.LoadFromFile("./images/linked.jpg")
			mainForm.HeziandVPN.SetPicture(pictuer)
			mainForm.VPNidLabel.SetCaption(Tgconf.VPN.NetWorkId)
		})

	} else if connected == 0 {
		//设置vpn页面的本机与盒子的连接状态图

		vcl.ThreadSync(func() { //非主线程访问ui
			pictuer := vcl.NewPicture()
			pictuer.LoadFromFile("./images/linking.jpg")
			mainForm.HeziandVPN.SetPicture(pictuer)
			mainForm.VPNidLabel.SetCaption(Tgconf.VPN.NetWorkId)

		})
	} else {
		//设置vpn页面的本机与盒子的连接状态图

		vcl.ThreadSync(func() { //非主线程访问ui
			pictuer := vcl.NewPicture()
			pictuer.LoadFromFile("./images/unlinked.jpg")
			mainForm.HeziandVPN.SetPicture(pictuer)
			mainForm.VPNidLabel.SetCaption("盒子未连接到子网！")
		})
	}
}
