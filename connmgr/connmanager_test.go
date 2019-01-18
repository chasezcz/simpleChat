package connmgr

import "testing"

func TestConnManager_Init(t *testing.T) {
	var connManager ConnManager
	ok := connManager.Init("192.168.12.1")

	if !ok {
		t.Error("错误")
	}

	

}