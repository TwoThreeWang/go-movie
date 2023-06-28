package easylog

import "testing"

func TestInit(t *testing.T) {
	Init()
	Log.Info(`测试日志Info`)
	Log.Warn(`测试日志Warn`)
	Log.Error(`测试日志Error`)
	Log.Debug(`测试日志Debug`)
}
