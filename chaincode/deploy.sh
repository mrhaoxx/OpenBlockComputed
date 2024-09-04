cd /home/star/fabric-samples/test-network

./network.sh down
./network.sh up createChannel -ca -c mychannel 
# ./addOrg3/addOrg3.sh up
./network.sh deployCC -ccn openbc -ccp ~/OpenBlockComputed/chaincode -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cccg '/home/star/OpenBlockComputed/chaincode/collections_config.json'  -ccep "OR('Org1MSP.peer','Org2MSP.peer')"


# export CORE_PEER_TLS_ENABLED=true
# export CORE_PEER_LOCALMSPID=Org1MSP
# export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
# export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
# export CORE_PEER_ADDRESS=localhost:7051
# export PATH=${PWD}/../bin:$PATH
# export FABRIC_CFG_PATH=${PWD}/../config

# peer chaincode query -C mychannel -n openbc -c '{"Args":["org.hyperledger.fabric:GetMetadata"]}' | jq

# peer chaincode invoke -C mychannel -n openbc -c  '{"Args":["CreateRootUser"]}' | jq

# ./network.sh cc invoke -org 2 -c mychannel -ccn openbc -ccic '{"Args":["CreateRootUser"]}'

# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n openbc --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"CreateRootUser","Args":[]}'
