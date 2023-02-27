package main

import (
	"BliveMusicBoard/pkg"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"strconv"
)

const (
	PORT = "4220"
)

func main() {
	pkg.InitSettings()
	// 初始化 SongList
	pkg.Songs = []pkg.Song{}

	pkg.GSync = pkg.SyncJson{
		Status: pkg.SongEnable,
		Data:   &pkg.Songs,
		Length: len(pkg.Songs),
	}
	// 创建一个路由器
	router := http.NewServeMux()

	// 注册 /songs 路由的处理函数
	router.HandleFunc("/songs", syncSongsHandler)

	// 注册 /songs/add 路由的处理函数
	router.HandleFunc("/songs/add", addSongHandler)
	// 注册 /songs/del 路由的处理函数
	router.HandleFunc("/songs/del", deleteSongHandler)

	router.HandleFunc("/songs/clear", clearSongsHandler)
	router.HandleFunc("/songs/setting/sync", syncSettingHandler)
	router.HandleFunc("/songs/setting/edit", editSettingHandler)

	router.HandleFunc("/songs/switch", switchHandler)

	// 注册 /admin 路由的处理函数
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	// 注册静态文件路由
	router.Handle("/", http.FileServer(http.Dir("./public")))
	router.HandleFunc("/songs/ws", pkg.HandleWebSocket)
	// 使用 CORS 中间件
	handler := handlers.CORS(
	//handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
	//handlers.AllowedMethods([]string{"GET", "POST", "DELETE"}),
	//handlers.AllowedOrigins([]string{"*"}),
	)(router)
	// 使用 log 包输出服务器启动信息
	log.Printf("服务已启动")
	log.Printf("OBS添加浏览器 http://localhost:%s \n", PORT)
	log.Printf("后台管理界面 http://localhost:%s/admin \n", PORT)
	log.Println("关注芋艿谢谢喵~ https://live.bilibili.com/326763")
	//初始化弹幕监听
	go pkg.SetDanmaku(pkg.GSetting.RoomId.Val.(string))

	// 使用 log 包启动服务器，如果启动失败，程序会直接退出
	go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), handler))
}

//处理 /songs 路由请求，返回所有的歌曲列表以及消息
func syncSongsHandler(w http.ResponseWriter, r *http.Request) {

	res := pkg.GSync
	jsonData, err := json.Marshal(res)
	if err != nil {
		// 输出错误信息并返回 500 错误码
		http.Error(w, "JSON marshal error", http.StatusInternalServerError)
		return
	}

	// 将 JSON 数据写入响应体
	w.Write(jsonData)

	//清空消息列表
	//pkg.GSync.Msg = make([]pkg.ResultMessage, 0)

}

//处理 /songs/add 路由请求，将一首新歌曲添加到列表中
func addSongHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取歌曲名和点歌人信息
	name := r.FormValue("name")
	user := r.FormValue("user")
	uidStr := r.FormValue("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		http.Error(w, "uid 错误", http.StatusBadRequest)
		return
	}
	pkg.AddSong(name, user, uid)
}

func deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Number int `json:"number"`
	}
	var data req
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid number", http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(pkg.DeleteSong(data.Number))
	// 将 JSON 数据写入响应体
	w.Write(jsonData)
}

func clearSongsHandler(w http.ResponseWriter, r *http.Request) {
	pkg.Songs = []pkg.Song{}
	log.Printf("All songs cleared")
}

func syncSettingHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, _ := json.Marshal(pkg.GSetting)
	// 将 JSON 数据写入响应体
	w.Write(jsonData)
}

func editSettingHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&pkg.GSetting)
	var res = pkg.ResultMessage{OK: 1, Msg: "保存成功"}
	if err != nil {
		res.OK = 0
		res.Msg = "保存失败，请检查设置"
	}
	jsonData, _ := json.Marshal(res)
	// 将 JSON 数据写入响应体
	w.Write(jsonData)
	pkg.SaveSetting()
}

func switchHandler(w http.ResponseWriter, r *http.Request) {

	msg := "开启点歌"
	if pkg.GSync.Status == pkg.SongEnable {
		pkg.GSync.Status = pkg.SongDisable
		msg = "关闭点歌"
	} else {
		pkg.GSync.Status = pkg.SongEnable
	}

	type ResultMessageWithStatus struct {
		OK     int    `json:"ok"`
		Msg    string `json:"msg"`
		Status int    `json:"status"`
	}

	res := ResultMessageWithStatus{OK: 1, Msg: msg, Status: pkg.GSync.Status}
	pkg.SendMessage(1, msg)
	jsonData, _ := json.Marshal(res)
	// 将 JSON 数据写入响应体
	w.Write(jsonData)
	log.Println(msg)
}
