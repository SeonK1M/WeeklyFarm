package main

// import modules

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// define class
type SmartContract struct{}

type Farmland struct {
	FarmID      string `json:"fid"`
	BusinessID  string `json:"bid"`
	Address     string `json:"address"`
	Participant string `json:"participant"`
	RentPrice   string `json:"rentPrice"`
}

type Record struct {
	ContractID string `json:"cid"`
	UserID     string `json:"uid"`
	FarmID     string `json:"fid"`
	PlantName  string `json:"plantName"`
	PlantState int    `json:"plantState"`
} //웹 서비스를 기획할 때 랜딩페이지와 블록체인 대시보드는 필수적임

// func init
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response { //모든 데이터를 초기화하기 위해 체인 코드 인스턴스화 중에 Init이 호출됨
	//when you want to instantiate or upgrade chaincode
	return shim.Success(nil)
}

// func invoke
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response { //트랜잭션 요청을 받을 때 호출되는 함수, (stub(인자 이름), shim.CSIf(자료형태))
	fn, args := stub.GetFunctionAndParameters()

	if fn == "registerFarmland" {
		return s.registerFarmland(stub, args)
	} else if fn == "newRecord" {
		return s.newRecord(stub, args)
	} else if fn == "updateRecord" {
		return s.updateRecord(stub, args)
	} else if fn == "getAllFarmland" {
		return s.getAllFarmland(stub, args)
	} else if fn == "getRecordInfo" {
		return s.getRecordInfo(stub, args)
	}
	return shim.Error("wrong function name")
}

// registerFarmland
// newContract
// updateRecord
func (s *SmartContract) registerFarmland(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// check args
	if len(args) != 4 { //4개 파라미터 받는다
		return shim.Error("registerFarmland: wrong args")
	}
	var farm = Farmland{FarmID: args[0], BusinessID: args[1], Address: args[2], RentPrice: args[3]}
	farmAsBytes, _ := json.Marshal(farm)
	stub.PutState(args[0], farmAsBytes) //마샬한 데이터를 키와 밸류 형태로 월드스테이트에 저장한다.

	return shim.Success([]byte("transaction has been submited"))
} // record.Plantstate +1 (조건문 써주기) // const type - req.body.type; (post) // if type == 1 // register // else (type) 이렇게 하지 말고 걍 하나 더 만드셈

func (s *SmartContract) newRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// check args
	if len(args) != 4 { //파라미터 4개인지 확인
		return shim.Error("newRecord: wrong args")
	}
	var record = Record{ContractID: args[0], UserID: args[1], FarmID: args[2], PlantName: args[3], PlantState: 0} //파라미터 받는 변수 선언
	recordAsBytes, _ := json.Marshal(record)                                                                      //마샬함
	stub.PutState(args[0], recordAsBytes)                                                                         //마샬한 데이터를 키와 밸류 형태로 월드스테이트에 저장한다.

	return shim.Success([]byte("transaction has been submited"))
} // 새 계약을 등록할 때 등록된 토지인지 확인하는 부분 추가하기

func (s *SmartContract) updateRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// check args
	if len(args) != 1 { //파라미터 3개, 계약(일지)id
		return shim.Error("updateRecord: wrong args")
	}
	recordAsBytes, err := stub.GetState(args[0]) //
	if err != nil {
		return shim.Error("function recordScore: error from GetState")
	} // if getstate bring error
	if recordAsBytes == nil { //if there's no id
		return shim.Error("no ID")
	}

	var record = Record{}
	_ = json.Unmarshal(recordAsBytes, &record) //언마샬함

	if record.PlantState < 5 {
		record.PlantState++
	}
	recordAsBytes, _ = json.Marshal(record) //데이터를 마샬해서
	stub.PutState(args[0], recordAsBytes)   //월드스테이트에

	return shim.Success([]byte("transaction has been submited"))
}

// getAllFarmLand
// getRecordInfo
// getCurrentStateOfPlant
func (s *SmartContract) getAllFarmland(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	startKey := "FARM0"
	endKey := "FARM999"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"value\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getAllFarmland:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) getRecordInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 { //파라미터 하나
		return shim.Error("function getRecordInfo: wrong args")
	}
	info, _ := stub.GetState(args[0])
	return shim.Success(info)
}

//main
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("Error while creating new SmartContract")
	}
}
