package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"strings"
	//"strconv"

	"strconv"

)

const (
	BLUE_COIN_VALUE_NONE = "none"
)

type BlockchainSample struct {
	Entity
	Name    string
	Age int32
	ErrMsg  string
}

type BlockchainSampleCompare func(BlockchainSample, BlockchainSample) (bool)

func createBlockchainSample(stub shim.ChaincodeStubInterface) BlockchainSample {
	bs := BlockchainSample{}
	e := Entity{Stub: stub, TableName: "BlockchainSample"}

	bs.Entity = e
	return bs
}

func (bs BlockchainSample) createTable() error {
	return createTable(bs.Entity, bs)
}

func (bs BlockchainSample) deleteTable() error {
	return deleteTable(bs.Entity)
}

func (bs BlockchainSample) insert() error {
	return insert(bs.Entity, bs)
}

func (bs BlockchainSample) replace() error {
	return replace(bs.Entity, bs)
}

func (bs BlockchainSample)delete()(error){
	return delete(bs.Entity, bs)
}

func (bs BlockchainSample) get() (BlockchainSample, error) {
	row, err := get(bs.Entity, bs)

	return bs.Convert(row), err
}

func (bs BlockchainSample) InitForCreate() []*shim.ColumnDefinition {
	return []*shim.ColumnDefinition{
		{Name: "UserID", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "UserMessage", Type: shim.ColumnDefinition_INT32, Key: false},
	}
}

func (bs BlockchainSample) InitForInsertAndReplace() []*shim.Column {
	return []*shim.Column{
		{Value: &shim.Column_String_{String_: bs.Name}},
		{Value: &shim.Column_Int32{Int32: bs.Age}},
	}
}

func (bs BlockchainSample) InitForDeleteAndGet() []shim.Column {
	return []shim.Column{
		{Value: &shim.Column_String_{String_: bs.Name}},
	}
}

func (bs BlockchainSample) Convert(row shim.Row) BlockchainSample {
	bs.Name = row.Columns[0].GetString_()
	bs.Age = row.Columns[1].GetInt32()
	return bs
}

func (bs BlockchainSample)generateKeyArray(num int8)([]shim.Column) {
	var col shim.Column

	switch(num){
	case 0:	return []shim.Column{}
	case 1:	col = shim.Column{Value: &shim.Column_String_{String_: bs.Name}}
	}

	return append(bs.generateKeyArray(num-1), col)
}

func (bs BlockchainSample)getAll(keyArr []shim.Column, bsCompare BlockchainSampleCompare)([]BlockchainSample, error) {
	var err error
	var rowChan <-chan shim.Row
	var bsArr []BlockchainSample

	rowChan, err = getAll(bs.Entity, bs, keyArr)

	if err != nil {
		return bsArr, err
	}

	for row := range rowChan {
		bsTemp :=  createBlockchainSample(bs.Stub)
		bsTemp = bsTemp.Convert(row)
		if bsCompare(bs, bsTemp) {
			bsArr = append(bsArr, bsTemp)
		}
	}

	return bsArr, nil
}

func (bs BlockchainSample)getAllByAge()([]BlockchainSample, error){
	var keyArr []shim.Column
	var bsCompare BlockchainSampleCompare

	keyArr = bs.generateKeyArray(0)

	bsCompare = func(bs1 BlockchainSample, bs2 BlockchainSample) (bool) {
		return bs1.Age == bs2.Age
	}

	bsArr, err := bs.getAll(keyArr, bsCompare)

	return bsArr, err
}


/*******************************

CHAINCODE HELPER FUNCTIONS

 *******************************/




func insertBlockchainSample(stub shim.ChaincodeStubInterface, args []string) (BlockchainSample, error) {
	bs := createBlockchainSample(stub)


	if len(args) != 2 {
		return bs, errors.New("Function insertBlockchainSample expects 2 arguments.")
	}

	name := args[0];
	age, _ := strconv.Atoi(args[1])

	bs.Name = name
	bs.Age = int32(age)

	err := bs.insert()
	if err != nil {
		bs.get()
		return bs, err
	}
	return bs, nil
}

func increaseAge(stub shim.ChaincodeStubInterface, args []string) (bsRet BlockchainSample, errRet error) {

	bs := createBlockchainSample(stub)


	if len(args) != 2 {
		return bs, errors.New("Function increaseAge expects 2 arguments.")
	}

	name := args[0];
	ageIncrease, _ := strconv.Atoi(args[1])

	bs.Name = name



	bs, err := bs.get()
	if err != nil {
		return bs, err
	}

	bs.Age = bs.Age + int32(ageIncrease)

	err = bs.replace()
	if err != nil {
		return bs, err
	}

	return bs, nil
}


func deleteBlockchainSample(stub shim.ChaincodeStubInterface, args []string) (bsRet BlockchainSample, errRet error) {

	bs := createBlockchainSample(stub)


	if len(args) != 1 {
		return bs, errors.New("Function deleteBlockchainSample expects 1 argument.")
	}

	name := args[0]

	bs.Name = name

	err := bs.delete()

	if err != nil {
		return bs, err
	}

	return bs, err
}


func getBlockchainSample(stub shim.ChaincodeStubInterface, args []string) (bsRet BlockchainSample, errRet error) {

	bs := createBlockchainSample(stub)


	if len(args) != 1 {
		return bs, errors.New("Function getBlockchainSample expects 1 argument.")
	}

	name := args[0]

	bs.Name = name

	bs, err := bs.get()

	if err != nil {
		return bs, err
	}

	return bs, err
}

func getAllBlockchainSampleByAge(stub shim.ChaincodeStubInterface, args []string) ([]BlockchainSample, error){
	var err error
	var bsArr []BlockchainSample

	bs := createBlockchainSample(stub)

	if len(args) != 1 {
		return bsArr, errors.New("Function getAllBlockchainSample expects 1 argument.")
	}



	age, _ := strconv.Atoi(args[0])

	bs.Age = int32(age)
	bsArr, err = bs.getAllByAge()

	return bsArr, err
}