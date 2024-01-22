package eip4361

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const _SIWE_DOMAIN = "(?P<domain>([^/?#]+)) wants you to sign in with your Ethereum account:\\n"
const _SIWE_ADDRESS = "(?P<address>0x[a-zA-Z0-9]{40})\\n\\n"
const _SIWE_STATEMENT = "((?P<statement>[^\\n]+)\\n)?\\n"
const _RFC3986 = "(([^ :/?#]+):)?(//([^ /?#]*))?([^ ?#]*)(\\?([^ #]*))?(#(.*))?"

var _SIWE_URI_LINE = fmt.Sprintf("URI: (?P<uri>%s?)\\n", _RFC3986)

const _SIWE_VERSION = "Version: (?P<version>1)\\n"
const _SIWE_CHAIN_ID = "Chain ID: (?P<chainId>[0-9]+)\\n"
const _SIWE_NONCE = "Nonce: (?P<nonce>[a-zA-Z0-9]{8,})\\n"
const _SIWE_DATETIME = "([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))"

var _SIWE_ISSUED_AT = fmt.Sprintf("Issued At: (?P<issuedAt>%s)", _SIWE_DATETIME)
var _SIWE_EXPIRATION_TIME = fmt.Sprintf("(\\nExpiration Time: (?P<expirationTime>%s))?", _SIWE_DATETIME)
var _SIWE_NOT_BEFORE = fmt.Sprintf("(\\nNot Before: (?P<notBefore>%s))?", _SIWE_DATETIME)

const _SIWE_REQUEST_ID = "(\\nRequest ID: (?P<requestId>[-._~!$&'()*+,;=:@%a-zA-Z0-9]*))?"

var _SIWE_RESOURCES = fmt.Sprintf("(\\nResources:(?P<resources>(\\n- %s)+))?", _RFC3986)

var _SIWE_MESSAGE = regexp.MustCompile(fmt.Sprintf("^%s%s%s%s%s%s%s%s%s%s%s%s$",
	_SIWE_DOMAIN,
	_SIWE_ADDRESS,
	_SIWE_STATEMENT,
	_SIWE_URI_LINE,
	_SIWE_VERSION,
	_SIWE_CHAIN_ID,
	_SIWE_NONCE,
	_SIWE_ISSUED_AT,
	_SIWE_EXPIRATION_TIME,
	_SIWE_NOT_BEFORE,
	_SIWE_REQUEST_ID,
	_SIWE_RESOURCES))

const (
	V1 string = "1"

	TimeLayout = "2006-01-02T15:04:05.000Z07:00"

	DomainMessage        = " wants you to sign in with your Ethereum account:"
	URIPrefix            = "URI: "
	VersionPrefix        = "Version: "
	ChainIDPrefix        = "Chain ID: "
	NoncePrefix          = "Nonce: "
	IssuedAtPrefix       = "Issued At: "
	ExpirationTimePrefix = "Expiration Time: "
	NotBeforePrefix      = "Not Before: "
	RequestIDPrefix      = "Request ID: "
	ResourcesPrefix      = "Resources:"
	ResourcePrefix       = "- "
)

type parser struct {
	scanner *bufio.Scanner
	msg     *Message
	prev    *string
	err     error
}

func (p *parser) parse() error {
	ok := p.ruleDomain() &&
		p.ruleAddress() &&
		p.ruleEmptyLine() &&
		p.ruleStatement() &&
		p.ruleEmptyLine() &&
		p.ruleURI() &&
		p.ruleVersion() &&
		p.ruleChainID() &&
		p.ruleNonce() &&
		p.ruleIssuedAt() &&
		p.ruleExpirationTime() &&
		p.ruleNotBefore() &&
		p.ruleRequestID() &&
		p.ruleResources()

	if !ok {
		if p.err != nil {
			return p.err
		} else {
			return io.ErrUnexpectedEOF
		}
	}

	return nil
}

func (p *parser) nextLine() (string, bool) {
	if p.err != nil {
		return "", false
	}

	if p.prev != nil {
		l := *p.prev
		p.prev = nil
		return l, true
	}

	if p.scanner.Scan() {
		return p.scanner.Text(), true
	}

	if err := p.scanner.Err(); err != nil {
		p.err = err
	}

	return "", false
}

func (p *parser) optionalRule(rule func(l string) (bool, error)) bool {
	l, ok := p.nextLine()
	if !ok {
		return true
	}

	if ok, err := rule(l); err != nil {
		p.err = err
		return false
	} else if !ok {
		p.prev = &l
	}

	return true
}

func (p *parser) ruleEmptyLine() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}

	if l != "" {
		p.err = errors.New("empty line expected")
		return false
	}

	return true
}

func (p *parser) ruleDomain() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}

	if !strings.HasSuffix(l, DomainMessage) {
		p.err = errors.New("invalid domain line")
		return false
	}

	d := strings.SplitN(l, " ", 2)[0]
	p.msg.Domain = d
	return true
}

func (p *parser) ruleAddress() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}
	p.msg.Address = common.HexToAddress(l)
	return true
}

func (p *parser) ruleStatement() bool {
	return p.optionalRule(func(l string) (bool, error) {
		if l == "" {
			return false, nil
		}

		p.msg.Statement = &l
		return true, nil
	})
}

func (p *parser) ruleURI() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}

	if !strings.HasPrefix(l, URIPrefix) {
		p.err = errors.New("invalid URI line")
		return false
	}

	tempUrl, _ := url.Parse(l[len(URIPrefix):])
	p.msg.Uri = *tempUrl
	return true
}

func (p *parser) ruleVersion() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}

	if !strings.HasPrefix(l, VersionPrefix) {
		p.err = errors.New("invalid Version line")
		return false
	}

	v := l[len(VersionPrefix):]
	if v != string(V1) {
		p.err = errors.New("version not supported")
		return false
	}

	p.msg.Version = V1
	return true
}

func (p *parser) ruleChainID() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}

	if !strings.HasPrefix(l, ChainIDPrefix) {
		p.err = errors.New("invalid ChainID line")
		return false
	}

	id, err := strconv.ParseInt(l[len(ChainIDPrefix):], 10, 64)
	if err != nil {
		p.err = err
		return false
	}

	p.msg.ChainID = int(id)
	return true
}

func (p *parser) ruleNonce() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}

	if !strings.HasPrefix(l, NoncePrefix) {
		p.err = errors.New("invalid Nonce line")
		return false
	}

	p.msg.Nonce = l[len(NoncePrefix):]
	return true
}

func (p *parser) ruleIssuedAt() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}

	if !strings.HasPrefix(l, IssuedAtPrefix) {
		p.err = errors.New("invalid IssuedAt line")
		return false
	}

	//v, err := time.Parse(TimeLayout, l[len(IssuedAtPrefix):])
	//if err != nil {
	//	p.err = err
	//	return false
	//}

	p.msg.IssuedAt = l[len(IssuedAtPrefix):]
	return true
}

func (p *parser) ruleExpirationTime() bool {
	return p.optionalRule(func(l string) (bool, error) {
		if !strings.HasPrefix(l, ExpirationTimePrefix) {
			return false, nil
		}

		//v, err := time.Parse(TimeLayout, l[len(ExpirationTimePrefix):])
		//if err != nil {
		//	return false, err
		//}
		stringData := l[len(ExpirationTimePrefix):]
		p.msg.ExpirationTime = &stringData
		return true, nil
	})
}

func (p *parser) ruleNotBefore() bool {
	return p.optionalRule(func(l string) (bool, error) {
		if !strings.HasPrefix(l, NotBeforePrefix) {
			return false, nil
		}

		//v, err := time.Parse(TimeLayout, l[len(NotBeforePrefix):])
		//if err != nil {
		//	return false, err
		//}
		notbefore := l[len(NotBeforePrefix):]
		p.msg.NotBefore = &notbefore
		return true, nil
	})
}

func (p *parser) ruleRequestID() bool {
	return p.optionalRule(func(l string) (bool, error) {
		if !strings.HasPrefix(l, RequestIDPrefix) {
			return false, nil
		}

		id := l[len(RequestIDPrefix):]
		p.msg.RequestID = &id
		return true, nil
	})
}

func (p *parser) ruleResources() bool {
	return p.optionalRule(func(l string) (bool, error) {
		if !strings.HasPrefix(l, ResourcesPrefix) {
			return false, nil
		}

		for p.ruleResource() {
		}

		return true, nil
	})
}

func (p *parser) ruleResource() bool {
	l, ok := p.nextLine()
	if !ok {
		return false
	}

	if !strings.HasPrefix(l, ResourcePrefix) {
		return false
	}
	tempData, _ := url.Parse(l[len(ResourcePrefix):])
	p.msg.Resources = append(p.msg.Resources, *tempData)
	return true
}
