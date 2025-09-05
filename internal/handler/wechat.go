package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"vista/config"
	"vista/internal/database"
	"vista/internal/repository"
	"vista/internal/service"
	"vista/pkg/response"
	"vista/pkg/utils"

	"github.com/gin-gonic/gin"
)

// WechatAuth 发起微信授权，重定向到微信授权页面
func WechatAuth(c *gin.Context) {
	// 生成 state 参数防止 CSRF 攻击
	state := utils.GenerateState()

	// 构建微信授权 URL
	authURL := buildWechatAuthURL(state)

	fmt.Printf("Generated auth URL: %s\n", authURL)

	// 重定向到微信授权页面
	c.Redirect(http.StatusFound, authURL)
}

// WechatCallback 处理微信授权回调
func WechatCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	fmt.Printf("WeChat callback received - Code: %s, State: %s\n", code, state)

	if code == "" {
		fmt.Println("Error: Missing authorization code")
		response.BadRequest(c, "authorization code is required")
		return
	}

	// 验证 state 参数（这里简化了，实际项目中应该存储在 session 或 cache 中）
	if state == "" {
		fmt.Println("Error: Missing state parameter")
		response.BadRequest(c, "invalid state parameter")
		return
	}

	// 创建服务实例
	userRepo := repository.NewWechatUserRepository(database.GetDB())
	wechatService := service.NewWechatOAuthService(userRepo)

	fmt.Println("Starting WeChat OAuth process...")

	// 使用 code 换取 access_token 和用户信息
	authResult, err := wechatService.ExchangeCodeForToken(code)
	if err != nil {
		fmt.Printf("Error exchanging code for token: %v\n", err)
		response.InternalServerError(c, fmt.Sprintf("微信登录失败: %v", err))
		return
	}

	fmt.Printf("Successfully got access token for OpenID: %s\n", authResult.OpenID)

	// 测试号支持获取完整用户信息
	userInfo, err := wechatService.GetUserInfo(authResult.AccessToken, authResult.OpenID)
	if err != nil {
		fmt.Printf("Error getting user info: %v\n", err)
		response.InternalServerError(c, fmt.Sprintf("获取用户信息失败: %v", err))
		return
	}

	fmt.Printf("Successfully got user info - OpenID: %s, Nickname: %s\n", userInfo.OpenID, userInfo.Nickname)

	// 保存用户授权信息
	if err := wechatService.SaveUserAuth(authResult, userInfo); err != nil {
		fmt.Printf("Error saving user auth: %v\n", err)
		response.InternalServerError(c, fmt.Sprintf("保存用户信息失败: %v", err))
		return
	}

	fmt.Println("Successfully saved user auth to database")

	// 检查是否是AJAX请求
	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" ||
		c.GetHeader("Accept") == "application/json" ||
		c.Query("format") == "json" {
		// 返回JSON响应
		response.Success(c, gin.H{
			"message": "微信登录成功",
			"user_info": gin.H{
				"openid":     authResult.OpenID,
				"nickname":   userInfo.Nickname,
				"headimgurl": userInfo.HeadImgURL,
				"sex":        userInfo.Sex,
				"language":   userInfo.Language,
				"country":    userInfo.Country,
				"province":   userInfo.Province,
				"city":       userInfo.City,
				"privilege":  userInfo.Privilege,
				"note":       "使用微信测试号获取完整用户信息",
			},
		})
		return
	}

	// 返回HTML页面，包含用户信息供前端JavaScript获取
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>微信登录成功</title>
    <style>
        body { font-family: Arial, sans-serif; padding: 20px; text-align: center; }
        .success { color: #52c41a; }
        .user-info { background: #f5f5f5; padding: 15px; margin: 20px 0; border-radius: 8px; }
        .avatar { width: 80px; height: 80px; border-radius: 50%%; margin: 10px 0; }
        button { background: #1890ff; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; margin: 5px; }
        .code { background: #f0f0f0; padding: 10px; margin: 10px 0; border-radius: 4px; font-family: monospace; }
    </style>
</head>
<body>
    <h2 class="success">✅ 微信登录成功</h2>
    
    <div class="user-info">
        <img class="avatar" src="%s" alt="头像" onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODAiIGhlaWdodD0iODAiIHZpZXdCb3g9IjAgMCA4MCA4MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iNDAiIGN5PSI0MCIgcj0iNDAiIGZpbGw9IiNmMGYwZjAiLz4KPHN2ZyB3aWR0aD0iODAiIGhlaWdodD0iODAiIHZpZXdCb3g9IjAgMCA4MCA4MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHRleHQgeD0iNDAiIHk9IjQ1IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LXNpemU9IjE0IiBmaWxsPSIjOTk5Ij7op6blg488L3RleHQ+Cjwvc3ZnPgo='"/>
        <h3>%s</h3>
        <p>OpenID: %s</p>
    </div>

    <div>
        <button onclick="getUserInfo()">获取用户信息</button>
        <button onclick="closeWindow()">关闭窗口</button>
        <button onclick="copyUserInfo()">复制用户信息</button>
    </div>

    <div class="code" id="userInfoJson" style="display:none;">
        <pre id="jsonContent"></pre>
    </div>

    <script>
        // 用户信息数据
        const userInfo = %s;
        
        // 将用户信息存储到 localStorage
        localStorage.setItem('wechat_user_info', JSON.stringify(userInfo));
        
        // 触发自定义事件，通知父页面
        if (window.parent !== window) {
            window.parent.postMessage({
                type: 'WECHAT_LOGIN_SUCCESS',
                data: userInfo
            }, '*');
        }
        
        // 如果是在微信浏览器中，尝试调用微信JS-SDK
        if (typeof WeixinJSBridge !== 'undefined') {
            WeixinJSBridge.call('closeWindow');
        }

        function getUserInfo() {
            const jsonElement = document.getElementById('userInfoJson');
            const contentElement = document.getElementById('jsonContent');
            
            if (jsonElement.style.display === 'none') {
                contentElement.textContent = JSON.stringify(userInfo, null, 2);
                jsonElement.style.display = 'block';
            } else {
                jsonElement.style.display = 'none';
            }
        }

        function copyUserInfo() {
            navigator.clipboard.writeText(JSON.stringify(userInfo, null, 2)).then(() => {
                alert('用户信息已复制到剪贴板');
            }).catch(() => {
                // 降级方案
                const textArea = document.createElement('textarea');
                textArea.value = JSON.stringify(userInfo, null, 2);
                document.body.appendChild(textArea);
                textArea.select();
                document.execCommand('copy');
                document.body.removeChild(textArea);
                alert('用户信息已复制到剪贴板');
            });
        }

        function closeWindow() {
            if (window.parent !== window) {
                // 在iframe中
                window.parent.postMessage({type: 'CLOSE_IFRAME'}, '*');
            } else {
                // 尝试关闭窗口
                window.close();
            }
        }

        // 3秒后自动尝试关闭（仅在微信浏览器中）
        setTimeout(() => {
            if (typeof WeixinJSBridge !== 'undefined') {
                WeixinJSBridge.call('closeWindow');
            }
        }, 3000);
    </script>
</body>
</html>`,
		userInfo.HeadImgURL,
		userInfo.Nickname,
		authResult.OpenID,
		func() string {
			userInfoJson := gin.H{
				"openid":     authResult.OpenID,
				"nickname":   userInfo.Nickname,
				"headimgurl": userInfo.HeadImgURL,
				"sex":        userInfo.Sex,
				"language":   userInfo.Language,
				"country":    userInfo.Country,
				"province":   userInfo.Province,
				"city":       userInfo.City,
				"privilege":  userInfo.Privilege,
			}
			jsonBytes, _ := json.Marshal(userInfoJson)
			return string(jsonBytes)
		}())

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 从URL参数获取openid
	openid := c.Param("openid")
	if openid == "" {
		// 也可以从query参数获取
		openid = c.Query("openid")
	}

	if openid == "" {
		response.BadRequest(c, "openid is required")
		return
	}

	// 从数据库获取用户信息
	repo := repository.NewWechatUserRepository(database.GetDB())
	user, err := repo.GetByOpenID(openid)
	if err != nil {
		response.Error(c, http.StatusNotFound, "user not found")
		return
	}

	// 返回用户信息
	response.SuccessWithMessage(c, "获取用户信息成功", gin.H{
		"user_info": gin.H{
			"openid":     user.OpenID,
			"nickname":   user.Nickname,
			"sex":        user.Sex,
			"language":   user.Language,
			"city":       user.City,
			"province":   user.Province,
			"country":    user.Country,
			"headimgurl": user.HeadImgURL,
			"unionid":    user.UnionID,
			"privilege":  user.Privilege,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// GetAllUsers 获取所有用户列表
func GetAllUsers(c *gin.Context) {
	repo := repository.NewWechatUserRepository(database.GetDB())
	users, err := repo.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get users")
		return
	}

	// 构建用户列表响应
	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"openid":     user.OpenID,
			"nickname":   user.Nickname,
			"sex":        user.Sex,
			"language":   user.Language,
			"city":       user.City,
			"province":   user.Province,
			"country":    user.Country,
			"headimgurl": user.HeadImgURL,
			"unionid":    user.UnionID,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	response.SuccessWithMessage(c, "获取用户列表成功", gin.H{
		"users": userList,
		"total": len(userList),
	})
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	// 检查数据库连接
	if err := database.HealthCheck(); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("database health check failed: %v", err))
		return
	}

	response.SuccessWithMessage(c, "service is healthy", gin.H{
		"status":   "ok",
		"service":  "vista-wechat-auth",
		"database": "connected",
	})
}

// buildWechatAuthURL 构建微信授权 URL
func buildWechatAuthURL(state string) string {
	cfg := config.Get()

	// 使用URL编码的回调地址
	redirectURI := url.QueryEscape(cfg.Wechat.RedirectURI)

	// 构建授权URL - 测试号支持完整功能，使用 snsapi_userinfo
	authURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect",
		cfg.Wechat.AppID,
		redirectURI,
		state,
	)

	fmt.Printf("Building WeChat Test Account auth URL:\n")
	fmt.Printf("  AppID: %s\n", cfg.Wechat.AppID)
	fmt.Printf("  RedirectURI: %s\n", cfg.Wechat.RedirectURI)
	fmt.Printf("  EncodedURI: %s\n", redirectURI)
	fmt.Printf("  State: %s\n", state)
	fmt.Printf("  Final URL: %s\n", authURL)

	return authURL
}
