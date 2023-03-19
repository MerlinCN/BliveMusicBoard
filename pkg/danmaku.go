package pkg

import (
	"fmt"
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	"log"
	"regexp"
	"time"
)

var c *client.Client

func SetDanmaku(RoomId string) {
	c = client.NewClient(RoomId) //关注芋艿谢谢喵
	c.OnDanmaku(handleSongRequest)

	err := c.Start()
	if err != nil {
		log.Fatal(err)
	}
	msg := fmt.Sprintf("房间号：%s 弹幕监听加载成功", RoomId)
	log.Printf(msg)
	SendMessage(1, msg)
}

func handleSongRequest(danmaku *message.Danmaku) {
	if danmaku.Type == message.EmoticonDanmaku {
		return
	}
	re := regexp.MustCompile(`^点歌\s*(.*)$`)
	match := re.FindStringSubmatch(danmaku.Content)
	if len(match) < 1 {
		return
	}

	if GSync.Status == SongDisable {
		msg := fmt.Sprintf("点歌失败，已关闭点歌")
		log.Printf(msg)
		SendMessage(0, msg)
		return
	}

	var interval = checkCD(danmaku.Sender.Uid)
	if interval > 0 {
		msg := fmt.Sprintf("%s的点歌正在CD中 剩余%d秒\n", danmaku.Sender.Uname, interval)
		log.Printf(msg)
		SendMessage(0, msg)
		return
	}
	songName := match[1]
	history := checkHistory(songName)
	if history != nil {
		msg := fmt.Sprintf("今天已经点过%s咯！，在%s由%s点\n", history.Name, history.RequestedAtStr, history.User)
		log.Printf(msg)
		SendMessage(0, msg)
		return
	}
	AddSong(songName, danmaku.Sender.Uname, danmaku.Sender.Uid)
}

//检查点歌的时间是否小于设置 如果是就返回剩余时间 如果不是就返回0
func checkCD(uid int) int {
	for _, song := range Songs {
		if song.UID == uid {
			lastSongInterval := int(time.Since(song.RequestedAt).Seconds()) //历史点歌时间和现在的间隔（秒）
			if lastSongInterval < int(GSetting.SongCd.Val.(float64)) {
				return int(GSetting.SongCd.Val.(float64)) - lastSongInterval
			}
		}
	}
	return 0
}

func checkHistory(name string) *Song {
	for _, song := range HistorySongs {
		if song.Name == name {
			songInterval := int(time.Since(song.RequestedAt).Seconds())
			if songInterval <= 86400 {
				return &song
			}
		}
	}
	return nil
}
