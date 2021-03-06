package ZetaDB

import (
	"ZETA_server/ZetaDB/execution"
	parser "ZETA_server/ZetaDB/parser"
	"ZETA_server/ZetaDB/storage"
	"ZETA_server/ZetaDB/utility"
	"os"
	"sync"
)

type Database struct{}

var dbInstance *Database
var dbOnce sync.Once

//call this function to get database
func GetDatabase(fileLocation string) *Database {

	//alter file location
	if fileLocation == "" { // set file location to defaut
		currentLocation, _ := os.Getwd()

		utility.DEFAULT_DATAFILE_LOCATION = currentLocation + "/ZetaDB/file/data.zdb"
		utility.DEFAULT_INDEXFILE_LOCATION = currentLocation + "/ZetaDB/file/index.zdb"
		utility.DEFAULT_LOGFILE_LOCATION = currentLocation + "/ZetaDB/file/log.zdb"
	} else {
		utility.DEFAULT_DATAFILE_LOCATION = fileLocation + "/data.zdb"
		utility.DEFAULT_INDEXFILE_LOCATION = fileLocation + "/index.zdb"
		utility.DEFAULT_LOGFILE_LOCATION = fileLocation + "/log.zdb"
	}

	dbOnce.Do(func() {
		dbInstance = &Database{}
	})
	return dbInstance
}

func (db *Database) ExecuteSql(currentUserId int32, sqlString string) (int32, string) {
	//get Parser, rewriter, executionEngine, transaction
	Parser := parser.GetParser()
	rewriter := execution.GetRewriter()
	executionEngine := execution.GetExecutionEngine()
	transaction := storage.GetTransaction()

	//parse this sql and get an AST, if sql syntax invalid, reply immediately
	sqlAstNode, parseErr := Parser.ParseSql(sqlString)
	if parseErr != nil {
		return -1, "error: sql syntax invalid"
	}

	//TODO unfinished, change userId
	//generate an executionPlan from current userId, AST and sql string
	executionPlan, rewriteErr := rewriter.ASTNodeToExecutionPlan(currentUserId, sqlAstNode, sqlString)
	if rewriteErr != nil {
		return -1, rewriteErr.Error()
	}

	//TODO debug
	if executionPlan == nil {
		return -1, "error: not supported currently"
	}

	//use executionEngine to execute this executionPlan, get a result string for reply
	executionResult := executionEngine.Execute(executionPlan)

	//update all modification into disk
	transaction.PushTransactionIntoDisk()

	//reply
	if len(executionResult) > 8 && executionResult[0:8] == "userId: " {
		//return logged userId
		return 1, executionResult[8:]
	} else if len(executionResult) > 7 && executionResult[0:7] == "error: " {
		return -1, executionResult
	} else if executionResult == "Execute OK, system halt" {
		return -2, executionResult
	} else {
		return 0, executionResult
	}
}

func GetParser() *parser.Parser {
	return parser.GetParser()
}
