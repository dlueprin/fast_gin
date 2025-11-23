package user_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/model"
	"fast_gin/utils/captcha"
	"fast_gin/utils/jwts"
	"fast_gin/utils/pwd"
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Username    string `json:"username" binding:"required" label:"用户名"`
	Password    string `json:"password" binding:"required" label:"密码"`
	CaptchaID   string `json:"captchaID"`
	CaptchaCode string `json:"captchaCode"` //用户输入的验证码
}

func (UserApi) LoginView(c *gin.Context) {
	//var cr LoginRequest
	////参数绑定就是获取请求中各种格式的参数到自己定义的结构体里的
	//if err := c.ShouldBindJSON(&cr); err != nil {
	//	msg := validate.ValidateError(err)
	//	logrus.Errorf("参数校验失败：%s", msg)
	//	res.FailWithError(err, c)
	//	return
	//}
	cr := middleware.GetBind[LoginRequest](c)
	if global.Config.Site.Login.Captcha {
		if cr.CaptchaID == "" || cr.CaptchaCode == "" {
			res.FailWithMsg("请输入图片验证码", c)
			return
		}
		if !captcha.CaptchaStore.Verify(cr.CaptchaID, cr.CaptchaCode, true) { //true代表验证通过了就销毁
			res.FailWithMsg("图片验证码校验失败", c)
			return
		}
	}
	var user model.UserModel
	err := global.DB.Take(&user, "username = ?", cr.Username).Error //妈的，这里要加个.Error,不然就是全找不到
	if err != nil {
		logrus.Errorf("用户输入的用户名不存在:%s", cr.Username)
		res.FailWithMsg("用户名或密码错误", c) //即使知道是用户名不存在也不能只提示用户名不存在，防止暴力破解用户名
		return
	}
	if !pwd.CompareHashAndPassword(user.Password, cr.Password) {
		logrus.Errorf("用户输入的密码错误：%s", cr.Password)
		res.FailWithMsg("用户名或密码错误", c)
	}
	token, err1 := jwts.SetToken(jwts.Claims{
		UserID: user.ID,
		RoleID: user.RoleID,
	})
	if err1 != nil {
		logrus.Errorf("token颁发失败%s", err)
		res.FailWithMsg("登录失败", c)
	}
	res.OkWithData(token, c)
	return
}
