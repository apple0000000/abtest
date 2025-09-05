package model

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// ABTest : 单条ab实验 / Single AB test
type ABTest struct {
	AbID      int64    `json:"-"`
	LayerID   int64    `json:"layer_id"`   // 分层ID / Layer ID
	SlotIDs   []int64  `json:"slot_ids"`   // 槽位ID列表 / Slot IDs
	WhiteList []string `json:"white_list"` // 白名单设备ID / Whitelist device IDs
}

var (
	// ABTests : [公共内存]所有的ab配置，[]*ABTest / [Shared memory] All AB configurations, []*ABTest
	ABTests atomic.Value
)

const (
	// MaxSlot : Slot最大值 / Maximum slot value
	MaxSlot = 1000
)

// AbTestConfig 定义AB测试配置结构 / AB test configuration structure
type AbTestConfig struct {
	AbId     int64  `json:"ab_id"`     // 实验ID / Experiment ID
	AbConfig string `json:"ab_config"` // 实验配置 / Experiment configuration
}

// ProcessABConfig : 处理配置更新 / Process configuration updates
func ProcessABConfig(ab_configs []*AbTestConfig) error {
	startTm := time.Now().UnixNano()

	// 临时变量（实验id组,白名单) / Temporary variable (experiment ID group, whitelist)
	var ABTestsTmp = make([]*ABTest, 0)

	for _, ab_config_presto := range ab_configs {
		if ab_config_presto.AbId <= 0 {
			continue
		}

		ab_config_list := make([]*ABTest, 0)
		err := json.Unmarshal([]byte(ab_config_presto.AbConfig), &ab_config_list)
		if err != nil {
			logrus.Errorf("[ProcessABConfig] json err, %v", err)
			continue
		}

		for _, ab_config := range ab_config_list {
			// 存ABTeststmp / Store to ABTestsTmp
			ABTestsTmp = append(ABTestsTmp, &ABTest{
				AbID:      ab_config_presto.AbId,
				LayerID:   ab_config.LayerID,
				SlotIDs:   ab_config.SlotIDs,
				WhiteList: ab_config.WhiteList,
			})
		}
	}

	// Tmp转正式 / Convert temporary to formal
	ABTests.Store(ABTestsTmp)

	endTm := time.Now().UnixNano()
	consumeTm := endTm - startTm

	logrus.Infof("[ProcessABConfig] time :%v us, config:%v", consumeTm/1000, util.ToJson(ABTestsTmp))
	return nil
}
