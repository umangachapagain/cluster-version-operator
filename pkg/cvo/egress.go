package cvo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"

	"k8s.io/apimachinery/pkg/api/errors"
)

// getHTTPSProxyURL returns a url.URL object for the configured
// https proxy only. It can be nil if does not exist or there is an error.
func (optr *Operator) getHTTPSProxyURL() (*url.URL, string, error) {
	proxy, err := optr.proxyLister.Get("cluster")

	if errors.IsNotFound(err) {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}

	if &proxy.Spec != nil {
		if proxy.Spec.HTTPSProxy != "" {
			proxyURL, err := url.Parse(proxy.Spec.HTTPSProxy)
			if err != nil {
				return nil, "", err
			}
			return proxyURL, proxy.Spec.TrustedCA.Name, nil
		}
	}
	return nil, "", nil
}

func (optr *Operator) getTLSConfig(cmNameRef string) (*tls.Config, error) {
	cm, err := optr.cmConfigLister.Get(cmNameRef)

	if err != nil {
		return nil, err
	}

	certPool, _ := x509.SystemCertPool()
	if certPool == nil {
		certPool = x509.NewCertPool()
	}

	if cm.Data["ca-bundle.crt"] != "" {
		if ok := certPool.AppendCertsFromPEM([]byte(cm.Data["ca-bundle.crt"])); !ok {
			return nil, fmt.Errorf("unable to add ca-bundle.crt certificates")
		}
	} else {
		return nil, nil
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	return config, nil
}
