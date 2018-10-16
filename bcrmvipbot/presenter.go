// Package main provides ...
package main

import "github.com/duminghui/go-tipservice/db"

var dbBcrm = db.NewDBSymbol("BCRM", "bcrm")
var dbGuild = db.NewDBGuild()

type presenter struct {
}

var psn = new(presenter)
