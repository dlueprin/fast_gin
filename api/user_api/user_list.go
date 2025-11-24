package user_api

import (
	"fast_gin/middleware"
	"fast_gin/model"
	"fast_gin/service/common"
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (UserApi) UserListView(c *gin.Context) {
	//cr:criteria request,请求的查询条件的信息，而查询条件是通过参数绑定获取的
	var cr = middleware.GetBind[model.PageInfo](c)

	////切片创建的必须参数有2个，第一个是切片的数据类型，第二个是切片长度，第三个可选，是切片容量，妈的，go知识又忘了
	//var userList = make([]model.UserModel, 0)
	////key不为空时，模糊匹配来查找关键词相关的数据库：where username like ’%zhang%‘
	//query := global.DB.Where("")
	//if cr.Key != "" {
	//	//链式构建，把上一次查询的条件构建后又可以重新写查询条件，但是只是写条件，没有真正执行
	//	query.Where("username like ? ", "%"+cr.Key+"%") //妈的，这里是like，不是=
	//}
	//offset := (cr.Page - 1) * cr.Limit
	//global.DB.Where(query).Limit(cr.Limit).Offset(offset).Order(cr.Order).Find(&userList)
	//
	//var count int64
	//global.DB.Model(model.UserModel{}).Where(query).Count(&count) //注意统计查询次数是要指定模型的

	userList, count, err := common.QueryList[model.UserModel](model.UserModel{}, common.QueryOption{
		PageInfo: cr,
		Likes:    []string{"username", "nickname"},
		Debug:    true,
	})
	if err != nil {
		logrus.Errorf("query failed,%s", err)
		res.FailWithMsg("query failed", c)
		return
	}
	res.OkWithList(userList, count, c)
}
