package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTLSClientAuthentication(t *testing.T) {
	serverCert, err := tls.LoadX509KeyPair("../certs/server.pem", "../certs/server.key")
	require.NoError(t, err)

	certpool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("../certs/ca.pem")
	require.NoError(t, err)

	assert.True(t, certpool.AppendCertsFromPEM(pem))

	serverConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certpool,
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client!")
	}))
	ts.Config.TLSConfig = serverConfig
	ts.StartTLS()

	clientCert, err := tls.LoadX509KeyPair("../certs/client.pem", "../certs/client.key")
	require.NoError(t, err)

	clientConfig := tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		InsecureSkipVerify: true,
	}

	client := ts.Client()
	client.Transport = &http.Transport{
		TLSClientConfig: &clientConfig,
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, []byte("Hello, client!\n"), body)
}
