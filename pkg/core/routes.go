package core

import (
	"github.com/aioncore/ouroboros/pkg/core/middleware"
	"github.com/aioncore/ouroboros/pkg/service/server/utils"
)

func GenerateRoutes(sw middleware.Switch) map[string]*utils.APIFunc {
	routes := map[string]*utils.APIFunc{
		// info API
		//"health":       NewAPIFunc(sw.Health, ""),
		"services":     utils.NewAPIFunc(sw.Services, ""),
		"use_service":  utils.NewAPIFunc(sw.UseService, "service_data"),
		"call_service": utils.NewAPIFunc(sw.CallService, "type,request"),
		//"status":               NewAPIFunc(Status, ""),
		//"net_info":             NewAPIFunc(NetInfo, ""),
		//"blockchain":           NewAPIFunc(BlockchainInfo, "minHeight,maxHeight"),
		//"genesis":              NewAPIFunc(Genesis, ""),
		//"genesis_chunked":      NewAPIFunc(GenesisChunked, "chunk"),
		//"block":                NewAPIFunc(Block, "height"),
		//"block_by_hash":        NewAPIFunc(BlockByHash, "hash"),
		//"block_results":        NewAPIFunc(BlockResults, "height"),
		//"commit":               NewAPIFunc(Commit, "height"),
		//"check_tx":             NewAPIFunc(CheckTx, "tx"),
		//"tx":                   NewAPIFunc(Tx, "hash,prove"),
		//"tx_search":            NewAPIFunc(TxSearch, "query,prove,page,per_page,order_by"),
		//"block_search":         NewAPIFunc(BlockSearch, "query,page,per_page,order_by"),
		//"validators":           NewAPIFunc(Validators, "height,page,per_page"),
		//"dump_consensus_state": NewAPIFunc(DumpConsensusState, ""),
		//"consensus_state":      NewAPIFunc(ConsensusState, ""),
		//"consensus_params":     NewAPIFunc(ConsensusParams, "height"),
		//"unconfirmed_txs":      NewAPIFunc(UnconfirmedTxs, "limit"),
		//"num_unconfirmed_txs":  NewAPIFunc(NumUnconfirmedTxs, ""),
		//
		//// tx broadcast API
		//"broadcast_tx_commit": NewAPIFunc(BroadcastTxCommit, "tx"),
		//"broadcast_tx_sync":   NewAPIFunc(BroadcastTxSync, "tx"),
		//"broadcast_tx_async":  NewAPIFunc(BroadcastTxAsync, "tx"),
		//
		//// abci API
		//"abci_query": NewAPIFunc(ABCIQuery, "path,data,height,prove"),
		//"abci_info":  NewAPIFunc(ABCIInfo, ""),
		//
		//// evidence API
		//"broadcast_evidence": NewAPIFunc(BroadcastEvidence, "evidence"),
	}
	return routes
}
