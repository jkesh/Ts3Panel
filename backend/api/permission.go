package api

import (
	"Ts3Panel/core"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/ts3"
)

// --- 结构体定义 ---

type PermReq struct {
	PermName  string `json:"perm_name" binding:"required"`
	PermValue int    `json:"perm_value"`
}

type ServerGroupPerm struct {
	PermID  int    `ts3:"permid" json:"permid"`
	Name    string `ts3:"permsid" json:"name"`
	Value   int    `ts3:"permvalue" json:"value"`
	Negated int    `ts3:"permnegated" json:"negated"`
	Skip    int    `ts3:"permskip" json:"skip"`
}

// 定义一个简单的权限定义结构，用于解析 permissionlist 命令
type PermDefinition struct {
	PermID   int    `ts3:"permid"`
	PermName string `ts3:"permname"`
	PermDesc string `ts3:"permdesc"`
}

// --- 缓存机制 ---

var (
	permCache     = make(map[int]string)
	permCacheLock sync.RWMutex
	permCached    = false
)

// ensurePermCache 确保权限缓存已加载
func ensurePermCache(ctx context.Context) {
	permCacheLock.RLock()
	if permCached {
		permCacheLock.RUnlock()
		return
	}
	permCacheLock.RUnlock()

	permCacheLock.Lock()
	defer permCacheLock.Unlock()

	// 双重检查
	if permCached {
		return
	}

	log.Println("[Permission] Loading permission list cache from TS3...")

	// 执行 permissionlist 命令
	// 注意：有些服务器可能需要 permissionlist -new
	resp, err := core.Client.Exec(ctx, "permissionlist")
	if err != nil {
		log.Printf("[Permission] Failed to load permission list: %v", err)
		return
	}

	var defs []PermDefinition
	if err := ts3.NewDecoder().Decode(resp, &defs); err != nil {
		log.Printf("[Permission] Failed to decode permission list: %v", err)
		return
	}

	for _, def := range defs {
		permCache[def.PermID] = def.PermName
	}

	permCached = true
	log.Printf("[Permission] Cache loaded. Total permissions: %d", len(permCache))
}

// getPermNameByID 从缓存中获取权限名
func getPermNameByID(id int) string {
	permCacheLock.RLock()
	defer permCacheLock.RUnlock()
	return permCache[id]
}

// --- 接口实现 ---

// ListServerGroupPerms 获取服务器组的当前权限列表
func ListServerGroupPerms(c *gin.Context) {
	sgid, _ := strconv.Atoi(c.Param("sgid"))

	core.Mutex.Lock()
	// 在锁住 core.Mutex 之前或之后都可以，这里为了简单直接执行
	// 注意：ensurePermCache 内部也会调用 Exec，Exec 内部会用 core.Mutex 吗？
	// 查看 core/ts3_manager.go，Exec 是 core.Client.Exec，它是线程安全的（Client 内部有锁）。
	// 但为了防止逻辑冲突，我们最好在获取组权限之前保证缓存加载。
	ensurePermCache(c.Request.Context())

	// 发送命令: servergrouppermlist sgid=xx -names
	cmd := fmt.Sprintf("servergrouppermlist sgid=%d -names", sgid)
	resp, err := core.Client.Exec(c.Request.Context(), cmd)

	core.Mutex.Unlock() // 尽早释放核心锁

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 调试日志：看看 TS3 到底返回了什么
	// log.Printf("[Debug] Raw Perm Resp: %s", resp)

	var perms []ServerGroupPerm
	if err := ts3.NewDecoder().Decode(resp, &perms); err != nil {
		c.JSON(200, gin.H{"data": []ServerGroupPerm{}})
		return
	}

	// [关键修复] 遍历列表，如果 Name 为空，则从缓存补全
	for i := range perms {
		if perms[i].Name == "" {
			cachedName := getPermNameByID(perms[i].PermID)
			if cachedName != "" {
				perms[i].Name = cachedName
			} else {
				// 如果缓存里也没有，说明这个 ID 可能很特殊，或者缓存加载失败
				// 保持为空，前端会显示 ID
			}
		}
	}

	c.JSON(200, gin.H{"data": perms})
}

// AddChannelPerm 修改频道权限
func AddChannelPerm(c *gin.Context) {
	cid, _ := strconv.Atoi(c.Param("cid"))
	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	if err := core.Client.ChannelAddPerm(c.Request.Context(), cid, req.PermName, req.PermValue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Updated"})
}

// AddServerGroupPerm 修改服务器组权限
func AddServerGroupPerm(c *gin.Context) {
	sgid, _ := strconv.Atoi(c.Param("sgid"))
	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	if err := core.Client.ServerGroupAddPerm(c.Request.Context(), sgid, req.PermName, req.PermValue, false, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Updated"})
}

// AddClientDbPerm 修改客户端(数据库ID)权限
func AddClientDbPerm(c *gin.Context) {
	cldbid, _ := strconv.Atoi(c.Param("cldbid"))
	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	if err := core.Client.ClientAddPerm(c.Request.Context(), cldbid, req.PermName, req.PermValue, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Updated"})
}
