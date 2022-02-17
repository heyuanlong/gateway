package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/selector"
	"golang.org/x/net/http2"

	config "github.com/go-kratos/gateway/api/gateway/config/v1"
)

const _globalClientPool = 10

var _ selector.Node = &node{}
var _globalClient = defaultClient()
var _globalH2Client = defaultH2Client()

func env2int(in string, def int) int {
	v := os.Getenv(in)
	if v == "" {
		return def
	}
	i, _ := strconv.ParseInt(v, 10, 64)
	return int(i)
}

var DefaultTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 600 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func defaultClient() *http.Client {
	tr := DefaultTransport.(*http.Transport).Clone()
	tr.MaxIdleConns = env2int("MAX_IDLE_CONNS", 100)                 //1000
	tr.MaxConnsPerHost = env2int("MAX_CONNS_PER_HOST", 100)          //100
	tr.MaxIdleConnsPerHost = env2int("MAX_IDLE_CONNS_PER_HOST", 100) //100
	tr.DisableCompression = true
	return &http.Client{Transport: tr}
}

func defaultH2Client() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			// So http2.Transport doesn't complain the URL scheme isn't 'https'
			AllowHTTP:          true,
			DisableCompression: true,
			// Pretend we are dialing a TLS endpoint.
			// Note, we ignore the passed tls.Config
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
}

func newNode(addr string, protocol config.Protocol, weight *int64, md map[string]string) *node {
	node := &node{
		protocol: protocol,
		address:  addr,
		weight:   weight,
		metadata: md,
	}
	if protocol == config.Protocol_GRPC {
		node.client = _globalH2Client
	} else {
		node.client = _globalClient
	}
	return node
}

type node struct {
	address  string
	name     string
	weight   *int64
	version  string
	metadata map[string]string

	client   *http.Client
	protocol config.Protocol
}

func (n *node) Address() string {
	return n.address
}

// ServiceName is service name
func (n *node) ServiceName() string {
	return n.name
}

// InitialWeight is the initial value of scheduling weight
// if not set return nil
func (n *node) InitialWeight() *int64 {
	return n.weight
}

// Version is service node version
func (n *node) Version() string {
	return n.version
}

// Metadata is the kv pair metadata associated with the service instance.
// version,namespace,region,protocol etc..
func (n *node) Metadata() map[string]string {
	return n.metadata
}
