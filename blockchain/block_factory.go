
package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

func CreateGenesisBlock(genesisconfFilePath string) (DefaultBlock, error) {

	//declare
	GenesisBlock := &DefaultBlock{}
	validator := DefaultValidator{}

	//set basic
	err := setBlockWithConfig(genesisconfFilePath, GenesisBlock)

	if err != nil {
		return DefaultBlock{}, ErrSetConfig
	}

	//build
	Seal, err := validator.BuildSeal(GenesisBlock.Timestamp, GenesisBlock.PrevSeal, GenesisBlock.TxSeal, GenesisBlock.Creator)

	if err != nil {
		return DefaultBlock{}, ErrBuildingSeal
	}

	//set seal
	GenesisBlock.SetSeal(Seal)

	return *GenesisBlock, nil
}

func setBlockWithConfig(filePath string, block *DefaultBlock) error {

	// load
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()

	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	GenesisConfig := &GenesisConfig{}

	err = json.Unmarshal(byteValue, GenesisConfig)
	if err != nil {
		return err
	}

	// set
	const longForm = "Jan 1, 2006 at 0:00am (MST)"

	timeStamp, err := time.Parse(longForm, GenesisConfig.TimeStamp)

	if err != nil {
		return err
	}

	block.SetPrevSeal(make([]byte, 0))
	block.SetHeight(uint64(GenesisConfig.Height))
	block.SetTxSeal(make([][]byte, 0))
	block.SetTimestamp(timeStamp)
	block.SetCreator(GenesisConfig.Creator)
	block.SetState(Created)

	return nil
}

type GenesisConfig struct {
	Organization string
	NedworkId    string
	Height       int
	TimeStamp    string
	Creator      string
}

func CreateProposedBlock(prevSeal []byte, height uint64, txList []*DefaultTransaction, Creator string) (DefaultBlock, error) {

	//declare
	ProposedBlock := &DefaultBlock{}
	validator := DefaultValidator{}
	TimeStamp := time.Now().Round(0)

	//build
	for _, tx := range txList {
		ProposedBlock.PutTx(tx)
	}

	txSeal, err := validator.BuildTxSeal(ConvertTxType(txList))

	if err != nil {
		return DefaultBlock{}, ErrBuildingTxSeal
	}

	Seal, err := validator.BuildSeal(TimeStamp, prevSeal, txSeal, Creator)

	if err != nil {
		return DefaultBlock{}, ErrBuildingSeal
	}

	//set
	ProposedBlock.SetSeal(Seal)
	ProposedBlock.SetPrevSeal(prevSeal)
	ProposedBlock.SetHeight(height)
	ProposedBlock.SetTxSeal(txSeal)
	ProposedBlock.SetTimestamp(TimeStamp)
	ProposedBlock.SetCreator(Creator)
	ProposedBlock.SetState(Created)

	return *ProposedBlock, nil
}
