package biubiulib

import (
	"Open_IM/pkg/biubiuens/contract/registry"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Registry is the structure for the registry contract
type Registry struct {
	backend      bind.ContractBackend
	Contract     *registry.Contract
	ContractAddr common.Address
}

// NewRegistry obtains the ENS registry
func NewRegistry(backend bind.ContractBackend) (*Registry, error) {
	address, err := RegistryContractAddress(backend)
	if err != nil {
		return nil, err
	}
	return NewRegistryAt(backend, address)
}

// NewRegistryAt obtains the ENS registry at a given address
func NewRegistryAt(backend bind.ContractBackend, address common.Address) (*Registry, error) {
	contract, err := registry.NewContract(address, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{
		backend:      backend,
		Contract:     contract,
		ContractAddr: address,
	}, nil
}

// Owner returns the address of the owner of a name
func (r *Registry) Owner(name string) (common.Address, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return UnknownAddress, err
	}
	return r.Contract.Owner(nil, nameHash)
}

// ResolverAddress returns the address of the resolver for a name
func (r *Registry) ResolverAddress(name string) (common.Address, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return UnknownAddress, err
	}
	return r.Contract.Resolver(nil, nameHash)
}

// SetResolver sets the resolver for a name
func (r *Registry) SetResolver(opts *bind.TransactOpts, name string, address common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	return r.Contract.SetResolver(opts, nameHash, address)
}

// Resolver returns the resolver for a name
func (r *Registry) Resolver(name string) (*Resolver, error) {
	address, err := r.ResolverAddress(name)
	if err != nil {
		return nil, err
	}
	return NewResolverAt(r.backend, name, address)
}

// SetOwner sets the ownership of a domain
func (r *Registry) SetOwner(opts *bind.TransactOpts, name string, address common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	return r.Contract.SetOwner(opts, nameHash, address)
}

// SetSubdomainOwner sets the ownership of a subdomain, potentially creating it in the process
func (r *Registry) SetSubdomainOwner(opts *bind.TransactOpts, name string, subname string, address common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	labelHash, err := LabelHash(subname)
	if err != nil {
		return nil, err
	}
	return r.Contract.SetSubnodeOwner(opts, nameHash, labelHash, address)
}

// RegistryContractAddress obtains the address of the registry contract for a chain.
// This is (currently) the same for all chains.
func RegistryContractAddress(backend bind.ContractBackend) (common.Address, error) {
	// Instantiate the registry contract.  The same for all chains.
	//00000000000C2E074eC69A0dFb2997BA6C7d2e1e
	return common.HexToAddress("0xc67Ea4083e2333F26DD916f10DB733fE05d2043d"), nil
}

// SetResolver sets the resolver for a name
func SetResolver(session *registry.ContractSession, name string, resolverAddr *common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	return session.SetResolver(nameHash, *resolverAddr)
}

// SetSubdomainOwner sets the owner for a subdomain of a name
func SetSubdomainOwner(session *registry.ContractSession, name string, subdomain string, ownerAddr *common.Address) (*types.Transaction, error) {
	nameHash, err := NameHash(name)
	if err != nil {
		return nil, err
	}
	labelHash, err := LabelHash(subdomain)
	if err != nil {
		return nil, err
	}
	return session.SetSubnodeOwner(nameHash, labelHash, *ownerAddr)
}
