package pkg

type SyncJson struct {
	Status int     `json:"status"` //点歌状态 0 停止点歌 1 允许点歌
	Length int     `json:"length"` //歌曲总数
	Data   *[]Song `json:"data"`   //歌曲列表
}

type ResultMessage struct {
	OK  int    `json:"ok"`
	Msg string `json:"msg"`
}

// GSync 与前端同步的消息列表
var GSync SyncJson

const (
	SongDisable = 0 // 停止点歌
	SongEnable  = 1 // 允许点歌
)
