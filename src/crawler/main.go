/*
 * Created  main.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午2:55
 */

package main

import (
	"crawler/src/crawler/c"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {

	go c.CrawlerProduct()
	go c.CrawlerWishId()
	go c.FeedCrawler()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		exec_shell("update.sh")
	})

	s := &http.Server{
		Addr: ":9529",
	}

	log.Fatal(s.ListenAndServe())
}

func exec_shell(s string) {
	cmd := exec.Command("sh", s)

	_, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output: ", err)
		return
	}
	os.Exit(1)
}
