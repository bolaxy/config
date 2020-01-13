package conf

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/boxproject/bolaxy/common"
	"github.com/boxproject/bolaxy/crypto"
)

//Peer is a struct that holds Peer data
type Peer struct {
	Alias     string `mapstructure:"alias"`
	PubKeyHex string `mapstructure:"pubkey"`
	Address   string `mapstructure:"address"`
	HttpPort  string `mapstructure:"httpport"`
	TcpPort   string `mapstructure:"tcpport"`

	id uint32
}

//NewPeer is a factory method for creating a new Peer instance
func NewPeer(pubKeyHex, netAddr, alias, httpPort, tcpPort string) *Peer {
	peer := &Peer{
		PubKeyHex: strings.ToUpper(pubKeyHex),
		Address:   netAddr,
		Alias:     alias,
		HttpPort:  httpPort,
		TcpPort:   tcpPort,
	}
	return peer
}

//ID returns an ID for the peer, calculating a hash is one is not available
//XXX Not very nice
func (p *Peer) ID() uint32 {
	if p.id == 0 {
		pubKeyBytes := p.PubKeyBytes()
		p.id = crypto.Hash32(pubKeyBytes)
	}
	return p.id
}

//PubKeyString returns the upper-case version of PubKeyHex. It is used for
//indexing in maps with string keys.
//XXX do something nicer
func (p *Peer) PubKeyString() string {
	return p.PubKeyHex
}

//PubKeyBytes converts hex string representation of the public key and returns a byte array
func (p *Peer) PubKeyBytes() []byte {
	return common.FromHex(p.PubKeyHex)
}

//Marshal marshals the json representation of the peer
//json encoding excludes the ID field
func (p *Peer) Marshal() ([]byte, error) {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)

	if err := enc.Encode(p); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

//Unmarshal generates a JSON representation of the peer
func (p *Peer) Unmarshal(data []byte) error {
	b := bytes.NewBuffer(data)

	dec := json.NewDecoder(b) //will read from b

	if err := dec.Decode(p); err != nil {
		return err
	}

	return nil
}

func (p *Peer) HttpAddress() string {
	return p.Address + ":" + p.HttpPort
}

func (p *Peer) TcpAddress() string {
	return p.Address + ":" + p.TcpPort
}

// ExcludePeer is used to exclude a single peer from a list of peers.
func ExcludePeer(peers []*Peer, peerID uint32) (int, []*Peer) {
	index := -1
	otherPeers := make([]*Peer, 0, len(peers))
	for i, p := range peers {
		if p.ID() != peerID {
			otherPeers = append(otherPeers, p)
		} else {
			index = i
		}
	}
	return index, otherPeers
}
