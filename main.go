package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"abtest/model"

	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化一些示例配置 / Initialize some sample configurations
	initSampleConfig()

	// 设置HTTP路由 / Set up HTTP routes
	http.HandleFunc("/abtest", ABInterface)

	// 启动HTTP服务器 / Start HTTP server
	logrus.Info("Starting ABTest server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logrus.Fatal("Server failed to start: ", err)
	}
}

// ABInterface HTTP处理函数 / HTTP handler function
func ABInterface(w http.ResponseWriter, req *http.Request) {
	tm := time.Now().UnixNano()

	head := 0
	device_id := req.URL.Query().Get("device_id")
	if device_id == "" {
		device_id = req.Header.Get("device_id")
		if device_id != "" {
			head = 1
		}
	}

	res := model.GetAbId(device_id)
	r, _ := json.Marshal(res)

	logrus.Infof("[ABInterface] device_id:%v, head:%v, tm:%v us, res:%v", device_id, head, (time.Now().UnixNano()-tm)/1000, res.Data.AbId)

	w.Header().Set("content-type", "application/json")
	io.WriteString(w, string(r))
}

// initSampleConfig 初始化示例配置 / Initialize sample configuration
func initSampleConfig() {
	// 创建一些示例AB测试配置 / Create some sample AB test configurations
	configs := []*model.AbTestConfig{
		{
			AbId: 1,
			AbConfig: `[
				{"layer_id": 1, "slot_ids": [0, 1, 2, 3, 4, 5, 6, 7, 8, 9], "white_list": ["test_device_1", "test_device_2"]},
				{"layer_id": 2, "slot_ids": [10, 11, 12, 13, 14, 15, 16, 17, 18, 19], "white_list": []},
				{"layer_id": 3, "slot_ids": [20, 21, 22, 23, 24, 25, 26, 27, 28, 29], "white_list": []}
			]`,
		},
		{
			AbId: 2,
			AbConfig: `[
				{"layer_id": 1, "slot_ids": [0, 1, 2, 3, 4], "white_list": ["test_device_3"]},
				{"layer_id": 2, "slot_ids": [5, 6, 7, 8, 9], "white_list": []}
			]`,
		},
	}

	// 处理配置 / Process configurations
	if err := model.ProcessABConfig(configs); err != nil {
		logrus.Error("Failed to process AB config: ", err)
	}
}
