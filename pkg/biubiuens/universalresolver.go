package biubiulib

import (
	"Open_IM/pkg/biubiuens/contract/universalresolver"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

func ReverseResolveEnsName(backend bind.ContractBackend, contractAddress common.Address, address common.Address) (string, common.Address, common.Address, common.Address, error) {
	resolver, err := universalresolver.NewContract(contractAddress, backend)
	if err != nil {
		return "", UnknownAddress, UnknownAddress, UnknownAddress, err
	}
	fmt.Println(resolver.Owner(nil))
	st := fmt.Sprintf("%s.addr.reverse", strings.ToLower(address.Hex())[2:])
	_, err = NameHash(fmt.Sprintf("%s.addr.reverse", address.Hex()[2:]))
	if err != nil {
		return "", UnknownAddress, UnknownAddress, UnknownAddress, err
	}
	fmt.Printf("%x\n", EncodeName(st))
	// Resolve the name

	return resolver.Reverse0(&bind.CallOpts{}, EncodeName(st))

}
