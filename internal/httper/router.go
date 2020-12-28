package httper

import (
	"github.com/ChangbaServer/nsqproxy/config"
	"net/http"
)

// 启动HTTP
func (h *Httper) router() {
	//获取状态
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		Success(w, "ok")
	})
	//获取角色
	http.HandleFunc("/getRole", func(w http.ResponseWriter, r *http.Request) {
		Success(w, config.SystemConfig.Role)
	})

	//后台
	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(h.statikFS)))
	//后台接口
	// 消费者管理
	consumeConfig := ConsumeConfig{}
	http.HandleFunc("/admin/api/consumeConfig/create", consumeConfig.Create)
	http.HandleFunc("/admin/api/consumeConfig/update", consumeConfig.Update)
	http.HandleFunc("/admin/api/consumeConfig/delete", consumeConfig.Delete)
	http.HandleFunc("/admin/api/consumeConfig/page", consumeConfig.Page)
	http.HandleFunc("/admin/api/consumeConfig/get", consumeConfig.Get)
	http.HandleFunc("/admin/api/consumeConfig/workList", consumeConfig.WorkList)
	//Worker机管理
	workServer := WorkServer{}
	http.HandleFunc("/admin/api/workServer/create", workServer.Create)
	http.HandleFunc("/admin/api/workServer/update", workServer.Update)
	http.HandleFunc("/admin/api/workServer/delete", workServer.Delete)
	http.HandleFunc("/admin/api/workServer/page", workServer.Page)
	http.HandleFunc("/admin/api/workServer/all", workServer.All)
	http.HandleFunc("/admin/api/workServer/get", workServer.Get)
	//消费者和Worker机关联关系管理
	consumeServerMap := ConsumeServerMap{}
	http.HandleFunc("/admin/api/consumeServerMap/create", consumeServerMap.Create)
	http.HandleFunc("/admin/api/consumeServerMap/update", consumeServerMap.Update)
	http.HandleFunc("/admin/api/consumeServerMap/delete", consumeServerMap.Delete)
	http.HandleFunc("/admin/api/consumeServerMap/getWork", consumeServerMap.GetWork)
}
