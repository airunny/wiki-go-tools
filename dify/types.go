package dify

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CheckResponse interface {
	Check() error
}

type CommonResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (s CommonResponse) Check() error {
	if s.Code != "" {
		return fmt.Errorf("%v[%v]", s.Message, s.Code)
	}
	return nil
}

// 索引方式
const (
	IndexingTechniqueHighQuality string = "high_quality" //  高质量：使用 embedding 模型进行嵌入，构建为向量数据库索引
	IndexingTechniqueEconomy     string = "economy"      // 经济：使用 Keyword Table Index 的倒排索引进行构建
)

const (
	ProcessRuleModeAutomatic string = "automatic" // 自动
	ProcessRuleModeCustom    string = "custom"    // 自定义
)

const (
	PreProcessingRuleIDRemoveExtraSpaces string = "remove_extra_spaces" // 替换连续空格、换行符、制表符
	PreProcessingRuleIDRemoveUrlsEmails  string = "remove_urls_emails"  // 删除 URL、电子邮件地址
)

type PreProcessingRule struct {
	ID      string `json:"id"`      // 预处理规则的唯一标识符
	Enabled bool   `json:"enabled"` // 是否选中该规则，不传入文档 ID 时代表默认值
}

type Segmentation struct {
	Separator string `json:"separator"`  // 自定义分段标识符，目前仅允许设置一个分隔符。默认为 \n
	MaxTokens int64  `json:"max_tokens"` // 最大长度 (token) 默认为 1000
}

type Rule struct {
	PreProcessingRules []*PreProcessingRule `json:"pre_processing_rules"` // 预处理规则
	Segmentation       *Segmentation        `json:"segmentation"`         // 分段规则
}

type ProcessRule struct {
	Mode  string `json:"mode"`  // 清洗、分段模式 ，automatic 自动 / custom 自定义
	Rules *Rule  `json:"rules"` // 自定义规则（自动模式下，该字段为空）
}

// ============= 创建空知识库 ===============

type AddDatasetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Permission  string `json:"permission"`
}

type AddDatasetResponse struct {
	CommonResponse
	Id                     string      `json:"id"`
	Name                   string      `json:"name"`
	Description            interface{} `json:"description"`
	Provider               string      `json:"provider"`
	Permission             string      `json:"permission"`
	DataSourceType         interface{} `json:"data_source_type"`
	IndexingTechnique      interface{} `json:"indexing_technique"`
	AppCount               int         `json:"app_count"`
	DocumentCount          int         `json:"document_count"`
	WordCount              int         `json:"word_count"`
	CreatedBy              string      `json:"created_by"`
	CreatedAt              int         `json:"created_at"`
	UpdatedBy              string      `json:"updated_by"`
	UpdatedAt              int         `json:"updated_at"`
	EmbeddingModel         interface{} `json:"embedding_model"`
	EmbeddingModelProvider interface{} `json:"embedding_model_provider"`
	EmbeddingAvailable     interface{} `json:"embedding_available"`
}

// ============= 获取知识库列表 ===============

type FindDatasetRequest struct {
	Page  string `json:"page"`  // 页码
	Limit string `json:"limit"` // 返回条数，默认 20，范围 1-100
}

type Dataset struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Permission        string `json:"permission"`
	DataSourceType    string `json:"data_source_type"`
	IndexingTechnique string `json:"indexing_technique"`
	AppCount          int    `json:"app_count"`
	DocumentCount     int    `json:"document_count"`
	WordCount         int    `json:"word_count"`
	CreatedBy         string `json:"created_by"`
	CreatedAt         int64  `json:"created_at"`
	UpdatedBy         string `json:"updated_by"`
	UpdatedAt         int64  `json:"updated_at"`
}

type FindDatasetResponse struct {
	CommonResponse
	Data    []*Dataset `json:"data"`
	HasMore bool       `json:"has_more"`
	Limit   int        `json:"limit"`
	Total   int        `json:"total"`
	Page    int        `json:"page"`
}

// ============= 删除知识库 ===============

type DeleteDatasetRequest struct {
	DatasetId string `json:"dataset_id"` // 知识库ID
}

type DeleteDatasetResponse struct {
	CommonResponse
}

// ============= 新增文档 ===============

type AddDocumentRequest struct {
	DatasetId         string       `json:"-"`                  // 知识库ID
	Name              string       `json:"name"`               // 文档名称
	Text              string       `json:"text"`               // 文档内容
	IndexingTechnique string       `json:"indexing_technique"` // 索引方式，默认:high_quality
	ProcessRule       *ProcessRule `json:"process_rule"`       // 处理规则
}

type AddDocumentResponse struct {
	CommonResponse
	Document struct {
		Id             string `json:"id"`
		Position       int    `json:"position"`
		DataSourceType string `json:"data_source_type"`
		DataSourceInfo struct {
			UploadFileId string `json:"upload_file_id"`
		} `json:"data_source_info"`
		DatasetProcessRuleId string      `json:"dataset_process_rule_id"`
		Name                 string      `json:"name"`
		CreatedFrom          string      `json:"created_from"`
		CreatedBy            string      `json:"created_by"`
		CreatedAt            int         `json:"created_at"`
		Tokens               int         `json:"tokens"`
		IndexingStatus       string      `json:"indexing_status"`
		Error                interface{} `json:"error"`
		Enabled              bool        `json:"enabled"`
		DisabledAt           interface{} `json:"disabled_at"`
		DisabledBy           interface{} `json:"disabled_by"`
		Archived             bool        `json:"archived"`
		DisplayStatus        string      `json:"display_status"`
		WordCount            int         `json:"word_count"`
		HitCount             int         `json:"hit_count"`
		DocForm              string      `json:"doc_form"`
	} `json:"document"`
	Batch string `json:"batch"`
}

// ============= 更新文档 ===============

type UpdateDocumentRequest struct {
	DatasetId         string       `json:"dataset_id"`         // 知识库ID
	DocumentId        string       `json:"document_id"`        // 文档ID
	Name              string       `json:"name"`               // 文档名称
	Text              string       `json:"text"`               // 文档内容
	IndexingTechnique string       `json:"indexing_technique"` // 索引方式，默认:high_quality
	ProcessRule       *ProcessRule `json:"process_rule"`       // 处理规则
}

type UpdateDocumentResponse struct {
	CommonResponse
	Document struct {
		Id             string `json:"id"`
		Position       int    `json:"position"`
		DataSourceType string `json:"data_source_type"`
		DataSourceInfo struct {
			UploadFileId string `json:"upload_file_id"`
		} `json:"data_source_info"`
		DatasetProcessRuleId string      `json:"dataset_process_rule_id"`
		Name                 string      `json:"name"`
		CreatedFrom          string      `json:"created_from"`
		CreatedBy            string      `json:"created_by"`
		CreatedAt            int         `json:"created_at"`
		Tokens               int         `json:"tokens"`
		IndexingStatus       string      `json:"indexing_status"`
		Error                interface{} `json:"error"`
		Enabled              bool        `json:"enabled"`
		DisabledAt           interface{} `json:"disabled_at"`
		DisabledBy           interface{} `json:"disabled_by"`
		Archived             bool        `json:"archived"`
		DisplayStatus        string      `json:"display_status"`
		WordCount            int         `json:"word_count"`
		HitCount             int         `json:"hit_count"`
		DocForm              string      `json:"doc_form"`
	} `json:"document"`
	Batch string `json:"batch"`
}

// ============= 删除文档 ===============

type DeleteDocumentRequest struct {
	DatasetId  string `json:"dataset_id"`  // 知识库ID
	DocumentId string `json:"document_id"` // 文档ID
}

type DeleteDocumentResponse struct {
	CommonResponse
	Result string `json:"result"`
}

// ============= 文本生成型应用 API ===============

type Inputs map[string]interface{}

type File struct {
	Type           string `json:"type"`                     // (string) 支持类型：图片 image（目前仅支持图片格式）
	TransferMethod string `json:"transfer_method"`          // 传递方式: remote_url: 图片地址。local_file: 上传文件
	URL            string `json:"url,omitempty"`            // 图片地址。（仅当传递方式为 remote_url 时）。
	UploadFileId   string `json:"upload_file_id,omitempty"` // 上传文件 ID。（仅当传递方式为 local_file 时）。
}

type CompletionMessagesRequest struct {
	Inputs       Inputs  `json:"inputs,omitempty"`
	ResponseMode string  `json:"response_mode,omitempty"`
	User         string  `json:"user,omitempty"`
	Files        []*File `json:"files,omitempty"`
}

type ChatCompletionResponse struct {
	MessageId string `json:"message_id"`
	Mode      string `json:"mode"`
	Answer    string `json:"answer"`
}

type ChunkChatCompletionResponse struct {
}

type CompletionMessagesResponse struct {
	CommonResponse
	ChatCompletionResponse
}

type UploadFileRequest struct {
	FilePath string `json:"file_path"`
	User     string `json:"user"`
}

type UploadFileResponse struct {
	CommonResponse
	Id        string `json:"id"`
	Name      string `json:"name"`       // 文件名
	Size      int64  `json:"size"`       // 文件大小（byte）
	Extension string `json:"extension"`  // 文件后缀
	MimeType  string `json:"mime_type"`  // 文件 mime-type
	CreatedBy string `json:"created_by"` // 上传人 ID
	CreatedAt int64  `json:"created_at"` // 上传时间
}

// ============= 聊天型应用 API ===============

type ChatMessagesRequest struct {
	Query            string  `json:"query"`
	Inputs           Inputs  `json:"inputs"`
	ResponseMode     string  `json:"response_mode,omitempty"`
	User             string  `json:"user,omitempty"`
	ConversationId   string  `json:"conversation_id,omitempty"`
	Files            []*File `json:"files,omitempty"`
	AutoGenerateName bool    `json:"auto_generate_name"`
}

type ChatMessagesResponse struct {
	CommonResponse
	ChatCompletionResponse
	Body []byte `json:"body"`
}

type ChunkMessage struct {
	Event          string `json:"event"`
	ConversationId string `json:"conversation_id"`
	MessageId      string `json:"message_id"`
	CreatedAt      int    `json:"created_at"`
	TaskId         string `json:"task_id"`
	Id             string `json:"id"`
	Answer         string `json:"answer"`
}

func (s *ChatMessagesResponse) ChunkData(v any) error {
	body := string(s.Body)
	splits := strings.Split(body, "\n")
	values := make([]string, 0, len(splits))
	for _, split := range splits {
		if strings.TrimSpace(split) == "" {
			continue
		}
		values = append(values, strings.TrimSpace(strings.TrimPrefix(split, "data: ")))
	}
	newBody := fmt.Sprintf("[%s]", strings.Join(values, ","))
	return json.Unmarshal([]byte(newBody), v)
}

func (s *ChatMessagesResponse) ChunkAnswer() (string, error) {
	var (
		values []*ChunkMessage
		result []string
	)
	err := s.ChunkData(&values)
	if err != nil {
		return "", err
	}

	for _, value := range values {
		if value.Event == "message_end" {
			break
		}
		result = append(result, value.Answer)
	}
	return strings.Join(result, ""), nil
}
