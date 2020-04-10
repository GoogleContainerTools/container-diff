/*
Copyright 2020 Google, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"crypto/tls"
	"crypto/x509"
	. "github.com/google/go-containerregistry/pkg/name"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var tlsConfiguration = struct {
	certifiedRegistries     map[string]string
	skipTLSVerifyRegistries map[string]struct{}
}{
	certifiedRegistries:     make(map[string]string),
	skipTLSVerifyRegistries: make(map[string]struct{}),
}

func ConfigureTLS(skipTsVerifyRegistries []string, registriesToCertificates map[string]string) {
	tlsConfiguration.skipTLSVerifyRegistries = make(map[string]struct{})
	for _, registry := range skipTsVerifyRegistries {
		tlsConfiguration.skipTLSVerifyRegistries[registry] = struct{}{}
	}
	tlsConfiguration.certifiedRegistries = make(map[string]string)
	for registry := range registriesToCertificates {
		tlsConfiguration.certifiedRegistries[registry] = registriesToCertificates[registry]
	}
}

func BuildTransport(registry Registry) http.RoundTripper {
	var tr http.RoundTripper = http.DefaultTransport.(*http.Transport).Clone()

	if _, present := tlsConfiguration.skipTLSVerifyRegistries[registry.RegistryStr()]; present {
		tr.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else if certificatePath := tlsConfiguration.certifiedRegistries[registry.RegistryStr()]; certificatePath != "" {
		systemCertPool := defaultX509Handler()
		if err := appendCertificate(systemCertPool, certificatePath); err != nil {
			logrus.WithError(err).Warnf("Failed to load certificate %s for %s\n", certificatePath, registry.RegistryStr())
		} else {
			tr.(*http.Transport).TLSClientConfig = &tls.Config{
				RootCAs: systemCertPool,
			}
		}
	}
	return tr
}

func appendCertificate(pool *x509.CertPool, path string) error {
	pem, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	pool.AppendCertsFromPEM(pem)
	return nil
}

func defaultX509Handler() *x509.CertPool {
	systemCertPool, err := x509.SystemCertPool()
	if err != nil {
		logrus.Warn("Failed to load system cert pool. Loading empty one instead.")
		systemCertPool = x509.NewCertPool()
	}
	return systemCertPool
}
