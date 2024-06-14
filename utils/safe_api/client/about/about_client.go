// Code generated by go-swagger; DO NOT EDIT.

package about

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// New creates a new about API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

// New creates a new about API client with basic auth credentials.
// It takes the following parameters:
// - host: http host (github.com).
// - basePath: any base path for the API client ("/v1", "/v3").
// - scheme: http scheme ("http", "https").
// - user: user for basic authentication header.
// - password: password for basic authentication header.
func NewClientWithBasicAuth(host, basePath, scheme, user, password string) ClientService {
	transport := httptransport.New(host, basePath, []string{scheme})
	transport.DefaultAuthentication = httptransport.BasicAuth(user, password)
	return &Client{transport: transport, formats: strfmt.Default}
}

// New creates a new about API client with a bearer token for authentication.
// It takes the following parameters:
// - host: http host (github.com).
// - basePath: any base path for the API client ("/v1", "/v3").
// - scheme: http scheme ("http", "https").
// - bearerToken: bearer token for Bearer authentication header.
func NewClientWithBearerToken(host, basePath, scheme, bearerToken string) ClientService {
	transport := httptransport.New(host, basePath, []string{scheme})
	transport.DefaultAuthentication = httptransport.BearerToken(bearerToken)
	return &Client{transport: transport, formats: strfmt.Default}
}

/*
Client for about API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption may be used to customize the behavior of Client methods.
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	V1AboutDeploymentsList(params *V1AboutDeploymentsListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutDeploymentsListOK, error)

	V1AboutEthereumRPCList(params *V1AboutEthereumRPCListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutEthereumRPCListOK, error)

	V1AboutEthereumTracingRPCList(params *V1AboutEthereumTracingRPCListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutEthereumTracingRPCListOK, error)

	V1AboutIndexingList(params *V1AboutIndexingListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutIndexingListOK, error)

	V1AboutList(params *V1AboutListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutListOK, error)

	V1AboutSingletonsList(params *V1AboutSingletonsListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutSingletonsListOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
V1AboutDeploymentsList Returns a list of safe deployments by version.
*/
func (a *Client) V1AboutDeploymentsList(params *V1AboutDeploymentsListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutDeploymentsListOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1AboutDeploymentsListParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_about_deployments_list",
		Method:             "GET",
		PathPattern:        "/v1/about/deployments/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1AboutDeploymentsListReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*V1AboutDeploymentsListOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_about_deployments_list: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
V1AboutEthereumRPCList Get information about the Ethereum RPC node used by the service
*/
func (a *Client) V1AboutEthereumRPCList(params *V1AboutEthereumRPCListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutEthereumRPCListOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1AboutEthereumRPCListParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_about_ethereum-rpc_list",
		Method:             "GET",
		PathPattern:        "/v1/about/ethereum-rpc/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1AboutEthereumRPCListReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*V1AboutEthereumRPCListOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_about_ethereum-rpc_list: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
V1AboutEthereumTracingRPCList Get information about the Ethereum Tracing RPC node used by the service (if any configured)
*/
func (a *Client) V1AboutEthereumTracingRPCList(params *V1AboutEthereumTracingRPCListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutEthereumTracingRPCListOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1AboutEthereumTracingRPCListParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_about_ethereum-tracing-rpc_list",
		Method:             "GET",
		PathPattern:        "/v1/about/ethereum-tracing-rpc/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1AboutEthereumTracingRPCListReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*V1AboutEthereumTracingRPCListOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_about_ethereum-tracing-rpc_list: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
V1AboutIndexingList Get current indexing status for ERC20/721 events
*/
func (a *Client) V1AboutIndexingList(params *V1AboutIndexingListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutIndexingListOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1AboutIndexingListParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_about_indexing_list",
		Method:             "GET",
		PathPattern:        "/v1/about/indexing/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1AboutIndexingListReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*V1AboutIndexingListOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_about_indexing_list: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
V1AboutList Returns information and configuration of the service
*/
func (a *Client) V1AboutList(params *V1AboutListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutListOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1AboutListParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_about_list",
		Method:             "GET",
		PathPattern:        "/v1/about/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1AboutListReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*V1AboutListOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_about_list: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
V1AboutSingletonsList v1 about singletons list API
*/
func (a *Client) V1AboutSingletonsList(params *V1AboutSingletonsListParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1AboutSingletonsListOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1AboutSingletonsListParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_about_singletons_list",
		Method:             "GET",
		PathPattern:        "/v1/about/singletons/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1AboutSingletonsListReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*V1AboutSingletonsListOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_about_singletons_list: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
