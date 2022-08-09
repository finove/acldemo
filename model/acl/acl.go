package acl

import (
	"github.com/casbin/casbin/v2"
)

var (
	e *casbin.Enforcer
)

func SetupEnforcer() (err error) {
	e, err = casbin.NewEnforcer("data/model.conf", "data/policy.csv")
	return
}

func VerifyAccess(sub, obj, act string) (yes bool) {
	yes, _ = e.Enforce(sub, obj, act)
	return
}

func TestEnforce() {
	var sub = "admin"
	var obj = "book"
	var act = "read"
	VerifyAccess(sub, obj, act)
}
