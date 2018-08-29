package service

import (
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"fmt"
	//"golang.org/x/net/html/atom"
)

// 发布票据
func (setup *FabricSetupService) SaveBill(bill Bill) (string, error) {
	var args []string
	args = append(args, "issue")
	b, _ := json.Marshal(bill)

	// 设置交易请求参数
	req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn: args[0], Args: [][]byte{b}}
	// 执行交易
	response, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("保存票据时发生错误: %v\n", err)
	}
	return response.TransactionID.ID, nil
}
// 根据当前持票人证件号码, 批量查询票据
func(setup *FabricSetupService) FindBills(holderCmId string) ([]byte, error) {

	req := chclient.Request{
		ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn: "queryBills",
		Args:[][]byte{[]byte(holderCmId)},
		}
	response,err := setup.Fabric.Client.Query(req)
	if err != nil {
		return nil , fmt.Errorf("%s", err.Error())
	}
	return response.Payload,nil
}

//根据票据编号查询票据流转记录
func(setup *FabricSetupService) QueryBillByNo(billInfoID string) ([]byte, error) {
	req := chclient.Request{
		ChaincodeID:setup.Fabric.ChaincodeID,
		Fcn:"queryBillByNo",
		Args:[][]byte{[]byte(billInfoID)},
	}
	response,err := setup.Fabric.Client.Query(req)
	if err != nil {
		return nil , fmt.Errorf("%s", err.Error())
	}
	return response.Payload,nil
}

//发起背书
func(setup *FabricSetupService) Endorse(billInfoID,waitEndorseAcct,waitEndorseCmId string) (string, error) {
	req := chclient.Request{
		ChaincodeID:setup.Fabric.ChaincodeID,
		Fcn:"endorse",
		Args:[][]byte{[]byte(billInfoID),[]byte(waitEndorseAcct),[]byte(waitEndorseCmId)},
	}

	response,err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	return response.TransactionID.ID,nil
}


//查询待背书列表
func(setup *FabricSetupService) QueryWaitBills(waitEndorseCmId string) ([]byte, error)  {
	req := chclient.Request{
		ChaincodeID:setup.Fabric.ChaincodeID,
		Fcn:"queryWaitBills",
		Args:[][]byte{[]byte(waitEndorseCmId)},
	}
	response,err := setup.Fabric.Client.Query(req)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	//fmt.Println("response.Payload",response.Payload)
	return response.Payload,nil
}

// 背书签名
func (setup *FabricSetupService) Accept(billInfoID,accptrAcct,accptrCmId string) (string,error) {

	req := chclient.Request{
		ChaincodeID:setup.Fabric.ChaincodeID,
		Fcn:"accept",
		Args:[][]byte{[]byte(billInfoID),[]byte(accptrAcct),[]byte(accptrCmId)},
	}
	response,err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	return response.TransactionID.ID,nil
}

// 拒绝背书
func (setup *FabricSetupService) Reject(billInfoID,rejectrAcct,rejectrCmId string) (string,error) {

	req := chclient.Request{
		ChaincodeID:setup.Fabric.ChaincodeID,
		Fcn:"reject",
		Args:[][]byte{[]byte(billInfoID),[]byte(rejectrAcct),[]byte(rejectrCmId)},
	}

	response,err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	return response.TransactionID.ID,nil
}







