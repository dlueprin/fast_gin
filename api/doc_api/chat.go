package doc_api

import (
	"encoding/json"
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DocChatRequest struct {
	Query string `json:"query" binding:"required"`
}

//	func (d DocApi) DocChatView(c *gin.Context) {
//		//1、绑定关键词
//		cr := middleware.GetBind[DocChatRequest](c)
//		//2、获取用户id
//		userID := middleware.GetAuth(c).UserID
//		//3、检索上下文（R）
//		contextText, err := docService.GetChatContext(userID, cr.Query)
//		if err != nil {
//			res.FailWithMsg("检索知识库失败", c)
//			return
//		}
//		//4、构造系统提示词
//		systemPrompt := fmt.Sprintf("你是一个专业的知识库助手。请根据提供的【参考资料】回答用户问题。" +
//			"要求:\n1.如果资料中没有提到，请诚实回答不知道，不要胡乱猜测。\n" +
//			"2.回答要简洁专业。\n" +
//			"【参考资料】如下:\n" + contextText)
//		//5、生成（G）
//		answer, err := doc_ser.GetChatResult(systemPrompt, cr.Query)
//		if err != nil {
//			res.FailWithMsg("ai思考失败"+err.Error(), c)
//		}
//		//6、返回结果
//		res.OkWithData(gin.H{
//			"answer":  answer.Choices[0].Message.Content,
//			"context": contextText,
//			"usage":   answer.Usage,
//		}, c)
//	}

func (d DocApi) DocChatView(c *gin.Context) {
	var fullReply string
	//1、绑定关键词
	cr := middleware.GetBind[DocChatRequest](c)
	//2、获取用户id
	userID := middleware.GetAuth(c).UserID
	//3、检索上下文（R）
	results, contextText, err := docService.GetChatContext(userID, cr.Query)
	if err != nil {
		c.SSEvent("error", "检索知识库失败")
		return
	}
	//4、构造系统提示词
	systemPrompt := fmt.Sprintf("你是一个专业的知识库助手。请根据提供的【参考资料】回答用户问题。" +
		"要求:\n1.如果资料中没有提到，请诚实回答不知道，不要胡乱猜测。\n" +
		"2.回答要简洁专业。\n" +
		"【排版要求】:\n1. 必须使用 Markdown 格式。\n2. 段落与段落之间，必须使用两个换行符(\\n\\n)进行分隔！这一点非常重要！\n3. 重要的要点请使用列表 (1. 2. 3. 或者 -) 输出。\n" +
		"【参考资料】:\n" + contextText)
	//5、设置流式的响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	//6、生成（G）并立即返回
	//先发参考资料包
	contextJson, _ := json.Marshal(results)
	c.SSEvent("content", string(contextJson))
	c.Writer.Flush()
	//再发sse正文
	err = docService.GetChatStream(systemPrompt, cr.Query, func(content string) {
		//每拿到一个字就会存起来并调用当前函数发送
		fullReply += content
		c.SSEvent("message", content)
		c.Writer.Flush() //强制刷新缓冲区
	})
	if err != nil {
		c.SSEvent("error", "ai响应中断"+err.Error())
		return
	}
	//7、保存聊天记录
	chatRecord := model.ChatHistoryModel{
		UserID:  userID,
		Query:   cr.Query,
		Answer:  fullReply,
		Context: string(contextJson),
	}
	if err = global.DB.Create(&chatRecord).Error; err != nil {
		logrus.Errorf("保存聊天记录失败:%v", err)
	}
	//8、回传对话id，用于历史记录查询
	c.SSEvent("id", chatRecord.ID)
}
