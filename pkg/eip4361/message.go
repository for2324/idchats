package eip4361

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/relvacode/iso8601"
)

type Message struct {
	Domain         string         `json:"domain"`
	Address        common.Address `json:"address"`
	Uri            url.URL        `json:"uri"`
	Version        string         `json:"version"`
	Statement      *string        `json:"statement"`
	Nonce          string         `json:"nonce"`
	ChainID        int            `json:"chainID"`
	IssuedAt       string         `json:"issuedAt"`
	ExpirationTime *string        `json:"expirationTime"`
	NotBefore      *string        `json:"notBefore"`
	RequestID      *string        `json:"requestID"`
	Resources      []url.URL      `json:"resources"`
}

func (m *Message) GetDomain() string {
	return m.Domain
}

func (m *Message) GetAddress() common.Address {
	return m.Address
}

func (m *Message) GetURI() url.URL {
	return m.Uri
}

func (m *Message) GetVersion() string {
	return m.Version
}

func (m *Message) GetStatement() *string {
	if m.Statement != nil {
		ret := *m.Statement
		return &ret
	}
	return nil
}

func (m *Message) GetNonce() string {
	return m.Nonce
}

func (m *Message) GetChainID() int {
	return m.ChainID
}

func (m *Message) GetIssuedAt() string {
	return m.IssuedAt
}

func (m *Message) getExpirationTime() *time.Time {
	if !isEmpty(m.ExpirationTime) {
		ret, _ := iso8601.ParseString(*m.ExpirationTime)
		return &ret
	}
	return nil
}

func (m *Message) GetExpirationTime() *string {
	if m.ExpirationTime != nil {
		ret := *m.ExpirationTime
		return &ret
	}
	return nil
}

func (m *Message) getNotBefore() *time.Time {
	if !isEmpty(m.NotBefore) {
		ret, _ := iso8601.ParseString(*m.NotBefore)
		return &ret
	}
	return nil
}

func (m *Message) GetNotBefore() *string {
	if m.NotBefore != nil {
		ret := *m.NotBefore
		return &ret
	}
	return nil
}

func (m *Message) GetRequestID() *string {
	if m.RequestID != nil {
		ret := *m.RequestID
		return &ret
	}
	return nil
}

func (m *Message) GetResources() []url.URL {
	return m.Resources
}

func buildAuthority(uri *url.URL) string {
	authority := uri.Host
	if uri.User != nil {
		authority = fmt.Sprintf("%s@%s", uri.User.String(), authority)
	}
	return authority
}

func validateDomain(domain *string) (bool, error) {
	if isEmpty(domain) {
		return false, &InvalidMessage{"`domain` must not be empty"}
	}

	validateDomain, err := url.Parse(fmt.Sprintf("https://%s", *domain))
	if err != nil {
		return false, &InvalidMessage{"Invalid format for field `domain`"}
	}

	authority := buildAuthority(validateDomain)
	if authority != *domain {
		return false, &InvalidMessage{"Invalid format for field `domain`"}
	}

	return true, nil
}

func validateURI(uri *string) (*url.URL, error) {
	if isEmpty(uri) {
		return nil, &InvalidMessage{"`uri` must not be empty"}
	}

	validateURI, err := url.Parse(*uri)
	if err != nil {
		return nil, &InvalidMessage{"Invalid format for field `uri`"}
	}

	return validateURI, nil
}

// InitMessage creates a Message object with the provided parameters
func InitMessage(domain, address, uri, nonce string, options map[string]interface{}) (*Message, error) {
	if ok, err := validateDomain(&domain); !ok {
		return nil, err
	}

	if isEmpty(&address) {
		return nil, &InvalidMessage{"`address` must not be empty"}
	}

	validateURI, err := validateURI(&uri)
	if err != nil {
		return nil, err
	}

	if isEmpty(&nonce) {
		return nil, &InvalidMessage{"`nonce` must not be empty"}
	}

	var statement *string
	if val, ok := options["statement"]; ok {
		value := val.(string)
		statement = &value
	}

	var chainId int
	if val, ok := options["chainId"]; ok {
		switch val.(type) {
		case float64:
			chainId = int(val.(float64))
		case int:
			chainId = val.(int)
		case string:
			parsed, err := strconv.Atoi(val.(string))
			if err != nil {
				return nil, &InvalidMessage{"Invalid format for field `chainId`, must be an integer"}
			}
			chainId = parsed
		default:
			return nil, &InvalidMessage{"`chainId` must be a string or a integer"}
		}
	} else {
		chainId = 1
	}

	var issuedAt string
	timestamp, err := parseTimestamp(options, "issuedAt")
	if err != nil {
		return nil, err
	}

	if timestamp != nil {
		issuedAt = *timestamp
	} else {
		issuedAt = time.Now().UTC().Format(time.RFC3339)
	}

	var expirationTime *string
	timestamp, err = parseTimestamp(options, "expirationTime")
	if err != nil {
		return nil, err
	}

	if timestamp != nil {
		expirationTime = timestamp
	}

	var notBefore *string
	timestamp, err = parseTimestamp(options, "notBefore")
	if err != nil {
		return nil, err
	}

	if timestamp != nil {
		notBefore = timestamp
	}

	var requestID *string
	if val, ok := isStringAndNotEmpty(options, "requestId"); ok {
		requestID = val
	}

	var resources []url.URL
	if val, ok := options["resources"]; ok {
		switch val.(type) {
		case []url.URL:
			resources = val.([]url.URL)
		default:
			return nil, &InvalidMessage{"`resources` must be a []url.URL"}
		}
	}

	return &Message{
		Domain:  domain,
		Address: common.HexToAddress(address),
		Uri:     *validateURI,
		Version: "1",

		Statement: statement,
		Nonce:     nonce,
		ChainID:   chainId,

		IssuedAt:       issuedAt,
		ExpirationTime: expirationTime,
		NotBefore:      notBefore,

		RequestID: requestID,
		Resources: resources,
	}, nil
}

func parseMessage(message string) (map[string]interface{}, error) {
	match := _SIWE_MESSAGE.FindStringSubmatch(message)

	if match == nil {
		return nil, &InvalidMessage{"Message could not be parsed"}
	}

	result := make(map[string]interface{})
	for i, name := range _SIWE_MESSAGE.SubexpNames() {
		if i != 0 && name != "" && match[i] != "" {
			result[name] = match[i]
		}
	}

	if _, ok := result["domain"]; !ok {
		return nil, &InvalidMessage{"`domain` must not be empty"}
	}
	domain := result["domain"].(string)
	if ok, err := validateDomain(&domain); !ok {
		return nil, err
	}

	if _, ok := result["uri"]; !ok {
		return nil, &InvalidMessage{"`domain` must not be empty"}
	}
	uri := result["uri"].(string)
	if _, err := validateURI(&uri); err != nil {
		return nil, err
	}

	originalAddress := result["address"].(string)
	parsedAddress := common.HexToAddress(originalAddress)
	if originalAddress != parsedAddress.String() {
		return nil, &InvalidMessage{"Address must be in EIP-55 format"}
	}

	if val, ok := result["resources"]; ok {
		resources := strings.Split(val.(string), "\n- ")[1:]
		validateResources := make([]url.URL, len(resources))
		for i, resource := range resources {
			validateResource, err := url.Parse(resource)
			if err != nil {
				return nil, &InvalidMessage{fmt.Sprintf("Invalid format for field `resources` at position %d", i)}
			}
			validateResources[i] = *validateResource
		}
		result["resources"] = validateResources
	}

	return result, nil
}

// ParseMessage returns a Message object by parsing an EIP-4361 formatted string
func ParseMessage(message string) (*Message, error) {
	result, err := parseMessage(message)
	if err != nil {
		return nil, err
	}

	parsed, err := InitMessage(
		result["domain"].(string),
		result["address"].(string),
		result["uri"].(string),
		result["nonce"].(string),
		result,
	)

	if err != nil {
		return nil, err
	}

	return parsed, nil
}

// ParseMessage returns a Message object by parsing an EIP-4361 formatted string
func ParseMessageFromParse(message string) (m *Message, err error) {

	m = &Message{}
	r := bytes.NewReader([]byte(message))
	p := parser{
		scanner: bufio.NewScanner(r),
		msg:     m,
	}
	err = p.parse()
	return m, err
}
func (m *Message) eip191Hash() common.Hash {
	// Ref: https://stackoverflow.com/questions/49085737/geth-ecrecover-invalid-signature-recovery-id
	data := []byte(m.String())
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
}

// ValidNow validates the time constraints of the message at current time.
func (m *Message) ValidNow() (bool, error) {
	return m.ValidAt(time.Now().UTC())
}

// ValidAt validates the time constraints of the message at a specific point in time.
func (m *Message) ValidAt(when time.Time) (bool, error) {
	if m.ExpirationTime != nil {
		if when.After(*m.getExpirationTime()) {
			return false, &ExpiredMessage{"Message expired"}
		}
	}

	if m.NotBefore != nil {
		if when.Before(*m.getNotBefore()) {
			return false, &InvalidMessage{"Message not yet valid"}
		}
	}

	return true, nil
}

// VerifyEIP191 validates the integrity of the object by matching it's signature.
func (m *Message) VerifyEIP191(signature string) (*ecdsa.PublicKey, error) {
	if isEmpty(&signature) {
		return nil, &InvalidSignature{"Signature cannot be empty"}
	}

	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return nil, &InvalidSignature{"Failed to decode signature"}
	}

	// Ref:https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	sigBytes[64] %= 27
	if sigBytes[64] != 0 && sigBytes[64] != 1 {
		return nil, &InvalidSignature{"Invalid signature recovery byte"}
	}

	pkey, err := crypto.SigToPub(m.eip191Hash().Bytes(), sigBytes)
	if err != nil {
		return nil, &InvalidSignature{"Failed to recover public key from signature"}
	}

	address := crypto.PubkeyToAddress(*pkey)

	if !strings.EqualFold(address.String(), m.Address.String()) {
		return nil, &InvalidSignature{"Signer address must match message address"}
	}

	return pkey, nil
}

// Verify validates time constraints and integrity of the object by matching it's signature.
func (m *Message) Verify(signature string, domain *string, nonce *string, timestamp *time.Time) (*ecdsa.PublicKey, error) {
	var err error

	if timestamp != nil {
		_, err = m.ValidAt(*timestamp)
	} else {
		_, err = m.ValidNow()
	}

	if err != nil {
		return nil, err
	}

	if domain != nil {
		if m.GetDomain() != *domain {
			return nil, &InvalidSignature{"Message domain doesn't match"}
		}
	}

	if nonce != nil {
		if m.GetNonce() != *nonce {
			return nil, &InvalidSignature{"Message nonce doesn't match"}
		}
	}

	return m.VerifyEIP191(signature)
}

func (m *Message) prepareMessage() string {
	greeting := fmt.Sprintf("%s wants you to sign in with your Ethereum account:", m.Domain)
	headerArr := []string{greeting, m.Address.String()}

	if isEmpty(m.Statement) {
		headerArr = append(headerArr, "\n")
	} else {
		headerArr = append(headerArr, fmt.Sprintf("\n%s\n", *m.Statement))
	}

	header := strings.Join(headerArr, "\n")

	uri := fmt.Sprintf("URI: %s", m.Uri.String())
	version := fmt.Sprintf("Version: %s", m.Version)
	chainId := fmt.Sprintf("Chain ID: %d", m.ChainID)
	nonce := fmt.Sprintf("Nonce: %s", m.Nonce)
	issuedAt := fmt.Sprintf("Issued At: %s", m.IssuedAt)

	bodyArr := []string{uri, version, chainId, nonce, issuedAt}

	if !isEmpty(m.ExpirationTime) {
		value := fmt.Sprintf("Expiration Time: %s", *m.ExpirationTime)
		bodyArr = append(bodyArr, value)
	}

	if !isEmpty(m.NotBefore) {
		value := fmt.Sprintf("Not Before: %s", *m.NotBefore)
		bodyArr = append(bodyArr, value)
	}

	if !isEmpty(m.RequestID) {
		value := fmt.Sprintf("Request ID: %s", *m.RequestID)
		bodyArr = append(bodyArr, value)
	}

	if len(m.Resources) > 0 {
		resourcesArr := make([]string, len(m.Resources))
		for i, v := range m.Resources {
			resourcesArr[i] = fmt.Sprintf("- %s", v.String())
		}

		resources := strings.Join(resourcesArr, "\n")
		value := fmt.Sprintf("Resources:\n%s", resources)

		bodyArr = append(bodyArr, value)
	}

	body := strings.Join(bodyArr, "\n")

	return strings.Join([]string{header, body}, "\n")
}

func (m *Message) String() string {
	return m.prepareMessage()
}
