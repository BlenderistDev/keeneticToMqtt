package keeneticdto

type DeviceInfoResponse struct {
	Mac           string                 `json:"mac"`
	Via           string                 `json:"via"`
	IP            string                 `json:"ip"`
	Hostname      string                 `json:"hostname"`
	Name          string                 `json:"name"`
	Interface     DeviceInfoInterface    `json:"interface"`
	Registered    bool                   `json:"registered"`
	Access        string                 `json:"access"`
	Schedule      string                 `json:"schedule"`
	Priority      int                    `json:"priority"`
	Active        bool                   `json:"active"`
	RxBytes       int64                  `json:"rxbytes"`
	TxBytes       int64                  `json:"txbytes"`
	FirstSeen     int64                  `json:"first-seen"`
	LastSeen      int64                  `json:"last-seen"`
	Link          string                 `json:"link"`
	SSID          string                 `json:"ssid"`
	AP            string                 `json:"ap"`
	PSM           bool                   `json:"psm"`
	Authenticated bool                   `json:"authenticated"`
	TxRate        int                    `json:"txrate"`
	Uptime        int64                  `json:"uptime"`
	HT            int                    `json:"ht"`
	Mode          string                 `json:"mode"`
	GI            int                    `json:"gi"`
	RSSI          int                    `json:"rssi"`
	MCS           int                    `json:"mcs"`
	TxSS          int                    `json:"txss"`
	EBF           bool                   `json:"ebf"`
	Security      string                 `json:"security"`
	TrafficShape  DeviceInfoTrafficShape `json:"traffic-shape"`
}

type DeviceInfoInterface struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeviceInfoTrafficShape struct {
	RX       int64  `json:"rx"`
	TX       int64  `json:"tx"`
	Mode     string `json:"mode"`
	Schedule string `json:"schedule"`
}

type DevicePolicy struct {
	Mac      string  `json:"mac"`
	Access   string  `json:"access"`
	Policy   *string `json:"policy"`
	Permit   bool    `json:"permit"`
	Priority int     `json:"priority"`
}
