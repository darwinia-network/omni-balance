// Code generated by go-swagger; DO NOT EDIT.

package notifications

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// New creates a new notifications API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

// New creates a new notifications API client with basic auth credentials.
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

// New creates a new notifications API client with a bearer token for authentication.
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
Client for notifications API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption may be used to customize the behavior of Client methods.
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	V1NotificationsDevicesCreate(params *V1NotificationsDevicesCreateParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1NotificationsDevicesCreateOK, error)

	V1NotificationsDevicesDelete(params *V1NotificationsDevicesDeleteParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1NotificationsDevicesDeleteNoContent, error)

	V1NotificationsDevicesSafesDelete(params *V1NotificationsDevicesSafesDeleteParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1NotificationsDevicesSafesDeleteNoContent, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
	V1NotificationsDevicesCreate Creates a new FirebaseDevice. If uuid is not provided a new device will be created.

If a uuid for an existing Safe is provided the FirebaseDevice will be updated with all the new data provided.
Safes provided on the request are always added and never removed/replaced
Signature must sign `keccack256('gnosis-safe{timestamp-epoch}{uuid}{cloud_messaging_token}{safes_sorted}':
  - `{timestamp-epoch}` must be an integer (no milliseconds)
  - `{safes_sorted}` must be checksummed safe addresses sorted and joined with no spaces
*/
func (a *Client) V1NotificationsDevicesCreate(params *V1NotificationsDevicesCreateParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1NotificationsDevicesCreateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1NotificationsDevicesCreateParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_notifications_devices_create",
		Method:             "POST",
		PathPattern:        "/v1/notifications/devices/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1NotificationsDevicesCreateReader{formats: a.formats},
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
	success, ok := result.(*V1NotificationsDevicesCreateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_notifications_devices_create: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
V1NotificationsDevicesDelete Remove a FirebaseDevice
*/
func (a *Client) V1NotificationsDevicesDelete(params *V1NotificationsDevicesDeleteParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1NotificationsDevicesDeleteNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1NotificationsDevicesDeleteParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_notifications_devices_delete",
		Method:             "DELETE",
		PathPattern:        "/v1/notifications/devices/{uuid}/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1NotificationsDevicesDeleteReader{formats: a.formats},
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
	success, ok := result.(*V1NotificationsDevicesDeleteNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_notifications_devices_delete: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
V1NotificationsDevicesSafesDelete Remove a Safe for a FirebaseDevice
*/
func (a *Client) V1NotificationsDevicesSafesDelete(params *V1NotificationsDevicesSafesDeleteParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*V1NotificationsDevicesSafesDeleteNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewV1NotificationsDevicesSafesDeleteParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "v1_notifications_devices_safes_delete",
		Method:             "DELETE",
		PathPattern:        "/v1/notifications/devices/{uuid}/safes/{address}/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &V1NotificationsDevicesSafesDeleteReader{formats: a.formats},
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
	success, ok := result.(*V1NotificationsDevicesSafesDeleteNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for v1_notifications_devices_safes_delete: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
