// Copyright Microsoft Corp.
// All rights reserved.

package ipam

const (
	// Libnetwork IPAM plugin endpoint type
	endpointType = "IpamDriver"

	// Libnetwork IPAM plugin remote API paths
	capabilitiesPath   = "/IpamDriver.GetCapabilities"
	addressSpacesPath  = "/IpamDriver.GetDefaultAddressSpaces"
	requestPoolPath    = "/IpamDriver.RequestPool"
	releasePoolPath    = "/IpamDriver.ReleasePool"
	requestAddressPath = "/IpamDriver.RequestAddress"
	releaseAddressPath = "/IpamDriver.ReleaseAddress"
)

// Response sent by plugin when registering its capabilities with libnetwork.
type getCapabilitiesResponse struct {
	RequiresMACAddress bool
}

// Response sent by plugin when returning the default address space names.
type getDefaultAddressSpacesResponse struct {
	LocalDefaultAddressSpace  string
	GlobalDefaultAddressSpace string
}

// Request sent by libnetwork when acquiring a reference to an address pool.
type requestPoolRequest struct {
	AddressSpace string
	Pool         string
	SubPool      string
	Options      map[string]string
	V6           bool
}

// Response sent by plugin when an address pool is successfully referenced.
type requestPoolResponse struct {
	PoolID string
	Pool   string
	Data   map[string]string
}

// Request sent by libnetwork when releasing a previously registered address pool.
type releasePoolRequest struct {
	PoolID string
}

// Response sent by plugin when an address pool is successfully released.
type releasePoolResponse struct {
}

// Request sent by libnetwork when reserving an address from a pool.
type requestAddressRequest struct {
	PoolID  string
	Address string
	Options map[string]string
}

// Response sent by plugin when an address is successfully reserved.
type requestAddressResponse struct {
	Address string
	Data    map[string]string
}

// Request sent by libnetwork when releasing an address back to the pool.
type releaseAddressRequest struct {
	PoolID  string
	Address string
}

// Response sent by plugin when an address is successfully released.
type releaseAddressResponse struct {
}