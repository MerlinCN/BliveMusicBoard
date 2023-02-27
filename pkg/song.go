package pkg

import (
	"log"
	"time"
)

type Song struct {
	Number         int       `json:"number"`
	Name           string    `json:"name"`
	User           string    `json:"user"`
	UID            int       `json:"uid"`
	RequestedAt    time.Time `json:"requestedAt"`
	RequestedAtStr string    `json:"requestedAtStr"`
}

var SongNumber = 0

var Songs []Song

var HistorySongs []Song

func AddSong(name string, user string, uid int) {
	nowTime := time.Now()
	SongNumber += 1
	// 创建一个新的 Song 实例
	song := Song{
		Number:         SongNumber,
		Name:           name,
		User:           user,
		UID:            uid,
		RequestedAtStr: nowTime.Format("2006-01-02 15:04:05"),
		RequestedAt:    nowTime,
	}

	// 将新的歌曲添加到 songs 列表中
	Songs = append(Songs, song)
	HistorySongs = append(HistorySongs, song)
	// 使用 log 包输出添加歌曲的信息
	log.Printf("添加歌曲: %+v \n", song)
	SyncSongs()
}

func DeleteSong(number int) ResultMessage {
	var index = -1
	for i, song := range Songs {
		if song.Number == number {
			index = i
			break
		}
	}

	var res = ResultMessage{OK: 1, Msg: "删除成功"}
	if index == -1 {
		log.Printf("不存在歌曲%d", number)
		res = ResultMessage{OK: 0, Msg: "歌曲不存在"}
	} else {
		log.Printf("删除歌曲: %+v\n", Songs[index])
		Songs = append(Songs[:index], Songs[index+1:]...)
	}
	SyncSongs()
	return res
}
