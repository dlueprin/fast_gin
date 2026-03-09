package doc_ser

import (
	"bytes"
	"encoding/json"
	"fast_gin/global"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// SearchResult chroma搜索返回内容的结构体
type SearchResult struct {
	Text     string  `json:"text"`
	DocID    int     `json:"doc_id"`
	Distance float64 `json:"distance"` //相似度距离
}

// EmbeddingResponse 用户输入的向量化返回
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// GetEmbeddings 获取用户输入的向量化接口
func (s *DocService) GetEmbeddings(chunks []string) ([][]float32, error) {
	//1、获取配置
	apiKey := global.Config.AI.Embedding.ApiKey
	baseURL := global.Config.AI.Embedding.BaseUrl
	model := global.Config.AI.Embedding.Model

	//2、新建resty客户端并发送请求
	client := resty.New()
	var result EmbeddingResponse
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(map[string]any{
			"model": model,
			"input": chunks,
		}).
		SetResult(&result).
		Post(baseURL)

	//3、处理错误
	if err != nil {
		return nil, fmt.Errorf("向量化接口请求失败，%v", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("网络请求错误：%s,错误码：%d", resp.String(), resp.StatusCode())
	}
	if result.Error.Message != "" {
		return nil, fmt.Errorf("api业务错误:%s", result.Error.Message)
	}

	//4、提取向量数据
	//由于是按顺序返回的，所以只需要按顺序赋值即可
	vectors := make([][]float32, len(chunks))
	for _, item := range result.Data {
		if item.Index < len(vectors) {
			vectors[item.Index] = item.Embedding
		}
	}

	return vectors, nil
}

// GetCollectionID 获取集合ID用于存库
func GetCollectionID() (string, error) {
	name := global.Config.Chroma.Collection

	listURL := getChromaAdminURL() + "/collections"

	resp, err := http.Get(listURL)
	if err != nil {
		return "", fmt.Errorf("连接chroma失败：%s", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var collections []map[string]any //包含多个map的切片，通常用来接受json数组相应
	json.Unmarshal(body, &collections)

	for _, col := range collections {
		if col["name"] == name {
			return col["id"].(string), nil
		}
	}

	logrus.Infof("[chroma]正在创建新集合:%s", name)
	createPayload, _ := json.Marshal(map[string]any{
		"name": name,
		//使用余弦相似度查询规则，适用条件是长度无关，而语义有关，适合文本
		"metadata": map[string]any{"hnsw:space": "cosine"},
	})

	cResp, err := http.Post(listURL, "application/json", bytes.NewBuffer(createPayload))
	if err != nil {
		return "", err
	}
	defer cResp.Body.Close() // 使用defer确保最后关闭

	var result map[string]any
	cBody, _ := io.ReadAll(cResp.Body)
	//if err != nil {
	//	return "", fmt.Errorf("读取响应体失败: %v", err)
	//}
	json.Unmarshal(cBody, &result)

	if id, ok := result["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("创建集合失败,%s", string(cBody))
}

// SyncToChroma chroma同步接口
func (s *DocService) SyncToChroma(colID string, docID int, userID int, chunks []string, vectors [][]float32) error {
	//1、构造插入url
	upsertURL := fmt.Sprintf("%s/collections/%s/upsert", getChromaAdminURL(), colID)

	//2、准备payload
	//Chroma 要求 IDs、Embeddings、Documents、Metadatas 四个数组长度必须一致
	ids := make([]string, len(chunks))
	metadatas := make([]map[string]any, len(chunks))
	for i := range chunks {
		ids[i] = fmt.Sprintf("doc_%d_chunk_%d", docID, i)
		metadatas[i] = map[string]any{
			"doc_id":  docID,
			"user_id": userID,
		}
	}
	payload := map[string]any{
		"ids":        ids,
		"embeddings": vectors,
		"documents":  chunks,
		"metadatas":  metadatas,
	}
	body, _ := json.Marshal(payload)
	//3、发送post请求
	resp, err := http.Post(upsertURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("请求Chroma失败,%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chroma响应错误：%s", string(respBody))
	}
	return nil
}

// Search chroma的用户语义搜索
func (s *DocService) Search(userID int, query string, nResult int) ([]SearchResult, error) {
	//1、向量化用户输入
	queryVectors, err := s.GetEmbeddings([]string{query})
	if err != nil || len(queryVectors) == 0 {
		return nil, fmt.Errorf("语义解析失败：%v", err)
	}

	//2、获取集合id,用于查询
	colID, err := GetCollectionID()
	if err != nil {
		return nil, err
	}

	//3、chroma的query请求，进行数据库查询
	queryURL := fmt.Sprintf("%s/collections/%s/query", getChromaAdminURL(), colID)
	payload := map[string]any{
		"query_embeddings": queryVectors,
		"n_results":        nResult,
		"where":            map[string]any{"user_id": userID},
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(queryURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//4、解析响应
	var queryRes struct {
		Documents [][]string         `json:"documents"`
		Metadatas [][]map[string]any `json:"metadatas"`
		Distances [][]float64        `json:"distances"`
	}
	json.NewDecoder(resp.Body).Decode(&queryRes)

	//5、组装返回格式
	var results []SearchResult
	if len(queryRes.Documents) > 0 {
		for i := range queryRes.Documents[0] {
			docID := int(queryRes.Metadatas[0][i]["doc_id"].(float64))
			results = append(results, SearchResult{
				Text:     queryRes.Documents[0][i],
				DocID:    docID,
				Distance: queryRes.Distances[0][i],
			})
		}
	}

	return results, nil
}

// GetChatContext 获取整合好的完整返回文本,现在新增一个分开的结果，用于结果的资料卡片单独展示
func (s *DocService) GetChatContext(userID int, query string) (results []SearchResult, contextText string, err error) {
	//1、获取相关度最高的结果
	results, err = s.Search(userID, query, 5)
	if err != nil {
		return nil, "", err
	}

	//2、拼接所有结果
	for i, result := range results {
		contextText += fmt.Sprintf("\n---资料片段%d---\n%s\n", i+1, result.Text)
	}
	return results, contextText, nil
}
func getChromaAdminURL() string {
	return fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database", global.Config.Chroma.Host)
}
