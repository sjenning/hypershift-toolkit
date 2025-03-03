package pki

import (
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"

	"github.com/openshift/hypershift-toolkit/pkg/pki/util"
)

type caSpec struct {
	name               string
	commonName         string
	organizationalUnit string
}

type certSpec struct {
	name         string
	ca           string
	commonName   string
	organization string
	hostNames    []string
	ips          []string
}

type kubeconfigSpec struct {
	certSpec
	serverAddress string
}

func generateCAs(caSpecs []caSpec) (map[string]*util.CA, error) {
	result := make(map[string]*util.CA)
	for _, caSpec := range caSpecs {
		ca, err := util.GenerateCA(caSpec.commonName, caSpec.commonName)
		if err != nil {
			return nil, err
		}
		result[caSpec.name] = ca
	}
	return result, nil
}

func generateKubeconfigs(kubeconfigSpecs []kubeconfigSpec, cas map[string]*util.CA) (map[string]*util.Kubeconfig, error) {
	result := make(map[string]*util.Kubeconfig)
	for _, spec := range kubeconfigSpecs {
		ca := cas[spec.ca]
		if ca == nil {
			return nil, fmt.Errorf("CA %s for kubeconfig %s not found", spec.ca, spec.name)
		}
		kubeconfig, err := util.GenerateKubeconfig(spec.serverAddress, spec.commonName, spec.organization, ca)
		if err != nil {
			return nil, err
		}
		result[spec.name] = kubeconfig
	}
	return result, nil
}

func generateCerts(certSpecs []certSpec, cas map[string]*util.CA) (map[string]*util.Cert, error) {
	result := make(map[string]*util.Cert)
	for _, spec := range certSpecs {
		ca := cas[spec.ca]
		if ca == nil {
			return nil, fmt.Errorf("CA %s for cert %s not found", spec.ca, spec.name)
		}
		cert, err := util.GenerateCert(spec.commonName, spec.organization, spec.hostNames, spec.ips, ca)
		if err != nil {
			return nil, err
		}
		result[spec.name] = cert
	}
	return result, nil
}

func ca(name, commonName, organizationalUnit string) caSpec {
	return caSpec{
		name:               name,
		commonName:         commonName,
		organizationalUnit: organizationalUnit,
	}
}

func cert(name, ca, commonName, organization string, hostNames, ips []string) certSpec {
	return certSpec{
		name:         name,
		ca:           ca,
		commonName:   commonName,
		organization: organization,
		hostNames:    hostNames,
		ips:          ips,
	}
}

func kubeconfig(name, serverAddress, ca, commonName, organization string) kubeconfigSpec {
	return kubeconfigSpec{
		certSpec: certSpec{
			name:         name,
			ca:           ca,
			commonName:   commonName,
			organization: organization,
		},
		serverAddress: serverAddress,
	}
}

func writeCerts(certMap map[string]*util.Cert, outputDir string) error {
	for k, v := range certMap {
		if err := v.WriteTo(filepath.Join(outputDir, k), false); err != nil {
			return err
		}
	}
	return nil
}

func writeKubeconfigs(kubeconfigMap map[string]*util.Kubeconfig, outputDir string) error {
	for k, v := range kubeconfigMap {
		if err := v.WriteTo(filepath.Join(outputDir, k)); err != nil {
			return err
		}
	}
	return nil
}

func writeCAs(caMap map[string]*util.CA, outputDir string) error {
	for k, v := range caMap {
		if err := v.WriteTo(filepath.Join(outputDir, k)); err != nil {
			return err
		}
	}
	return nil
}

func writeCombinedCA(cas []string, caMap map[string]*util.CA, outputDir, fileName string) error {
	var caList util.CAList
	for _, c := range cas {
		ca := caMap[c]
		if ca == nil {
			return fmt.Errorf("CA not found: %s", c)
		}
		caList = append(caList, ca)
	}
	if err := caList.WriteTo(filepath.Join(outputDir, fileName)); err != nil {
		return err
	}
	return nil
}

func writeRSAKey(outputDir, name string) error {
	fileName := filepath.Join(outputDir, name+".pem")
	if util.FileExists(fileName) {
		return nil
	}
	key, err := util.PrivateKey()
	if err != nil {
		return err
	}
	b := util.PrivateKeyToPem(key)
	if err := ioutil.WriteFile(fileName, b, 0644); err != nil {
		return err
	}
	return nil
}

func writeDHParams(outputDir, name string) error {
	fileName := filepath.Join(outputDir, name+".pem")
	if util.FileExists(fileName) {
		return nil
	}
	b, err := util.GenerateDHParams()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fileName, b, 0644); err != nil {
		return err
	}
	return nil
}

func nextIP(ip net.IP) net.IP {
	nextIP := net.IP(make([]byte, len(ip)))
	copy(nextIP, ip)
	for j := len(nextIP) - 1; j >= 0; j-- {
		nextIP[j]++
		if nextIP[j] > 0 {
			break
		}
	}
	return nextIP
}

func firstIP(network *net.IPNet) net.IP {
	return nextIP(network.IP)
}
