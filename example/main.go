package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/whuanle/lsm"
	"github.com/whuanle/lsm/config"
	"net/http"
	"time"
)

type TestValue struct {
	A int64  `json:"a"`
	B int64  `json:"b"`
	C int64  `json:"c"`
	D string `json:"d"`
}

func main() {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	lsm.Start(config.Config{
		DataDir:       `./example/mygo`,
		Level0Size:    100,
		PartSize:      4,
		Threshold:     3000,
		CheckInterval: 3,
	})

	router := gin.Default()

	router.POST("/set", setHandler)
	router.POST("/search", searchHandler)

	router.Run(":8080")
}

func setHandler(c *gin.Context) {
	var request struct {
		Key   string    `json:"key"`
		Value TestValue `json:"value"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	lsm.Set(request.Key, request.Value)
	elapse := time.Since(start)

	c.JSON(http.StatusOK, gin.H{"message": "Insert completed", "elapsed_time": elapse.String()})
}

func searchHandler(c *gin.Context) {
	var request struct {
		Key string `json:"key"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	v, _ := lsm.Get[TestValue](request.Key)
	elapse := time.Since(start)

	c.JSON(http.StatusOK, gin.H{
		"key":          request.Key,
		"value":        v,
		"elapsed_time": elapse.String(),
	})
}

//
//import (
//	"bufio"
//	"fmt"
//	"github.com/whuanle/lsm"
//	"github.com/whuanle/lsm/config"
//	"os"
//	"time"
//)
//
//type TestValue struct {
//	A int64
//	B int64
//	C int64
//	D string
//}
//
//func main() {
//	defer func() {
//		r := recover()
//		if r != nil {
//			fmt.Println(r)
//			inputReader := bufio.NewReader(os.Stdin)
//			_, _ = inputReader.ReadString('\n')
//		}
//	}()
//	lsm.Start(config.Config{
//		DataDir:       `./example/mygo`,
//		Level0Size:    100,
//		PartSize:      4,
//		Threshold:     3000,
//		CheckInterval: 3,
//	})
//
//	insert()
//	query()
//
//}
//
//func query() {
//	start := time.Now()
//	v, _ := lsm.Get[TestValue]("aaaa")
//	elapse := time.Since(start)
//	fmt.Println("查找 aaa 完成，消耗时间：", elapse)
//	fmt.Println(v)
//
//	start = time.Now()
//	v, _ = lsm.Get[TestValue]("czzz")
//	elapse = time.Since(start)
//	fmt.Println("查找 aaz 完成，消耗时间：", elapse)
//	fmt.Println(v)
//}
//func insert() {
//
//	// 64 个字节
//	testV := TestValue{
//		A: 1,
//		B: 1,
//		C: 3,
//		D: "00000000000000000000000000000000000000",
//	}
//
//	//testVData, _ := json.Marshal(testV)
//	//// 131 个字节
//	//kvData, _ := kv.Encode(kv.Value{
//	//	Key:     "abcdef",
//	//	Value:   testVData,
//	//	Deleted: false,
//	//})
//	//fmt.Println(len(kvData))
//	//position := ssTable.Position{}
//	//// 35 个字节
//	//positionData, _ := json.Marshal(position)
//	//fmt.Println(len(positionData))
//	//
//	count := 0
//	start := time.Now()
//	key := []byte{'a', 'a', 'a', 'a'}
//	lsm.Set(string(key), testV)
//	for a := 0; a < 4; a++ {
//		for b := 0; b < 26; b++ {
//			for c := 0; c < 26; c++ {
//				for d := 0; d < 26; d++ {
//					key[0] = 'a' + byte(a)
//					key[1] = 'a' + byte(b)
//					key[2] = 'a' + byte(c)
//					key[3] = 'a' + byte(d)
//					lsm.Set(string(key), testV)
//					count++
//				}
//			}
//		}
//	}
//	elapse := time.Since(start)
//	fmt.Println("插入完成，数据量：", count, ",消耗时间：", elapse)
//}
