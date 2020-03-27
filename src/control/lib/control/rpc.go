//
// (C) Copyright 2020 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//
package control

import (
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/daos-stack/daos/src/control/common/proto/convert"
	"github.com/daos-stack/daos/src/control/lib/hostlist"
	"github.com/daos-stack/daos/src/control/logging"
	"github.com/daos-stack/daos/src/control/security"
)

var defaultLogger = logging.NewCombinedLogger("", ioutil.Discard)

type (
	unaryRPC func(context.Context, grpc.ClientConnInterface) (proto.Message, error)

	targetSetter interface {
		SetHostList(hl []string)
	}

	targetChooser interface {
		getHostList() []string
		isMSRequest() bool
	}

	unaryRPCGetter interface {
		getRPC() unaryRPC
	}

	request struct {
		HostList []string
	}

	msRequest struct{}

	unaryRequest struct {
		request
		rpc unaryRPC
	}

	// UnaryInvoker defines an interface to be implemented by clients
	// capable of invoking a unary RPC (1 response for 1 request).
	UnaryInvoker interface {
		debugLogger
		InvokeUnaryRPC(ctx context.Context, req UnaryRequest) (*UnaryResponse, error)
		InvokeUnaryRPCAsync(ctx context.Context, req UnaryRequest) (HostResponseChan, error)
	}

	// Invoker defines an interface to be implemented by clients
	// capable of invoking unary or stream RPCs.
	Invoker interface {
		UnaryInvoker
		//StreamInvoker
		SetClientConfig(*ClientConfig)
	}

	debugLogger interface {
		Debug(string)
		Debugf(string, ...interface{})
	}

	Client struct {
		config *ClientConfig
		log    debugLogger
	}

	ClientOption func(c *Client)
)

func (r *request) getHostList() []string {
	return r.HostList
}

func (r *request) SetHostList(hl []string) {
	r.HostList = hl
}

func (r *request) isMSRequest() bool {
	return false
}

func (r *unaryRequest) getRPC() unaryRPC {
	return r.rpc
}

func (r *unaryRequest) setRPC(rpc unaryRPC) {
	r.rpc = rpc
}

func (r *msRequest) isMSRequest() bool {
	return true
}

func WithClientLogger(log debugLogger) ClientOption {
	return func(c *Client) {
		c.log = log
	}
}

func WithClientConfig(cfg *ClientConfig) ClientOption {
	return func(c *Client) {
		c.config = cfg
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		config: DefaultClientConfig(),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.log == nil {
		WithClientLogger(defaultLogger)(c)
	}

	return c
}

func DefaultClient() *Client {
	return NewClient(WithClientLogger(defaultLogger))
}

func (c *Client) SetClientConfig(cfg *ClientConfig) {
	c.config = cfg
}

func (c *Client) Debug(msg string) {
	c.log.Debug(msg)
}

func (c *Client) Debugf(fmtStr string, args ...interface{}) {
	c.log.Debugf(fmtStr, args...)
}

// getRequestHosts returns a list of control plane addresses for
// the request. The logic for determining the list is as follows:
//
// 1.  If the request has an explicit hostlist set, use it regardless
//     of request type. This allows for flexibility and config overrides,
//     but generally should be used for non-AP requests to specific hosts.
// 2.  If there is no hostlist set on the request, check the request type
//     and use the following decision tree:
// 2a. If the request is destined for an Access Point (DAOS MS Replica),
//     pick the first host in the configuration's hostlist. By convention
//     this host will be an AP and the request should succeed. If for some
//     reason the request fails, a future mechanism will attempt to find
//     a working AP replica to service the request.
// 2b. If the request is not destined for an AP, then the request is sent
//     to the entire hostlist set in the configuration.
//
// Will always return at least 1 host, or an error.
func getRequestHosts(cfg *ClientConfig, req targetChooser) ([]string, error) {
	hosts := make([]string, len(req.getHostList()))
	copy(hosts, req.getHostList())

	if len(hosts) == 0 {
		if len(cfg.HostList) == 0 {
			return nil, FaultConfigEmptyHostList
		}

		hosts = make([]string, len(cfg.HostList))
		copy(hosts, cfg.HostList)
		if req.isMSRequest() {
			// pick first host as AP, by convention
			hosts = hosts[:1]
		}
	}

	for i, host := range hosts {
		if !strings.Contains(host, ":") {
			if cfg.ControlPort == 0 {
				return nil, FaultConfigBadControlPort
			}
			hosts[i] = fmt.Sprintf("%s:%d", host, cfg.ControlPort)
		}
	}

	return hosts, nil
}

type (
	// UnaryRequest defines an interface to be implemented by
	// unary request types (1 response to 1 request).
	UnaryRequest interface {
		targetChooser
		targetSetter
		unaryRPCGetter
	}

	// HostResponse contains a single host's response to an unary RPC, or
	// an error if the host was unable to respond successfully.
	HostResponse struct {
		Addr    string
		Error   error
		Message proto.Message
	}

	// HostErrorsMap provides a mapping from error strings to a set of
	// hosts to which the error applies.
	HostErrorsMap map[string]*hostlist.HostSet

	// HostResponseChan defines a channel of *HostResponse items returned
	// from asynchronous unary RPC invokers.
	HostResponseChan chan *HostResponse

	// UnaryResponse contains a slice of *HostResponse items returned
	// from synchronous unary RPC invokers.
	UnaryResponse struct {
		Responses []*HostResponse
		fromMS    bool
	}
)

// Add creates or updates the err/addr keyval pair.
func (hem HostErrorsMap) Add(hostAddr string, hostErr error) (err error) {
	if hostErr == nil {
		return nil
	}

	errStr := hostErr.Error() // stringify the error as map key
	if _, exists := hem[errStr]; !exists {
		hem[errStr], err = hostlist.CreateSet(hostAddr)
		return
	}
	_, err = hem[errStr].Insert(hostAddr)
	return
}

// Keys returns a stable sorted slice of the errors map keys.
func (hem HostErrorsMap) Keys() []string {
	setToKeys := make(map[string]map[string]struct{})
	for errStr, set := range hem {
		rs := set.RangedString()
		if _, exists := setToKeys[rs]; !exists {
			setToKeys[rs] = make(map[string]struct{})
		}
		setToKeys[rs][errStr] = struct{}{}
	}

	sets := make([]string, 0, len(hem))
	for set := range setToKeys {
		sets = append(sets, set)
	}
	sort.Strings(sets)

	keys := make([]string, 0, len(hem))
	for _, set := range sets {
		setKeys := make([]string, 0, len(setToKeys[set]))
		for key := range setToKeys[set] {
			setKeys = append(setKeys, key)
		}
		sort.Strings(setKeys)
		keys = append(keys, setKeys...)
	}
	return keys
}

// getMSResponse is a helper method to return the MS response
// message from a UnaryResponse.
func (ur *UnaryResponse) getMSResponse() (proto.Message, error) {
	if ur == nil {
		return nil, errors.Errorf("nil %T", ur)
	}

	if !ur.fromMS {
		return nil, errors.New("response did not come from management service")
	}

	if len(ur.Responses) == 0 {
		return nil, errors.New("response did not contain a management service response")
	}

	msResp := ur.Responses[0]
	if msResp.Error != nil {
		return nil, msResp.Error
	}

	return msResp.Message, nil
}

// convertMSResponse is a helper function to extract the MS response
// message from a generic UnaryResponse. The out parameter must be
// a reference to a concrete type (e.g. PoolQueryResp).
func convertMSResponse(ur *UnaryResponse, out interface{}) error {
	msResp, err := ur.getMSResponse()
	if err != nil {
		return errors.Wrap(err, "failed to get MS response")
	}

	return convert.Types(msResp, out)
}

// dialOptions is a helper method to return a set of gRPC
// client dialer options.
func (c *Client) dialOptions() ([]grpc.DialOption, error) {
	opts := []grpc.DialOption{
		streamErrorInterceptor(),
		unaryErrorInterceptor(),
		grpc.FailOnNonTempDialError(true),
	}

	creds, err := security.DialOptionForTransportConfig(c.config.TransportConfig)
	if err != nil {
		return nil, err
	}
	opts = append(opts, creds)

	return opts, nil
}

// InvokeUnaryRPCAsync performs an asynchronous invocation of the given RPC
// across all hosts in the provided host list. The returned HostResponseChan
// provides access to a stream of HostResponse items as they are received, and
// is closed when no more responses are expected.
func (c *Client) InvokeUnaryRPCAsync(parent context.Context, req UnaryRequest) (HostResponseChan, error) {
	hosts, err := getRequestHosts(c.config, req)
	if err != nil {
		return nil, err
	}

	c.Debugf("request hosts: %v", hosts)

	// TODO: Make this a const, set appropriately
	// Requests should optionally be able to set a timeout?
	// Per-server vs per-sync-req timeouts?
	requestTimeout := 5 * time.Second

	// TODO: Explore strategies for rate-limiting or batching as necessary
	// in order to perform adequately at scale.
	respChan := make(HostResponseChan, len(hosts))
	go func() {
		ctx, cancel := context.WithTimeout(parent, requestTimeout)
		defer cancel()

		var wg sync.WaitGroup
		for _, host := range hosts {
			wg.Add(1)
			go func(hostAddr string) {
				var msg proto.Message
				opts, err := c.dialOptions()
				if err == nil {
					var conn grpc.ClientConnInterface
					conn, err = grpc.DialContext(ctx, hostAddr, opts...)
					if err == nil {
						msg, err = req.getRPC()(ctx, conn)
					}
				}

				select {
				case <-ctx.Done():
				case respChan <- &HostResponse{Addr: hostAddr, Error: err, Message: msg}:
				}
				wg.Done()
			}(host)
		}
		wg.Wait()
		close(respChan)
	}()

	return respChan, nil
}

// InvokeUnaryRPC performs a synchronous (blocking) invocation of the request's
// RPC across all hosts in the request. The response contains a slice of HostResponse
// items which represent the success or failure of the RPC invocation for each host
// in the request.
func (c *Client) InvokeUnaryRPC(ctx context.Context, req UnaryRequest) (*UnaryResponse, error) {
	respChan, err := c.InvokeUnaryRPCAsync(ctx, req)
	if err != nil {
		return nil, err
	}

	ur := &UnaryResponse{
		fromMS: req.isMSRequest(),
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case hr := <-respChan:
			if hr == nil {
				return ur, nil
			}
			ur.Responses = append(ur.Responses, hr)
		}
	}
}
