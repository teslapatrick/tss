package common_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	. "github.com/binance-chain/tss/common"
	"github.com/stretchr/testify/assert"
)

// example.json and node_key required within this file is generated by run setup_lan.sh
// with dump each peer's localPartySaveData in client.go as:
// 		bb, _ := json.Marshal(&msg)
//		err = ioutil.WriteFile(path.Join(client.config.Home, "example.json"), bb, 0400)

func TestSaveAndLoad(t *testing.T) {
	var wPriv bytes.Buffer
	var wPub bytes.Buffer
	passphrase := "1234qwerasdf"

	var expectedMsg keygen.LocalPartySaveData
	localSaveDataBytes, err := ioutil.ReadFile("example.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(localSaveDataBytes, &expectedMsg)
	if err != nil {
		t.Fatal(err)
	}
	expectedNodeKey, err := ioutil.ReadFile("node_key")
	if err != nil {
		t.Fatal(err)
	}

	Save(&expectedMsg, expectedNodeKey, DefaultKDFConfig(), passphrase, &wPriv, &wPub)
	result, nodeKey, _ := Load(passphrase, &wPriv, &wPub)

	if !reflect.DeepEqual(*result, expectedMsg) {
		t.Fatal("local saved data is not expected")
	}

	if !bytes.Equal(nodeKey, expectedNodeKey) {
		t.Fatal("local saved node key is not expected")
	}
}

func TestSaveAndLoadNeg(t *testing.T) {
	var wPriv bytes.Buffer
	var wPub bytes.Buffer
	passphrase := "1234qwerasdf"

	var expectedMsg keygen.LocalPartySaveData
	bytes, err := ioutil.ReadFile("example.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(bytes, &expectedMsg)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = ioutil.ReadFile("node_key")
	if err != nil {
		t.Fatal(err)
	}

	Save(&expectedMsg, bytes, DefaultKDFConfig(), passphrase, &wPriv, &wPub)
	_, _, err = Load("12345678", &wPriv, &wPub) // load saved data with a wrong passphrase would not success
	assert.Error(t, err)
}
