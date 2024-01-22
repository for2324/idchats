package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type VerifyMessageWithAddrRequest struct {

	// The message to be signed. When using REST, this field must be encoded as
	// base64.
	Msg []byte `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	// The compact ECDSA signature to be verified over the given message
	// ecoded in base64.
	Signature string `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	// The address which will be used to look up the public key and verify the
	// the signature.
	Addr     string `protobuf:"bytes,3,opt,name=addr,proto3" json:"addr,omitempty"`
	NetParam string
}

type VerifyMessageWithAddrResponse struct {
	// Whether the signature was valid over the given message.
	Valid bool `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	// The pubkey recovered from the signature.
	Pubkey []byte `protobuf:"bytes,2,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
}

// doubleHashMessage creates the double hash (sha256) of a message
// prepended with a specified prefix.
func doubleHashMessage(prefix string, msg string) ([]byte, error) {
	var buf bytes.Buffer
	err := wire.WriteVarString(&buf, 0, prefix)
	if err != nil {
		return nil, err
	}

	err = wire.WriteVarString(&buf, 0, msg)
	if err != nil {
		return nil, err
	}
	fmt.Println(hex.EncodeToString(buf.Bytes()))
	digest := chainhash.DoubleHashB(buf.Bytes())
	fmt.Println(hex.EncodeToString(digest))
	return digest, nil
}

const msgSignaturePrefix = "Bitcoin Signed Message:\n"

func VerifyMessageWithAddr(_ context.Context,
	req *VerifyMessageWithAddrRequest) (*VerifyMessageWithAddrResponse,
	error) {

	sig, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		return nil, fmt.Errorf("malformed base64 encoding of "+
			"the signature: %w", err)
	}

	digest, err := doubleHashMessage(msgSignaturePrefix, string(req.Msg))
	if err != nil {
		return nil, err
	}

	pk, wasCompressed, err := ecdsa.RecoverCompact(sig, digest)
	if err != nil {
		return nil, fmt.Errorf("unable to recover public key "+
			"from compact signature: %w", err)
	}

	var serializedPubkey []byte
	if wasCompressed {
		serializedPubkey = pk.SerializeCompressed()
	} else {
		serializedPubkey = pk.SerializeUncompressed()
	}
	//for unisat sign
	serializedPubkey = pk.SerializeCompressed()

	netParam := &chaincfg.TestNet3Params
	switch req.NetParam {
	case "testnet":
		netParam = &chaincfg.TestNet3Params
	case "livenet", "mainnet":
		netParam = &chaincfg.MainNetParams

	}
	addr, err := btcutil.DecodeAddress(req.Addr, netParam)
	if err != nil {
		return nil, fmt.Errorf("unable to decode address: %w", err)
	}

	if !addr.IsForNet(netParam) {
		return nil, fmt.Errorf("encoded address is for"+
			"the wrong network %s", req.Addr)
	}

	var (
		address    btcutil.Address
		pubKeyHash = btcutil.Hash160(serializedPubkey)
	)

	// Ensure the address is one of the supported types.
	switch addr.(type) {
	case *btcutil.AddressPubKeyHash:
		address, err = btcutil.NewAddressPubKeyHash(
			pubKeyHash, netParam,
		)
		if err != nil {
			return nil, err
		}

	case *btcutil.AddressWitnessPubKeyHash:
		address, err = btcutil.NewAddressWitnessPubKeyHash(
			pubKeyHash, netParam,
		)
		if err != nil {
			return nil, err
		}

	case *btcutil.AddressScriptHash:
		// Check if address is a Nested P2WKH (NP2WKH).
		address, err = btcutil.NewAddressWitnessPubKeyHash(
			pubKeyHash, netParam,
		)
		if err != nil {
			return nil, err
		}

		witnessScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return nil, err
		}

		address, err = btcutil.NewAddressScriptHashFromHash(
			btcutil.Hash160(witnessScript), netParam,
		)
		if err != nil {
			return nil, err
		}

	case *btcutil.AddressTaproot:
		// Only addresses without a tapscript are allowed because
		// the verification is using the internal key.
		tapKey := txscript.ComputeTaprootKeyNoScript(pk)
		address, err = btcutil.NewAddressTaproot(
			schnorr.SerializePubKey(tapKey),
			netParam,
		)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported address type")
	}

	return &VerifyMessageWithAddrResponse{
		Valid:  req.Addr == address.EncodeAddress(),
		Pubkey: serializedPubkey,
	}, nil
}
