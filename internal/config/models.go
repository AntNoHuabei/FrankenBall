package config

import "fmt"

type Provider string

const (
	Qwen   Provider = "qwen"
	Ollama Provider = "ollama"
)

type ChatModelDefine struct {
	Model           string   `json:"model"`
	Provider        Provider `json:"provider"`
	SupportThinking bool     `json:"support_thinking"` //是否支持思考
	IsMultimodal    bool     `json:"is_multimodal"`    //是否是多模态模型
}

//func GetDefaultModels(provider Provider) ([]ChatModelDefine, error) {
//
//}

func GetDefaultEndpoint(provider Provider) (string, error) {

	switch provider {

	case Qwen:
		return "https://dashscope.aliyuncs.com/compatible-mode/v1", nil

	case Ollama:
		return "https://localhost:11434", nil
	}

	return "", fmt.Errorf("unknown provider: %s", provider)

}
