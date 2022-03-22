package main

import (
	"testing"

	db "github.com/squirlyfoxy/ronny/database"
)

func TestContains(t *testing.T) { //Strings
	s := []string{"a", "b", "c"}
	if !db.Contains(s, "a") {
		t.Error("Contains() failed")
	}
	if db.Contains(s, "d") {
		t.Error("Contains() failed")
	}
}

func TestContainsOnType(t *testing.T) {
	s := []db.OnType{db.CAN_GLOBALLY_TAKE, db.CAN_TAKE}
	if !db.ContainsOnType(s, db.CAN_GLOBALLY_TAKE) {
		t.Error("ContainsOnType() failed")
	}
	if db.ContainsOnType(s, db.CAN_REMOVE) {
		t.Error("ContainsOnType() failed")
	}
}

func TestCheckTableRule(t *testing.T) {
	table := db.Table{
		Name: "test",
		Rule: db.Rule{
			RuleTypes: []db.RuleType{
				{
					Can: []db.OnType{db.CAN_GLOBALLY_TAKE},
				},
			},
		},
	}
	if !db.CheckTableRule(table, db.CAN_GLOBALLY_TAKE) {
		t.Error("CheckTableRule() failed")
	}
	if db.CheckTableRule(table, db.CAN_REMOVE) {
		t.Error("CheckTableRule() failed")
	}
}

func TestRemove(t *testing.T) { //Strings
	s := []string{"a", "b", "c"}
	if len(db.Remove(s, "a")) != 2 {
		t.Error("Remove() failed")
	}
	if len(db.Remove(s, "d")) != 3 {
		t.Error("Remove() failed")
	}
}
