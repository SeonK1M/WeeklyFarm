#!/bin/bash

docker exec cli peer chaincode install -n farm -v 0.9 -p github.com/farm

docker exec cli peer chaincode instantiate -n farm -v 0.9 -C mychannel -c '{"Args":[]}' -P 'OR("Org1MSP.member")'
sleep 3

docker exec cli peer chaincode invoke -n farm  -C mychannel -c '{"Args":["registerFarmland","FARM1","01-1111111","111-111","1,200,000"]}'
echo registering new farmland : FARM1
sleep 3

docker exec cli peer chaincode invoke -n farm  -C mychannel -c '{"Args":["newRecord","CID-1111","USER1","FARM1","TOMATO"]}'
echo new contract : CID-1111
sleep 3

docker exec cli peer chaincode invoke -n farm  -C mychannel -c '{"Args":["updateRecord","CID-1111"]}'
echo update contract : CID-1111
sleep 3

docker exec cli peer chaincode query -n farm  -C mychannel -c '{"Args":["getRecordInfo","CID-1111"]}'
sleep 3

docker exec cli peer chaincode query -n farm  -C mychannel -c '{"Args":["getAllFarmland"]}'
sleep 3