package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

type Setting struct {
	RoomId       SettingItem `json:"room_id"`
	SongCd       SettingItem `json:"song_cd"`
	FansMedalMin SettingItem `json:"fans_medal_min"`
}

type SettingItem struct {
	Label string      `json:"label"`
	Val   interface{} `json:"val"`
}

var GSetting Setting

var RoomIdCache string

func InitSettings() {
	// 读取设置文件
	data, err := ioutil.ReadFile("setting.json")
	if err != nil {
		log.Fatalf("无法读取设置文件: %v", err)
	}

	// 解析设置文件
	err = json.Unmarshal(data, &GSetting)
	if err != nil {
		log.Fatalf("解析设置文件失败: %v", err)
	}
	log.Println("正在加载配置 ......")
	// 打印设置文件
	v := reflect.ValueOf(GSetting)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fmt.Printf("\t %s: %v\n", fieldType.Name, field.Interface())
	}

	RoomIdCache = GSetting.RoomId.Val.(string)

}

func SaveSetting() {

	data, err := ioutil.ReadFile("setting.json")
	if err != nil {
		log.Fatalf("无法读取设置文件: %v", err)
	}
	// 保存设置到文件
	data, err = json.MarshalIndent(GSetting, "", "    ")
	if err != nil {
		fmt.Println("Error: failed to serialize setting")
		return
	}
	err = ioutil.WriteFile("setting.json", data, 0644)
	if err != nil {
		fmt.Println("Error: failed to save setting.json")
		return
	}

	if GSetting.RoomId.Val.(string) != RoomIdCache { //重置弹幕监听
		RoomIdCache = GSetting.RoomId.Val.(string)
		go SetDanmaku(GSetting.RoomId.Val.(string))
	}

}
