package flags

import (
	"fast_gin/global"
	"fast_gin/model"
	"fast_gin/utils/pwd"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type User struct {
}

// Create 创建用户
// 记得小写是不能导出的，还有在函数名前面的括号里面的，代表是这个结构体所属对象的方法，也就是要user点的方法
func (User) Create() {
	var user model.UserModel
	fmt.Println("请选择角色： 1 管理员 2 普通用户")
	_, err := fmt.Scanln(&user.RoleID)
	if err != nil {
		fmt.Printf("输入错误：%s", err)
		return
	}
	if user.RoleID != 1 && user.RoleID != 2 {
		fmt.Println("用户角色输入错误")
		return
	}
	fmt.Println("请输入用户名：")
	fmt.Scanln(&user.Username)
	//检测重名
	var u model.UserModel
	err = global.DB.Take(&u, "username = ?", user.Username).Error
	if err == nil {
		fmt.Println("用户已存在")
		return
	}
	fmt.Println("请输入密码：")                                      //使用终端密码不可见，防止日志泄露密码
	password, err := terminal.ReadPassword(int(os.Stdin.Fd())) //这里stdin是终端输入，fd是类似输入的东西的所存的文件句柄
	if err != nil {
		fmt.Printf("密码读取时出错：%s", err)
		return
	}
	fmt.Println("请再次输入密码：")
	rePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("密码读取时出错：%s", err)
		return
	}
	if string(password) != string(rePassword) {
		fmt.Println("密码输入不一致")
		return
	}
	//加密成哈希值
	hashPwd := pwd.GenerateFromPassword(string(password)) //密码本身是[]byte类型，所以要转换成string类型
	err = global.DB.Create(&model.UserModel{
		Username: user.Username,
		//Nickname: user.Username,
		Password: hashPwd,
		RoleID:   user.RoleID,
	}).Error
	if err != nil {
		logrus.Errorf("用户创建失败 %s", err)
		return
	}
	logrus.Infof("用户创建成功")
}
func (User) List() {
	var userList []model.UserModel
	global.DB.Order("created_at desc").Limit(10).Find(&userList)
	for _, user := range userList {
		fmt.Printf("用户id: %d 用户名: %s 昵称: %s 用户角色: %d 创建时间: %s\n",
			user.ID,
			user.Username,
			user.Nickname,
			user.RoleID,
			user.CreatedAt.Format("2006-01-02 15:04:05")) //格式化时间，显示记得写一下这个
	}
}
