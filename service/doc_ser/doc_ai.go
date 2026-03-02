package doc_ser

import (
	"context"
	"fast_gin/global"
	"fast_gin/model"
	"fast_gin/utils/textclean"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

//下面是用resty实现的ai接口调用，ai说不推荐使用，但是可以作为参考
//type AiChatResponse struct {
//	Choices []struct {
//		Message struct {
//			Content string `json:"content"`
//		} `json:"message"`
//	} `json:"choices"`
//}
//
//func (s *DocService) AsyncAnalyze(doc *model.DocumentModel) {
//	//1、读数据库文档，为ai阅读做准备
//	contentBytes, err := os.ReadFile(doc.Path)
//	if err != nil {
//		logrus.Errorf("[ai]无法读取文件：%v", err)
//		return
//	}
//	content := string(contentBytes)
//
//	//2、使用resty调用ai
//	client := resty.New()
//	var result AiChatResponse
//	resp, err := client.R().
//		SetHeader("Authorization", "Bearer "+global.Config.AI.ApiKey).
//		SetBody(map[string]any{
//			"model": global.Config.AI.Model,
//			"messages": []map[string]string{
//				{"role": "system", "content": "你是一个专业的知识库助手，请对用户提供的文档内容进行100字以内的摘要。"},
//				{"role": "user", "content": content},
//			},
//		}).
//		SetResult(&result).
//		Post(global.Config.AI.BaseUrl + "/chat/completions")
//	//if err != nil || len(result.Choices) == 0 {
//	//	logrus.Errorf("[ai]接口调用异常%v", err)
//	//	return
//	//}
//	// --- 开始详细排查日志 ---
//	if err != nil {
//		logrus.Errorf("[AI] 网络请求失败: %v", err)
//		return
//	}
//
//	// 如果没有 choices，打印出到底回了什么
//	if len(result.Choices) == 0 {
//		logrus.Errorf("[AI] 接口未返回结果。状态码: %d, 响应原文: %s", resp.StatusCode(), resp.String())
//		return
//	}
//
//	//3.接受结果
//	summary := result.Choices[0].Message.Content
//
//	//4.异步更新数据库
//	err = global.DB.Model(doc).Updates(map[string]any{
//		"summary":     summary,
//		"content":     content,
//		"is_analyzed": true,
//	}).Error
//
//	if err != nil {
//		logrus.Errorf("[ai]数据库更新失败%v", err)
//	} else {
//		logrus.Infof("[ai]文章《%s》ai摘要生成完成", doc.Title)
//	}
//}

// 下面是用火山引擎的sdk实现调用
func (s *DocService) AsyncAnalyze(doc *model.DocumentModel, ext string) {
	//1、调用服务层解析方法读本地的文字内容
	var content string
	var err error
	switch ext {
	case ".md", ".txt":
		content, err = ReadPlainText(doc.Path)
	case ".pdf":
		content, err = ReadPdfText(doc.Path)
	case ".docx":
		content, err = ReadDocxText(doc.Path)
	default:
		err = fmt.Errorf("暂不支持当前格式：%s", ext)
	}
	if err != nil {
		logrus.Errorf("[ai]解析文件失败:%s", err)
		return
	}

	//内部转换带来的无效格式（多余的换行符或类似\x00）的处理
	content = textclean.CleanInvalidUTF8(content)
	//多余换行符和空格处理
	content = textclean.CleanSpacing(content)
	//过长文档处理，限制5000字，减少token消耗
	content = smartTruncate(content, 5000)

	////2、初始化ai链接配置
	//config := openai.DefaultConfig(global.Config.AI.LLM.ApiKey)
	//config.BaseURL = global.Config.AI.LLM.BaseUrl
	//client := openai.NewClientWithConfig(config)
	//
	////3、发送请求
	//logrus.Infof("[ai]正在调用ai接口进行摘要生成，使用模型：%s", global.Config.AI.LLM.Model)
	//resp, err := client.CreateChatCompletion(
	//	context.Background(),
	//	openai.ChatCompletionRequest{
	//		Model: global.Config.AI.LLM.Model,
	//		Messages: []openai.ChatCompletionMessage{
	//			{
	//				Role:    openai.ChatMessageRoleSystem,
	//				Content: "你是专业文档分析助手，用中文输出不超过100字的精准摘要，无需开场白和结束语，直接提供核心要点。",
	//			},
	//			{
	//				Role:    openai.ChatMessageRoleUser,
	//				Content: content,
	//			},
	//		},
	//	},
	//)

	//2、初始化ai链接配置,发送请求
	resp, err := GetChatResult("你是专业文档分析助手，用中文输出不超过100字的精准摘要，无需开场白和结束语，直接提供核心要点。", content)
	logrus.Infof("[ai]正在调用ai接口进行摘要生成，使用模型：%s", global.Config.AI.LLM.Model)

	//4、异常处理
	if err != nil {
		logrus.Errorf("[ai]SDK调用失败：%s", err)
		return
	}
	if len(resp.Choices) <= 0 {
		logrus.Errorf("[ai]接口未返回结果")
		return
	}

	//5、提取并更新数据库
	summary := resp.Choices[0].Message.Content
	err = global.DB.Model(doc).Updates(map[string]any{
		"content":     content,
		"summary":     summary,
		"is_analyzed": true,
	}).Error

	if err != nil {
		logrus.Errorf("[ai]数据库更新失败%s", err)
		return
	} else {
		logrus.Infof("[ai]文章《%s》ai摘要生成完成，摘要已存入数据库", doc.Title)
	}

	// --- 第二阶段：RAG 向量入库 ---
	logrus.Infof("[RAG]开始为文档《%s》建立向量索引", doc.Title)

	//1、创建或获取集合id
	colId, err := GetCollectionID()
	if err != nil {
		logrus.Errorf("[RAG]获取chroma集合失败:%v", err)
		return
	}

	//2、正文分块
	chunks := MakeChunks(content, 500)

	//3、获取向量
	vectors, err := s.GetEmbeddings(chunks)
	if err != nil {
		logrus.Errorf("[RAG]向量化失败：%v", err)
		return
	}

	//4、向量化入库
	err = s.SyncToChroma(colId, doc.ID, doc.UserID, chunks, vectors)
	if err != nil {
		logrus.Errorf("[RAG]Chroma入库失败：%v", err)
		return
	}

	logrus.Infof("[分析完成]文档《%s》已经处理完毕，摘要已经存入mysql，%d个向量块已经存入chroma", doc.Title, len(chunks))
}

func smartTruncate(text string, maxLen int) string {
	//之前没用rune切片，是用字节为单位的，如果切的时候恰好出现把汉字切断的情况，就会导致乱码，然后数据库就敏感了报错
	//现在rune是按字符为单位，就不会出现切断汉字的情况，涉及汉字切片都建议用rune处理
	runes := []rune(text)

	if len(runes) <= maxLen {
		return text
	}
	headLen := int(float64(maxLen) * 0.7)
	tailLen := maxLen - headLen

	head := string(runes[:headLen])
	tail := string(runes[len(runes)-tailLen:])

	return fmt.Sprintf("%s\n\n...[此处省略中间部分]...\n\n%s", head, tail)
}

// 分块函数
func MakeChunks(text string, chunkSize int) (chunks []string) {
	runes := []rune(text)
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return
}

func GetChatResult(systemPrompt string, userPrompt string) (response openai.ChatCompletionResponse, err error) {
	//初始化ai链接配置
	config := openai.DefaultConfig(global.Config.AI.LLM.ApiKey)
	config.BaseURL = global.Config.AI.LLM.BaseUrl
	client := openai.NewClientWithConfig(config)

	//发送请求
	response, err = client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: global.Config.AI.LLM.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)
	return
}
