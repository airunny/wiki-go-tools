package dify

import (
	"bytes"
	"context"
	"fmt"
	"os"
)

// AddDataset 新增知识库
func (s *Client) AddDataset(ctx context.Context, in *AddDatasetRequest) (*AddDatasetResponse, error) {
	var out AddDatasetResponse
	_, err := s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetBody(in).
		SetResult(&out).
		Post(fmt.Sprintf("%v/datasets", s.url))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// FindDataset 获取知识库列表
func (s *Client) FindDataset(ctx context.Context, in *FindDatasetRequest) (*FindDatasetResponse, error) {
	var out FindDatasetResponse
	_, err := s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetResult(&out).
		Get(fmt.Sprintf("%v/datasets?page=%s&limit=%s", s.url, in.Page, in.Limit))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDataset 删除知识库
func (s *Client) DeleteDataset(ctx context.Context, in *DeleteDatasetRequest) (*DeleteDatasetResponse, error) {
	var out DeleteDatasetResponse
	_, err := s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetResult(&out).
		Delete(fmt.Sprintf("%v/datasets/%s", s.url, in.DatasetId))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// AddDocument 新增文档
func (s *Client) AddDocument(ctx context.Context, in *AddDocumentRequest) (*AddDocumentResponse, error) {
	var out AddDocumentResponse
	_, err := s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetBody(in).
		SetResult(&out).
		Post(fmt.Sprintf("%v/datasets/%s/document/create_by_text", s.url, in.DatasetId))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDocument 更新文档
func (s *Client) UpdateDocument(ctx context.Context, in *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	var out UpdateDocumentResponse
	_, err := s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetBody(in).
		SetResult(&out).
		Post(fmt.Sprintf("%v/datasets/%s/documents/%s/update_by_text", s.url, in.DatasetId, in.DocumentId))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDocument 删除文档
func (s *Client) DeleteDocument(ctx context.Context, in *DeleteDocumentRequest) (*DeleteDocumentResponse, error) {
	var out DeleteDocumentResponse
	_, err := s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetBody(in).
		SetResult(&out).
		Delete(fmt.Sprintf("%v/datasets/%s/documents/%s", s.url, in.DatasetId, in.DocumentId))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CompletionMessages 文本生成型应用 API
func (s *Client) CompletionMessages(ctx context.Context, in *CompletionMessagesRequest) (*CompletionMessagesResponse, error) {
	var out CompletionMessagesResponse
	_, err := s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetBody(in).
		SetResult(&out).
		Post(fmt.Sprintf("%s/completion-messages", s.url))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ChatMessages 聊天型应用 API
func (s *Client) ChatMessages(ctx context.Context, in *ChatMessagesRequest) (*ChatMessagesResponse, error) {
	var out ChatMessagesResponse
	res, err := s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetBody(in).
		SetResult(&out).
		Post(fmt.Sprintf("%s/chat-messages", s.url))
	if err != nil {
		return nil, err
	}
	out.Body = res.Body()
	return &out, nil
}

// Upload 上传图片
func (s *Client) Upload(ctx context.Context, in *UploadFileRequest) (*UploadFileResponse, error) {
	profileImgBytes, err := os.ReadFile(in.FilePath)
	if err != nil {
		return nil, err
	}

	var out UploadFileResponse
	_, err = s.httpClient.R().
		SetHeader("Authorization", in.APIKey).
		SetContext(ctx).
		SetFileReader("file", in.FilePath, bytes.NewReader(profileImgBytes)).
		SetFormData(map[string]string{
			"user": in.User,
		}).
		SetResult(&out).
		Post(fmt.Sprintf("%s/files/upload", s.url))
	if err != nil {
		return nil, err
	}
	return &out, nil
}
