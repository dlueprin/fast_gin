package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"reflect"

	// 在 import 语句中添加中文翻译包
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)

// 存储翻译器实例，后面翻译备用
var trans ut.Translator

// 初始化自动执行
func init() {
	//新建中文翻译器
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("zh")
	//注册验证引擎
	v, ok := binding.Validator.Engine().(*validator.Validate) //断言并转化为验证类型的引擎
	if ok {
		//注册默认的中文翻译规则
		_ = zh_translations.RegisterDefaultTranslations(v, trans) //将验证类型的引擎和中文翻译器关联，实现注册
	} else {
		logrus.Errorf("验证引擎注册失败")
	}
	//实现字段翻译
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
}

func ValidateError(err error) string {
	//判断是否是验证类型的错误，不是或者json格式不对就直接返回原错误，不翻译
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		logrus.Errorf("参数绑定的错误不是指定的验证错误类型")
		return err.Error()
	}
	var list []string
	//遍历每一个错误并翻译后添加到返回列表里
	for _, e := range errs {
		list = append(list, e.Translate(trans))
	}
	return strings.Join(list, ";")
}
