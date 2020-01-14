package conf

import (
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/viper"

	"github.com/bolaxy/common/hexutil"
	"github.com/bolaxy/crypto"
	"github.com/bolaxy/rlp"
)

type Genesis struct {
	CoinBase          string   `mapstructure:"coinbase"`
	ChainID           string   `mapstructure:"chain-id"`
	ConsensusAccounts []string `mapstructure:"consensus-accounts"`
	Alloc             []Alloc  `mapstructure:"alloc"`
	Poa               *PoaMap  `mapstructure:"poa"`
	Launcher          *PoaMap  `mapstructure:"launcher"`
}

type Alloc struct {
	Account     string            `mapstructure:"account"`
	Balance     string            `mapstructure:"balance"`
	Code        string            `mapstructure:"code"`
	Storage     map[string]string `mapstructure:"storage"`
	Authorising bool              `mapstructure:"authorising"`
}

type PoaMap struct {
	Address string            `mapstructure:"address"`
	Balance string            `mapstructure:"balance"`
	Abi     string            `mapstructure:"abi"`
	SubAbi  string            `mapstructure:"subabi"`
	Code    string            `mapstructure:"code"`
	Storage map[string]string `mapstructure:"storage"`
}

type genesis struct {
	ChainID           string
	ExtraData         string
	ConsensusAccounts []string `rlp:"nil"`
	Allocs            []alloc  `rlp:"nil"`
	Poa               *poaMap  `rlp:"nil"`
	Launcher          *poaMap  `rlp:"nil"`
}

type alloc struct {
	Account     string
	Balance     string
	Code        string
	Storage     [][2]string `rlp:"nil"`
	Authorising bool
}

type poaMap struct {
	Address string
	Balance string
	ABI     string
	SubAbi  string
	Code    string
	Storage [][2]string `rlp:"nil"`
}

func (g *Genesis) Hash() ([]byte, error) {
	if g == nil {
		return nil, nil
	}

	buf, err := EncodeRLPGenesis(g)
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(buf), nil
}

func (g *Genesis) HexHash() (string, error) {
	if g == nil {
		return "", nil
	}

	hash, err := g.Hash()
	if err != nil {
		return "", err
	}

	return hexutil.Encode(hash), nil
}

func (g *Genesis) EncodeRLP(w io.Writer) error {
	genesis := &genesis{
		ChainID:           g.ChainID,
		ConsensusAccounts: g.ConsensusAccounts,
	}

	if len(g.Alloc) > 0 {
		genesis.Allocs = translateFromAlloc(g.Alloc)
	}

	if g.Poa != nil {
		genesis.Poa = translateFromPoaMap(g.Poa)
	}

	if g.Launcher != nil {
		genesis.Launcher = translateFromPoaMap(g.Launcher)
	}

	return rlp.Encode(w, genesis)
}

func translateFromAlloc(original []Alloc) []alloc {
	allocs := make([]alloc, 0, len(original))
	for _, a := range original {
		item := alloc{
			Account:     a.Account,
			Balance:     a.Balance,
			Code:        a.Code,
			Authorising: a.Authorising,
		}

		if len(a.Storage) > 0 {
			item.Storage = translateFromStorage(a.Storage)
		}

		allocs = append(allocs, item)
	}

	return allocs
}

func translateFromPoaMap(original *PoaMap) (poa *poaMap) {
	poa = &poaMap{
		Address: original.Address,
		Balance: original.Balance,
		ABI:     original.Abi,
		SubAbi:  original.SubAbi,
		Code:    original.Code,
	}

	if len(original.Storage) > 0 {
		poa.Storage = translateFromStorage(original.Storage)
	}

	return
}

func translateFromStorage(store map[string]string) (target [][2]string) {
	target = make([][2]string, 0, len(store))
	for k, v := range store {
		target = append(target, [2]string{k, v})
	}

	sort.Sort(alphabetic(target))
	return
}

func (g *Genesis) DecodeRLP(s *rlp.Stream) error {
	var genesis genesis
	if err := s.Decode(&genesis); err != nil {
		return err
	}

	g.ChainID = genesis.ChainID
	g.ConsensusAccounts = genesis.ConsensusAccounts

	if len(genesis.Allocs) > 0 {
		g.Alloc = translateToAlloc(genesis.Allocs)
	}

	if genesis.Poa != nil {
		g.Poa = translateToPoaMap(genesis.Poa)
	}

	if genesis.Launcher != nil {
		g.Launcher = translateToPoaMap(genesis.Launcher)
	}

	return nil
}

func translateToAlloc(source []alloc) (target []Alloc) {
	target = make([]Alloc, len(source))
	for i, item := range source {
		data := Alloc{
			Account:     item.Account,
			Balance:     item.Balance,
			Code:        item.Code,
			Authorising: item.Authorising,
		}

		if len(item.Storage) > 0 {
			data.Storage = translateToStorage(item.Storage)
		}

		target[i] = data
	}
	return
}

func translateToPoaMap(source *poaMap) (target *PoaMap) {
	target = &PoaMap{
		Address: source.Address,
		Balance: source.Balance,
		Abi:     source.ABI,
		SubAbi:  source.SubAbi,
		Code:    source.Code,
	}

	if len(source.Storage) > 0 {
		target.Storage = translateToStorage(source.Storage)
	}

	return
}

func translateToStorage(source [][2]string) (store map[string]string) {
	store = make(map[string]string, len(source))
	for _, tuple := range source {
		store[tuple[0]] = tuple[1]
	}

	return
}

func EncodeRLPGenesis(g *Genesis) ([]byte, error) {
	return rlp.EncodeToBytes(g)
}

func DecodeRLPGenesis(genesis []byte) (*Genesis, error) {
	var g Genesis
	err := rlp.DecodeBytes(genesis, &g)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func GetGenesisFromFile(filePath string) (*Genesis, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}

	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var gensis Genesis
	if err := viper.Unmarshal(&gensis); err != nil {
		return nil, err
	}

	return &gensis, nil
}

type alphabetic [][2]string

func (list alphabetic) Len() int      { return len(list) }
func (list alphabetic) Swap(i, j int) { list[i], list[j] = list[j], list[i] }
func (list alphabetic) Less(i, j int) bool {
	var (
		si      = list[i]
		sj      = list[j]
		siLower = strings.ToLower(si[0])
		sjLower = strings.ToLower(sj[0])
	)

	if siLower == sjLower {
		return si[0] < sj[0]
	}

	return siLower < sjLower
}
