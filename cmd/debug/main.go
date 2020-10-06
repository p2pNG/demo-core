package main

import (
	"encoding/json"
	"fmt"
	"git.ixarea.com/p2pNG/p2pNG-core/components/file_store"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	//F:/Installer/System/cn_windows_10_consumer_editions_version_2004_updated_sep_2020_x64_dvd_049d70ee.iso
	data, err := file_store.StatLocalFile("D:/temp/bank-proj/Release/package", 0)
	spew.Dump(data, err)
	dataJson, err := json.Marshal(data)
	fmt.Print(string(dataJson))
}
