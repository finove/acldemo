package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

func main() {
	var sub = "admin"
	var obj = "book"
	var act = "read"
	e, err := casbin.NewEnforcer("data/model.conf", "data/policy.csv")
	log.Printf("new enforcer %v", err)
	if res, _ := e.Enforce(sub, obj, act); res {
		fmt.Printf("action ok\n")
	} else {
		fmt.Printf("action deny\n")
	}
	fmt.Printf("done\n")
}
