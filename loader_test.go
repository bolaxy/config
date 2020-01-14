package conf

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var (
	testFilePath = "./testdata"
)

func TestTryLoadGenesis(t *testing.T) {
	genesis, err := TryLoadGenesis(testFilePath)
	if err != nil {
		t.Fatalf("failed to load genesis file. cause: %v\n", err)
	}

	h1, err := genesis.Hash()
	if err != nil {
		t.Fatalf("failed to get genesis hash, cause: %v\n", err)
	}

	h2, _ := genesis.Hash()
	if !bytes.Equal(h1, h2) {
		t.Fatal("some genesis get difference hash")
	}
	// for k, v := range genesis.Allocs {
	//	fmt.Println("k:", k, "v:", v)
	// }
	spew.Dump(genesis)
}

var genesisNames = []string{
	"genesis",
	"genesis_all",
	"genesis_no_launcher",
}

func TestEncodeGenesis(t *testing.T) {
	for _, name := range genesisNames {
		runEncodeGenesis(t, name)
	}
}

func runEncodeGenesis(t *testing.T, name string) {
	genesis, err := TryLoadGenesis(testFilePath, name)
	if err != nil {
		t.Fatalf("failed to load genesis file. cause: %v\n", err)
	}

	h1, err := genesis.Hash()
	if err != nil {
		t.Fatalf("failed to get genesis hash, cause: %v\n", err)
	}

	h2, _ := genesis.Hash()
	if !bytes.Equal(h1, h2) {
		t.Fatal("some genesis get difference hash")
	}

	buf := bytes.NewBuffer(nil)
	if err = genesis.EncodeRLP(buf); err != nil {
		t.Fatalf("failed to encode genesis. cause: %v\n", err)
	}

	t.Logf("RLP: %v\n", buf.Bytes())
}

func TestDecodeGenesis(t *testing.T) {
	for _, name := range genesisNames {
		runDecodeGenesis(t, name)
	}
}

func runDecodeGenesis(t *testing.T, name string) {
	t.Logf("running in %s", name)
	genesis, err := TryLoadGenesis(testFilePath, name)
	if err != nil {
		t.Fatalf("failed to load genesis file. cause: %v\n", err)
	}

	spew.Println("------before-----------------------------------------------------------------")
	spew.Dump(genesis)
	spew.Println("------before-----------------------------------------------------------------")

	buf := bytes.NewBuffer(nil)
	if err = genesis.EncodeRLP(buf); err != nil {
		t.Fatalf("failed to encode genesis. cause: %v\n", err)
	}

	g, err := DecodeRLPGenesis(buf.Bytes())
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	spew.Println("------after------------------------------------------------------------------")
	spew.Dump(g)
	spew.Println("------after------------------------------------------------------------------")
}

func TestDecodeConfig(t *testing.T) {
	TryLoadConfig(testFilePath)
	// fmt.Println("---:", Global.Eth.CacheSize)
	spew.Dump(Global)
}
