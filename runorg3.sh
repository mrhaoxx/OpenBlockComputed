#!/bin/bash

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=/home/star/fabric-samples/test-network/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=/home/star/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
export PEER0_ORG2_CA=/home/star/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
export PEER0_ORG3_CA=/home/star/fabric-samples/test-network/organizations/peerOrganizations/org3.example.com/tlsca/tlsca.org3.example.com-cert.pem
export CORE_PEER_MSPCONFIGPATH=/home/star/fabric-samples/test-network/organizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
export CORE_PEER_ADDRESS=dns:///localhost:11051
export CORE_PEER_TLS_ROOTCERT_FILE=/home/star/fabric-samples/test-network/organizations/peerOrganizations/org3.example.com/tlsca/tlsca.org3.example.com-cert.pem
export CORE_PEER_LOCALMSPID=Org3MSP
export CORE_GATEWAY_PEER=peer0.org3.example.com
export PORT="8082"
go run .