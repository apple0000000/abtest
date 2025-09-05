package model

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"strconv"

	"abtest/util"
)

// ABResponse 定义AB测试响应结构 / AB test response structure
type ABResponse struct {
	Data *ABResponseData `json:"data"`
}

// ABResponseData 定义AB测试响应数据 / AB test response data
type ABResponseData struct {
	AbId string `json:"ab_id"` // AB测试ID / AB test ID
}

// GetAbId 根据设备ID获取AB测试分组 / Get AB test group based on device ID
func GetAbId(device_id string) *ABResponse {
	res := &ABResponse{
		Data: &ABResponseData{},
	}

	if device_id == "" {
		res.Data.AbId = "0-0"
		return res
	}

	// abtest整体序列的全局字符串 / Global string for the entire AB test sequence
	var strABIDAll bytes.Buffer

	// MapTestIDisSet : [实验id 互斥] 判断某个实验id是否被命中过的map,若命中过,则该实验id置1
	// MapTestIDisSet: [Experiment ID mutual exclusion] Map to determine if an experiment ID has been hit
	var MapTestIDisSet = make(map[int64]byte, 20)

	var TargetslotID int64

	// 临时变量（实验id组,白名单) / Temporary variable (experiment ID group, whitelist)
	var ABTestsTmp []*ABTest
	ABTestsTmpRaw := ABTests.Load()
	if ABTestsTmpRaw != nil {
		ABTestsTmp = ABTestsTmpRaw.([]*ABTest)
	}

	// 1) 先把白名单全部通过,白名单出现过的层,不再判断 / First, pass all whitelists, layers that have appeared in the whitelist are no longer judged
	for _, abtest := range ABTestsTmp {
		if util.SliceContainsString(abtest.WhiteList, device_id) {
			MapTestIDisSet[abtest.AbID] = 1

			if strABIDAll.Len() != 0 {
				strABIDAll.WriteByte(',')
			}
			strABIDAll.WriteString(util.GetAssertString(abtest.AbID))
			strABIDAll.WriteByte('-')
			strABIDAll.WriteString(util.GetAssertString(abtest.LayerID))
		}
	}

	// 正常循环判断 / Normal loop judgment
	for _, abtest := range ABTestsTmp {
		// 2）如果这个实验id已经命中过了,就无需再判断了 / If this experiment ID has been hit, no need to judge again
		if MapTestIDisSet[abtest.AbID] == 1 {
			continue
		}

		// 5) 验证第1000组 / Verify the 1000th group
		slotlen := len(abtest.SlotIDs)
		if slotlen == 0 {
			continue
		}

		if abtest.SlotIDs[slotlen-1] == MaxSlot {
			MapTestIDisSet[abtest.AbID] = 1

			if strABIDAll.Len() != 0 {
				strABIDAll.WriteByte(',')
			}
			strABIDAll.WriteString(util.GetAssertString(abtest.AbID))
			strABIDAll.WriteByte('-')
			strABIDAll.WriteString(util.GetAssertString(abtest.LayerID))
			continue
		}

		// 6）正常计算 / Normal calculation
		buff := make([]byte, 0, 40)
		buff = append(buff, []byte(device_id)...)
		buff = append(buff, []byte(strconv.FormatUint(uint64(abtest.AbID), 10))...)
		md5Binary := md5.Sum(buff)
		TargetslotID = int64(binary.BigEndian.Uint64(md5Binary[:8]) % MaxSlot)

		if util.SliceContainsInt64(abtest.SlotIDs, TargetslotID) {
			MapTestIDisSet[abtest.AbID] = 1
			if strABIDAll.Len() != 0 {
				strABIDAll.WriteByte(',')
			}
			strABIDAll.WriteString(util.GetAssertString(abtest.AbID))
			strABIDAll.WriteByte('-')
			strABIDAll.WriteString(util.GetAssertString(abtest.LayerID))
		}
	}

	// 均未命中,0-0 / None hit, 0-0
	if strABIDAll.Len() == 0 {
		strABIDAll.WriteString("0-0")
	}

	res.Data.AbId = strABIDAll.String()

	return res
}
