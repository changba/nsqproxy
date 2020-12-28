package model

import (
	"encoding/json"
	"testing"
)

func TestNewDb(t *testing.T) {
	NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
	type version struct{
		Version string
	}
	v := version{}
	result := db.Raw("SELECT @@version AS version").Scan(&v)
	if result.Error != nil{
		t.Fatalf("error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1{
		t.Fatalf("RowsAffected != 1")
	}
	if len(v.Version) <= 0{
		t.Fatalf("v.Version is empty")
	}
	t.Log(v.Version)
}

func TestGetAvailableConsumeList(t *testing.T) {
	NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
	consumeConfigList, err := GetAvailableConsumeList()
	if err != nil{
		t.Fatalf("error: %s", err.Error())
	}
	j, err := json.Marshal(consumeConfigList)
	if err != nil{
		t.Fatalf("error: %s", err.Error())
	}
	if len(j) < 100{
		t.Fatalf("json is short")
	}
}