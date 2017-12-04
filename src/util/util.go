/*
 * Created  util.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午4:52
 */

package util

import (
	"hash/fnv"
	"time"
)

/*
循环任务定时器：
hour : 每天几点启动任务
f: 任务函数
*/
func LoopTimer(hour int, f func()) {
	go func() {
		for {
			f()
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), hour, 0, 0, 0, time.Local)
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

func FNV(s string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return hash.Sum32()
}
