package main

import (
	"ZETA_server/ZetaDB"
	"fmt"
)

func main() {
	fmt.Println(ZetaDB.ExecuteSql(0, "initialize;"))
}
