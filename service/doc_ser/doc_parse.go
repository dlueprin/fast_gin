package doc_ser

import (
	"bytes"
	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"os"
	"regexp"
)

// 纯文本提取
func ReadPlainText(path string) (string, error) {
	content, err := os.ReadFile(path)
	return string(content), err
}

// pdf文本提取
func ReadPdfText(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		text, _ := p.GetPlainText(nil)
		buf.WriteString(text)
	}

	//正则表达式清洗多余换行符和空格
	re := regexp.MustCompile(`\s+`)
	res := re.ReplaceAllString(buf.String(), " ")

	return res, nil
}

// docx文本提取
func ReadDocxText(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	f, err := os.Open(path) //这里打开的是docx文件句柄
	if err != nil {
		return "", err
	}
	defer f.Close() //关闭文件句柄f

	//这里是打开docx压缩包的临时文件句柄
	d, err := docx.ReadDocxFromMemory(f, fi.Size())
	if err != nil {
		return "", err
	}
	defer d.Close() //关闭临时文件

	//获取包含文档内容的xml格式的内容
	rawXml := d.Editable().GetContent()
	//定义正则表达式，匹配<w:t>标签中的内容，这个标签里面就是正文
	re := regexp.MustCompile(`<w:t.*?>(.*?)</w:t>`)
	//查找所有匹配项
	matches := re.FindAllStringSubmatch(rawXml, -1)

	var res string
	for _, match := range matches {
		if len(match) > 1 {
			res += match[1] //这个match[1]是每次匹配到的内容，单个标签内的文本
		}
	}
	return res, nil
}
