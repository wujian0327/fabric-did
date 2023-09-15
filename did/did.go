package did

import (
	"encoding"
	"encoding/json"
	"errors"
	"fabric-did/ssi"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var _ fmt.Stringer = DID{}
var _ encoding.TextMarshaler = DID{}

var didPattern = regexp.MustCompile(`^did:([a-z0-9]+):((?:(?:[a-zA-Z0-9.\-_:])+|(?:%[0-9a-fA-F]{2})+)+)(/.*?|)(\?.*?|)(#.*|)$`)

// DIDContextV1 contains the JSON-LD context for a DID Document
const DIDContextV1 = "https://www.w3.org/ns/did/v1"

// DID represent a Decentralized Identifier as specified by the DID Core specification (https://www.w3.org/TR/did-core/#identifier).
type DID struct {
	// Method is the DID method, e.g. "example".
	Method string
	// ID is the method-specific ID, in escaped form.
	ID string
	// DecodedID is the method-specific ID, in unescaped form.
	// It is only set during parsing, and not used by the String() method.
	DecodedID string
	// Path is the DID path without the leading '/', in escaped form.
	Path string
	// DecodedPath is the DID path without the leading '/', in unescaped form.
	// It is only set during parsing, and not used by the String() method.
	DecodedPath string
	// Query contains the DID query key-value pairs, in unescaped form.
	// String() will escape the values again, and order the keys alphabetically.
	Query url.Values
	// Fragment is the DID fragment without the leading '#', in escaped form.
	Fragment string
	// DecodedFragment is the DID fragment without the leading '#', in unescaped form.
	// It is only set during parsing, and not used by the String() method.
	DecodedFragment string
}

// Empty checks whether the DID is set or not
func (d DID) Empty() bool {
	return d.Method == ""
}

// String returns the DID as formatted string.
func (d DID) String() string {
	if d.Empty() {
		return ""
	}
	result := "did:" + d.Method + ":" + d.ID
	if d.Path != "" {
		result += "/" + d.Path
	}
	if len(d.Query) > 0 {
		result += "?" + d.Query.Encode()
	}
	if d.Fragment != "" {
		result += "#" + d.Fragment
	}
	return result
}

// URI converts the DID to an URI.
// URIs are used in Verifiable Credentials
func (d DID) URI() ssi.URI {
	return ssi.URI{
		URL: url.URL{
			Scheme:   "did",
			Opaque:   fmt.Sprintf("%s:%s", d.Method, url.PathEscape(d.ID)),
			Fragment: d.Fragment,
		},
	}
}

// MarshalText implements encoding.TextMarshaler
func (d DID) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// Equals checks whether the DID equals to another DID.
// When the DIDs
// The check is case-sensitive.
func (d DID) Equals(other DID) bool {
	return d.cleanup().String() == other.cleanup().String()
}

func (d DID) cleanup() DID {
	if len(d.Query) == 0 {
		d.Query = nil
	}
	return d
}

// UnmarshalJSON unmarshals a DID encoded as JSON string, e.g.:
// "did:nuts:c0dc584345da8a0e1e7a584aa4a36c30ebdb79d907aff96fe0e90ee972f58a17"
func (d *DID) UnmarshalJSON(bytes []byte) error {
	var didString string
	err := json.Unmarshal(bytes, &didString)
	if err != nil {
		return ErrInvalidDID.wrap(err)
	}
	tmp, err := ParseDIDURL(didString)
	if err != nil {
		return err
	}
	*d = *tmp
	return nil
}

func (d *DID) IsURL() bool {
	return d.Fragment != "" || len(d.Query) != 0 || d.Path != ""
}

// MarshalJSON marshals the DID to a JSON string
func (d DID) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// WithoutURL returns a copy of the DID without URL parts (fragment, query, path).
func (d DID) WithoutURL() DID {
	return DID{
		Method:    d.Method,
		ID:        d.ID,
		DecodedID: d.DecodedID,
	}
}

// DIDContextV1URI returns DIDContextV1 as a URI
func DIDContextV1URI() ssi.URI {
	return ssi.MustParseURI(DIDContextV1)
}

// ParseDIDURL parses a DID URL.
// https://www.w3.org/TR/did-core/#did-url-syntax
// A DID URL is a URL that builds on the DID scheme.
func ParseDIDURL(input string) (*DID, error) {
	// There are 6 submatches (base 0)
	// 0. complete DID
	// 1. method
	// 2. id
	// 3. path (starting with '/')
	// 4. query (starting with '?')
	// 5. fragment (starting with '#')
	matches := didPattern.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, ErrInvalidDID
	}

	result := DID{
		Method:   matches[1],
		ID:       matches[2],
		Path:     strings.TrimPrefix(matches[3], "/"),
		Fragment: strings.TrimPrefix(matches[5], "#"),
	}
	var err error
	result.DecodedID, err = url.PathUnescape(result.ID)
	if err != nil {
		return nil, ErrInvalidDID.wrap(fmt.Errorf("invalid ID: %w", err))
	}
	result.DecodedPath, err = url.PathUnescape(result.Path)
	if err != nil {
		return nil, ErrInvalidDID.wrap(fmt.Errorf("invalid path: %w", err))
	}
	result.DecodedFragment, err = url.PathUnescape(result.Fragment)
	if err != nil {
		return nil, ErrInvalidDID.wrap(fmt.Errorf("invalid fragment: %w", err))
	}
	result.Query, err = url.ParseQuery(strings.TrimPrefix(matches[4], "?"))
	if err != nil {
		return nil, ErrInvalidDID.wrap(err)
	}
	result = result.cleanup()
	return &result, nil
}

// ParseDID parses a raw DID.
// If the input contains a path, query or fragment, use the ParseDIDURL instead.
// If it can't be parsed, an error is returned.
func ParseDID(input string) (*DID, error) {
	did, err := ParseDIDURL(input)
	if err != nil {
		return nil, err
	}
	if did.IsURL() {
		return nil, ErrInvalidDID.wrap(errors.New("DID can not have path, fragment or query params"))
	}
	return did, nil
}

// must accepts a function like Parse and returns the value without error or panics otherwise.
func must(fn func(string) (*DID, error), input string) DID {
	v, err := fn(input)
	if err != nil {
		panic(err)
	}
	return *v
}

// MustParseDID is like ParseDID but panics if the string cannot be parsed.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustParseDID(input string) DID {
	return must(ParseDID, input)
}

// MustParseDIDURL is like ParseDIDURL but panics if the string cannot be parsed.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustParseDIDURL(input string) DID {
	return must(ParseDIDURL, input)
}

// ErrInvalidDID is returned when a parser function is supplied with a string that can't be parsed as DID.
var ErrInvalidDID = ParserError{msg: "invalid DID"}

// ParserError is used when returning DID-parsing related errors.
type ParserError struct {
	msg string
	err error
}

func (w ParserError) wrap(err error) error {
	return ParserError{msg: fmt.Sprintf("%s: %s", w.msg, err.Error()), err: err}
}

// Is checks whether the given error is a ParserError
func (w ParserError) Is(other error) bool {
	_, ok := other.(ParserError)
	return ok
}

// Unwrap returns the underlying error.
func (w ParserError) Unwrap() error {
	return w.err
}

// Error returns the message of the error.
func (w ParserError) Error() string {
	return w.msg
}