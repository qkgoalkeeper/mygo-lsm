package lsm

import (
	"github.com/whuanle/lsm/skipList"
	"github.com/whuanle/lsm/ssTable"
	"github.com/whuanle/lsm/wal"
)

type Database struct {
	// 内存表
	//MemoryTree *sortTree.Tree
	MemoryTree *skipList.SkipList
	// SSTable 列表
	TableTree *ssTable.TableTree
	// WalF 文件句柄
	Wal *wal.Wal
}

// 数据库，全局唯一实例
var database *Database
