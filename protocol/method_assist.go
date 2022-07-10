package protocol

import (
	"ZETA_server/ZetaDB"
	"ZETA_server/ZetaDB/parser"
)

//assist function
func (n *Node) Check_AddressIncluded(address string) bool {

	for _, bufferAddr := range n.AddressBuffer {
		if bufferAddr == address {
			return true
		}
	}

	for _, heartbeatAddr := range n.AddressInMessage {
		if heartbeatAddr == address {
			return true
		}
	}

	for _, logAddr := range n.AddressLog {
		if logAddr == address {
			return true
		}
	}
	return false
}

func (n *Node) Check_SqlReadOnly(sql string) bool {

	//get sql parser
	p := ZetaDB.GetParser()

	ast, _ := p.ParseSql(sql)

	if ast.Type == parser.AST_DQL {
		return true
	}

	if ast.Type == parser.AST_DCL {

		dclType := ast.Dcl.Type

		if dclType == parser.DCL_SHOW_TABLES ||
			dclType == parser.DCL_SHOW_ASSERTIONS ||
			dclType == parser.DCL_SHOW_VIEWS ||
			dclType == parser.DCL_SHOW_INDEXS ||
			dclType == parser.DCL_SHOW_TRIGGERS ||
			dclType == parser.DCL_SHOW_FUNCTIONS ||
			dclType == parser.DCL_SHOW_PROCEDURES {
			return true
		}
	}

	return false
}
