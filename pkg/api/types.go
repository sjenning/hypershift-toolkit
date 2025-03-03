package api

type ClusterParams struct {
	Namespace               string `json:"namespace"`
	ExternalAPIDNSName      string `json:"externalAPIDNSName"`
	ExternalAPIPort         uint   `json:"externalAPIPort"`
	ExternalAPIIPAddress    string `json:"externalAPIAddress"`
	ExternalOpenVPNDNSName  string `json:"externalVPNDNSName"`
	ExternalOpenVPNPort     uint   `json:"externalVPNPort"`
	ServiceCIDR             string `json:"serviceCIDR"`
	PodCIDR                 string `json:"podCIDR"`
	ReleaseImage            string `json:"releaseImage"`
	APINodePort             uint   `json:"apiNodePort"`
	OAuthRoute              string `json:"oauthRoute"`
	IngressSubdomain        string `json:"ingressSubdomain"`
	OpenShiftAPIClusterIP   string `json:"openshiftAPIClusterIP"`
	ImageRegistryHTTPSecret string `json:"imageRegistryHTTPSecret"`
	RouterNodePortHTTP      string `json:"routerNodePortHTTP"`
	RouterNodePortHTTPS     string `json:"routerNodePortHTTPS"`
	OpenVPNNodePort         string `json:"openVPNNodePort"`
	BaseDomain              string `json:"baseDomain"`
}
