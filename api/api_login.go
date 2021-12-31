package api

import (
	"encoding/json"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools/cache"
	"github.com/xfyun/xns/tools/uid"
	"time"
)

const (
	keySession    = "session"
	codeNeedLogin = 1
	keyUser = "key_user"
	headerToken   = "X-HYDRANT-TOKEN"
)

const (
	adminUser = "admin"
	readOnly  = "readonly"
)

const (
	codeUnAuthorization = 10401
)

var (
	sessionCache   = cache.NewCache()
	SessionTimeout = 7 * 24 * time.Hour // will be  overwrite by session_timeout in config
	UserDao        dao.User
)

type session struct {
	username string
	mode     string
}

// 校验登陆权限
func validateUser(ctx *fastserver.Context) {
	sessionKey := ctx.GetRequestHeader(headerToken)
	sess, ok := sessionCache.Get(sessionKey)
	if !ok {
		ctx.AbortWithStatusJson(255, resp{
			Code:    codeUnAuthorization,
			Message: "need login",
		})
		return
	}
	// 更新session
	sessionCache.Set(sessionKey, sess, SessionTimeout)
	ctx.SetUserValue(keyUser,sess)
	ctx.Next()

}

type logresp struct {
	Token  string `json:"token"`
	Header string `json:"header"`
}

// 创建user
func createUser(ctx *fastserver.Context) {
	sessionKey := ctx.GetRequestHeader(headerToken)
	ses, ok := sessionCache.Get(sessionKey)
	req := &createUserReq{}
	err := json.Unmarshal(ctx.FastCtx.Request.Body(),req)
	if err != nil{
		ctx.AbortWithStatusJson(400,resp{Message: "bind json error:"+err.Error(),Code: codeUnAuthorization})
		return
	}
	if !ok {
		counter ,err := UserDao.GetCount()
		if err != nil{
			ctx.AbortWithStatusJson(500,&resp{
				Message: err.Error(),
				Code: codeUnAuthorization,
			})
			return
		}
		if counter > 0{
			ctx.AbortWithStatusJson(403,&resp{
				Message: "you have no right to create user",
				Code: codeUnAuthorization,
			})
			return
		}

		err = UserDao.Create(&models.User{
			Username: req.Username,
			Password: req.Password,
			Type:     adminUser,
		})
		if err != nil{
			ctx.AbortWithStatusJson(400,&resp{Message: "create user error:"+err.Error(),Code: codeUnAuthorization})
			return
		}
		ctx.AbortWithStatusJson(200,&resp{Message: "create user success",Code: codeUnAuthorization})
		return
	}
	sess := ses.(*session)
	if sess.mode != adminUser{
		ctx.AbortWithStatusJson(403,&resp{
			Message: "you have no right to create user",
			Code: codeUnAuthorization,
		})
		return
	}

	if !checkMode(req.Mode){
		ctx.AbortWithStatusJson(403,&resp{
			Message: "invalid mode:"+req.Mode,
			Code: codeUnAuthorization,
		})
		return
	}
	err = UserDao.Create(&models.User{
		Username: req.Username,
		Password: req.Password,
		Type:     req.Mode,
	})
	if err != nil{
		ctx.AbortWithStatusJson(400,&resp{Message: "create user error:"+err.Error()})
		return
	}
	ctx.AbortWithStatusJson(200,&resp{Message: "create user success"})
}

// 获取已经登陆的用户信息
func getUser(ctx *fastserver.Context){
	sess , _ := ctx.GetUserValue(keyUser).(*session)
	if sess == nil{
		ctx.AbortWithStatusJson(200,&resp{
			Message: "not login ",
			Code: codeUnAuthorization,
		})
		return
	}

	ctx.AbortWithStatusJson(200,&createUserResp{
		Username: sess.username,
		Mode:     sess.mode,
	})
}

// 校验是否有写的权限
func validateWriteAccessRight(ctx *fastserver.Context){
	if ctx.Method == GET{
		ctx.Next()
		return
	}

	sess ,ok := ctx.GetUserValue(keyUser).(*session)
	if !ok{
		ctx.AbortWithStatusJson(403,resp{Message: " you have no right to write "})
		return
	}
	if sess.mode != adminUser{
		ctx.AbortWithStatusJson(403,resp{Message: " you have no right to write, because your mode is "+sess.mode})
		return
	}
	ctx.Next()
}

func checkMode(mode string)bool{
	switch mode {
	case adminUser,readOnly:
		return true
	}
	return false
}
// a u1 je
func login(ctx *fastserver.Context) {
	token := ctx.GetRequestHeader(headerToken)
	_ , ok := sessionCache.Get(token)
	if ok {
		// 已经登陆
		ctx.AbortWithStatusJson(200, &logresp{
			Token:  token,
			Header: headerToken,
		})
		return
	}
	// todo  checklogin

	args := &loginArgs{}

	err := json.Unmarshal(ctx.FastCtx.Request.Body(), args)
	if err != nil {
		ctx.AbortWithStatusJson(400, &resp{
			Code:    400,
			Message: "parse body error. not json",
		})
		return
	}
	mode, ok := loginCheck(args.UserName, args.Password)
	if !ok{
		ctx.AbortWithStatusJson(400, &resp{
			Code:    10003,
			Message: "invalid username or password ",
		})
		return
	}
	token = uid.UUid()
	sessionCache.Set(token, &session{
		username: args.UserName,
		mode:mode,
	}, SessionTimeout)

	ctx.AbortWithStatusJson(200, &logresp{
		Token:  token,
		Header: headerToken,
	})
}

func loginCheck(username, password string) (mode string, ok bool) {
	user ,err := UserDao.Get(username)
	if err != nil{
		return "",false
	}
	if user.Password != password{
		return "",false
	}
	return user.Type,true
}

type loginArgs struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


type createUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Mode string `json:"mode"`
}


type createUserResp struct {
	Username string `json:"username"`
	Mode string `json:"mode"`
}
