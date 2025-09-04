# vista 
微信授权服务端

概述：

要在微信中打开网页并获取微信用户的 **ID（openid/unionid）、昵称、头像等信息**，需要借助微信的 **OAuth2 授权登录**，大致流程是这样的：

---

## 🔑 实现步骤

1. **注册微信公众平台账号**

   * 你需要一个 **微信公众号（服务号）**，订阅号是不支持网页授权获取用户信息的。
   * 在 [微信公众平台](https://mp.weixin.qq.com/) 申请服务号，并完成认证。

2. **配置网页授权域名**

   * 在微信公众平台后台，找到 **设置 → 公众号设置 → 功能设置 → 网页授权域名**。
   * 填写你网页所在的域名（必须备案过的）。

3. **前端跳转获取 code**
   在网页里，用户访问时，跳转到微信的 OAuth2 授权链接：

   ```text
   https://open.weixin.qq.com/connect/oauth2/authorize?
   appid=你的公众号AppID
   &redirect_uri=你的后端接收code的URL（URLEncode）
   &response_type=code
   &scope=snsapi_userinfo
   &state=123#wechat_redirect
   ```

   * `scope=snsapi_base` 只能获取用户的 openid（静默授权，不弹窗）。
   * `scope=snsapi_userinfo` 可以获取昵称、头像、性别等（需要用户点击授权）。

4. **后端用 code 换取 access\_token + openid**
   用户同意后，微信会带着 `code` 回调到 `redirect_uri`，你在服务端用这个 `code` 请求：

   ```text
   https://api.weixin.qq.com/sns/oauth2/access_token?
   appid=你的AppID
   &secret=你的AppSecret
   &code=上一步拿到的code
   &grant_type=authorization_code
   ```

   返回结果会包含：

   ```json
   {
     "access_token": "ACCESS_TOKEN",
     "expires_in": 7200,
     "refresh_token": "REFRESH_TOKEN",
     "openid": "用户的OPENID",
     "scope": "snsapi_userinfo"
   }
   ```

5. **拉取用户信息（昵称、头像等）**
   如果 scope 是 `snsapi_userinfo`，还可以用 access\_token 获取更多信息：

   ```text
   https://api.weixin.qq.com/sns/userinfo?
   access_token=ACCESS_TOKEN
   &openid=OPENID
   &lang=zh_CN
   ```

   返回结果：

   ```json
   {
     "openid":"OPENID",
     "nickname":"用户昵称",
     "sex":1,
     "province":"Guangdong",
     "city":"Guangzhou",
     "country":"China",
     "headimgurl":"用户头像URL",
     "privilege":[]
   }
   ```

---

## ⚠️ 注意事项

* 必须是 **服务号** 才能获取用户信息，订阅号不行。
* 域名必须备案并在公众号后台配置。
* 如果只是想识别用户身份（不需要昵称头像），可以用 `snsapi_base`，用户无感知。

---
