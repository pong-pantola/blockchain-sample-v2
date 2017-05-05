package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type BlockchainSampleChaincode struct {
}


type Attribute struct {
	AttributeName string
	AttributeVal string
	ErrMsg string
}

/*
args:
[0] - Attribute Name
 */
func (*BlockchainSampleChaincode) getAttribute(stub shim.ChaincodeStubInterface, args []string) (Attribute, error){
	var attr Attribute
	var err error

	if len(args) != 1 {
		return attr, errors.New("Function getAttribute expects 1 argument.")
	}


	index := -1

	index++
	attributeName := formatInput(args[index])

	attributeVal, err := getCertAttribute(stub, attributeName)
	if (err != nil){
		err = errors.New("getAttribute cannot get the attribute ["+ attributeName + "]")

		return attr, err
	}

	//attr  = Attribute{AttributeName: attributeName, AttributeVal: attributeVal}
	attr = Attribute{}
	attr.AttributeName = attributeName
	attr.AttributeVal = attributeVal
	return attr, nil
}

func (*BlockchainSampleChaincode) deleteTable(stub shim.ChaincodeStubInterface) (error){
	bs := createBlockchainSample(stub)

	_ = bs.deleteTable()

	return nil
}

func (*BlockchainSampleChaincode) createTable(stub shim.ChaincodeStubInterface) (error){
	var err error

	bs := createBlockchainSample(stub)

	err = bs.createTable()
	if err != nil{
		return err
	}

	return nil
}

/*
Init is called when chaincode is deployed.
*/
func (cc *BlockchainSampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	if function == "createTable" {

		err = cc.createTable(stub)

		return nil, err

	}

	return nil, errors.New("Unknown function " + function + ".")
}

func (cc *BlockchainSampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	if function == "getAttribute" {
		var attr Attribute

		attr, err = cc.getAttribute(stub, args)

		if err != nil {
			attr.ErrMsg = err.Error()
			return formatOutput(attr)
		}

		return formatOutput(attr)
	}else
	
	if function == "getBlockchainSample" { 
		var bs BlockchainSample

		bs, err = getBlockchainSample(stub, args)

		if err != nil {
			bs.ErrMsg = err.Error()
			return formatOutput(bs)
		}

		//jsonObj, err := json.Marshal(bs)

		//return []byte(string(jsonObj)), err
		return formatOutput(bs)
	}else
	if function == "getAllBlockchainSampleByAge" {
		var bsArr []BlockchainSample

		bsArr, err = getAllBlockchainSampleByAge(stub, args)

		if err != nil {
			return nil, err
		}

		//jsonObj, err := json.Marshal(icArr)

		return formatOutput(bsArr)
	}
	
	return nil, errors.New("Unknown function " + function + ".")
}

func (cc *BlockchainSampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	if function == "deleteTable" {
		err = cc.deleteTable(stub)

		return nil, err
	}else
	if function == "insertBlockchainSample" {
		var bs BlockchainSample

		bs, err = insertBlockchainSample(stub, args)

		if err != nil {
			bs.ErrMsg = err.Error();
			stub.SetEvent("insertBlockchainSample", formatPayload(bs))
			return nil, err
		}

		err = stub.SetEvent("insertBlockchainSample", formatPayload(bs))
		if err != nil {
			return nil, err
		}
		return nil, nil

	}else
	if function == "increaseAge" {
		var bs BlockchainSample

		bs, err = increaseAge(stub, args)

		if err != nil {
			bs.ErrMsg = err.Error();
			stub.SetEvent("increaseAge", formatPayload(bs))
			return nil, err
		}

		err = stub.SetEvent("increaseAge", formatPayload(bs))
		if err != nil {
			return nil, err
		}
		return nil, nil

	}else
	if function == "deleteBlockchainSample" {
		var bs BlockchainSample

		bs, err = deleteBlockchainSample(stub, args)

		if err != nil {
			bs.ErrMsg = err.Error();
			stub.SetEvent("deleteBlockchainSample", formatPayload(bs))
			return nil, err
		}

		err = stub.SetEvent("deleteBlockchainSample", formatPayload(bs))
		if err != nil {
			return nil, err
		}
		return nil, nil

	}

	return nil, nil
}

func main() {
	err := shim.Start(new(BlockchainSampleChaincode))
	if err != nil {
		fmt.Printf("Error creationing BlockchainSampleChaincode: %s", err)
	}
}