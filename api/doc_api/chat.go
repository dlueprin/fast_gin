package doc_api

import (
	"fast_gin/middleware"
	"fast_gin/service/doc_ser"
	"fast_gin/utils/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DocChatRequest struct {
	Query string `json:"query" binding:"required"`
}

func (d DocApi) DocChatView(c *gin.Context) {
	//1、绑定关键词
	cr := middleware.GetBind[DocChatRequest](c)
	//2、获取用户id
	userID := middleware.GetAuth(c).UserID
	//3、检索上下文（R）
	contextText, err := docService.GetChatContext(userID, cr.Query)
	if err != nil {
		res.FailWithMsg("检索知识库失败", c)
		return
	}
	//4、构造系统提示词
	systemPrompt := fmt.Sprintf("你是一个专业的知识库助手。请根据提供的【参考资料】回答用户问题。" +
		"要求:\n1.如果资料中没有提到，请诚实回答不知道，不要胡乱猜测。\n" +
		"2.回答要简洁专业。\n" +
		"【参考资料】如下:\n" + contextText)
	//5、生成（G）
	answer, err := doc_ser.GetChatResult(systemPrompt, cr.Query)
	if err != nil {
		res.FailWithMsg("ai思考失败"+err.Error(), c)
	}
	//6、返回结果
	res.OkWithData(gin.H{
		"answer":  answer.Choices[0].Message.Content,
		"context": contextText,
		"usage":   answer.Usage,
	}, c)
}
