package main

import (
	"ZETA_server/ZetaDB"
	"fmt"
)

func main() {
	db := ZetaDB.GetDatabase("")

	fmt.Println(db.ExecuteSql(0, "initialize;"))
}
