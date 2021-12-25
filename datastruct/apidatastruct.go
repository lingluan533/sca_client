package datastruct


type HeziVPNTokenResponse struct{
	Auth_token  string `json:auth_token`
}

//Central API :修改子网中的用户信息 包括（用户名、设备类型等） 返回
type ModifyMemberResponse struct {
	ID           string `json:"id"`
	Clock        int64  `json:"clock"`
	NetworkID    string `json:"networkId"`
	NodeID       string `json:"nodeId"`
	ControllerID string `json:"controllerId"`
	Hidden       bool   `json:"hidden"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Config       struct {
		ActiveBridge         bool     `json:"activeBridge"`
		Authorized           bool     `json:"authorized"`
		Capabilities         []int    `json:"capabilities"`
		CreationTime         int64    `json:"creationTime"`
		ID                   string   `json:"id"`
		Identity             string   `json:"identity"`
		IPAssignments        []string `json:"ipAssignments"`
		LastAuthorizedTime   int64    `json:"lastAuthorizedTime"`
		LastDeauthorizedTime int      `json:"lastDeauthorizedTime"`
		NoAutoAssignIps      bool     `json:"noAutoAssignIps"`
		Revision             int      `json:"revision"`
		Tags                 [][]int  `json:"tags"`
		VMajor               int      `json:"vMajor"`
		VMinor               int      `json:"vMinor"`
		VRev                 int      `json:"vRev"`
		VProto               int      `json:"vProto"`
	} `json:"config"`
	LastOnline          int64  `json:"lastOnline"`
	PhysicalAddress     string `json:"physicalAddress"`
	ClientVersion       string `json:"clientVersion"`
	ProtocolVersion     int    `json:"protocolVersion"`
	SupportsRulesEngine bool   `json:"supportsRulesEngine"`
}

//Central API :修改子网中的用户信息 包括（用户名、设备类型等）
type ModifyMemberRequest struct {
	Hidden      bool   `json:"hidden"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Config      struct {
		ActiveBridge    bool     `json:"activeBridge"`
		Authorized      bool     `json:"authorized"`
		Capabilities    []int    `json:"capabilities"`
		IPAssignments   []string `json:"ipAssignments"`
		NoAutoAssignIps bool     `json:"noAutoAssignIps"`
		Tags            [][]int  `json:"tags"`
	} `json:"config"`
}



// 参数：子网id 认证：bearer token 查询某个子网的所有成员信息，返回
type MemberListResponse []struct {
	ID           string `json:"id"`
	Clock        int64  `json:"clock"`
	NetworkID    string `json:"networkId"`
	NodeID       string `json:"nodeId"`
	ControllerID string `json:"controllerId"`
	Hidden       bool   `json:"hidden"`
	Name         string `json:"name"`
	Online		 bool  `json:"online"`
	Description  string `json:"description"`
	Config       struct {
		ActiveBridge         bool     `json:"activeBridge"`
		Authorized           bool     `json:"authorized"`
		Capabilities         []int    `json:"capabilities"`
		CreationTime         int64    `json:"creationTime"`
		ID                   string   `json:"id"`
		Identity             string   `json:"identity"`
		IPAssignments        []string `json:"ipAssignments"`
		LastAuthorizedTime   int64    `json:"lastAuthorizedTime"`
		LastDeauthorizedTime int      `json:"lastDeauthorizedTime"`
		NoAutoAssignIps      bool     `json:"noAutoAssignIps"`
		Revision             int      `json:"revision"`
		Tags                 [][]int  `json:"tags"`
		VMajor               int      `json:"vMajor"`
		VMinor               int      `json:"vMinor"`
		VRev                 int      `json:"vRev"`
		VProto               int      `json:"vProto"`
	} `json:"config"`
	LastOnline          int64  `json:"lastOnline"`
	PhysicalAddress     string `json:"physicalAddress"`
	ClientVersion       string `json:"clientVersion"`
	ProtocolVersion     int    `json:"protocolVersion"`
	SupportsRulesEngine bool   `json:"supportsRulesEngine"`
}



//查询已加入的网络的配置信息，请求
type GetJoinedNetworkConfigResponse struct {
	AllowDNS          bool     `json:"allowDNS"`
	AllowDefault      bool     `json:"allowDefault"`
	AllowGlobal       bool     `json:"allowGlobal"`
	AllowManaged      bool     `json:"allowManaged"`
	AssignedAddresses []string `json:"assignedAddresses"`
	Bridge            bool     `json:"bridge"`
	BroadcastEnabled  bool     `json:"broadcastEnabled"`
	DNS               struct {
		Domain  string   `json:"domain"`
		Servers []string `json:"servers"`
	} `json:"dns"`
	ID                     string `json:"id"`
	Mac                    string `json:"mac"`
	Mtu                    int    `json:"mtu"`
	MulticastSubscriptions []struct {
		Adi int    `json:"adi"`
		Mac string `json:"mac"`
	} `json:"multicastSubscriptions"`
	Name            string `json:"name"`
	NetconfRevision int    `json:"netconfRevision"`
	PortDeviceName  string `json:"portDeviceName"`
	PortError       int    `json:"portError"`
	Routes          []struct {
		Flags  int    `json:"flags"`
		Metric int    `json:"metric"`
		Target string `json:"target"`
		Via    string `json:"via"`
	} `json:"routes"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

type JoinNetworkResponse struct {
	AllowDNS          bool     `json:"allowDNS"`
	AllowDefault      bool     `json:"allowDefault"`
	AllowGlobal       bool     `json:"allowGlobal"`
	AllowManaged      bool     `json:"allowManaged"`
	AssignedAddresses []string `json:"assignedAddresses"`
	Bridge            bool     `json:"bridge"`
	BroadcastEnabled  bool     `json:"broadcastEnabled"`
	Dhcp              bool     `json:"dhcp"`
	DNS               struct {
		Domain  string   `json:"domain"`
		Servers []string `json:"servers"`
	} `json:"dns"`
	ID                     string `json:"id"`
	Mac                    string `json:"mac"`
	Mtu                    int    `json:"mtu"`
	MulticastSubscriptions []struct {
		Adi int64  `json:"adi"`
		Mac string `json:"mac"`
	} `json:"multicastSubscriptions"`
	Name            string `json:"name"`
	NetconfRevision int    `json:"netconfRevision"`
	Nwid            string `json:"nwid"`
	PortDeviceName  string `json:"portDeviceName"`
	PortError       int    `json:"portError"`
	Routes          []struct {
		Flags  int    `json:"flags"`
		Metric int    `json:"metric"`
		Target string `json:"target"`
		Via    string `json:"via"`
	} `json:"routes"`
	Status string `json:"status"`
	Type   string `json:"type"`
}


//申请连接vpn
type JoinNetworkRequest struct {
	AllowDNS          bool     `json:"allowDNS"`
	AllowDefault      bool     `json:"allowDefault"`
	AllowGlobal       bool     `json:"allowGlobal"`
	AllowManaged      bool     `json:"allowManaged"`
	AssignedAddresses []string `json:"assignedAddresses"`
	Bridge            bool     `json:"bridge"`
	BroadcastEnabled  bool     `json:"broadcastEnabled"`
	DNS               struct {
		Domain  string   `json:"domain"`
		Servers []string `json:"servers"`
	} `json:"dns"`
	ID                     string `json:"id"`
	Mac                    string `json:"mac"`
	Mtu                    int    `json:"mtu"`
	MulticastSubscriptions []struct {
		Adi int    `json:"adi"`
		Mac string `json:"mac"`
	} `json:"multicastSubscriptions"`
	Name            string `json:"name"`
	NetconfRevision int    `json:"netconfRevision"`
	PortDeviceName  string `json:"portDeviceName"`
	PortError       int    `json:"portError"`
	Routes          []struct {
		Flags  int    `json:"flags"`
		Metric int    `json:"metric"`
		Target string `json:"target"`
		Via    string `json:"via"`
	} `json:"routes"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

//查看本机vpn状态的返回
type StatusRespone struct {
	Address string `json:address`
	Clock   int64  `json:clock`
	Config  struct {
		Physical struct{} `json:physical`
		Settings struct {
			AllowTcpFallbackRelay bool   `json:allowTcpFallbackRelay`
			PortMappingEnabled    bool   `json:portMappingEnabled`
			PrimaryPort           int    `json:primaryPort`
			SoftwareUpdate        string `json:softwareUpdate`
			SoftwareUpdateChannel string `json:softwareUpdateChannel`
		} `json:settings`
	} `json:config`
	Online               bool   `json:online`
	PlanetWorldId        int64  `json:planetWorldId`
	PlanetWorldTimestamp int64  `json:planetWorldTimestamp`
	PublicIdentity       string `json:publicIdentity`
	TcpFallbackActive    bool   `json:tcpFallbackActive`
	Version              string `json:version`
	VersionBuild         int    `json:versionBuild`
	VersionMajor         int    `json:versionMajor`
	VersionMinor         int    `json:versionMinor`
	VersionRev           int    `json:versionRev`
}
