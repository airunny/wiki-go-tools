package feishu

import "fmt"

type (
	CommonResponse struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}

	AppAccessTokenRequest struct {
		AppID     string `json:"app_id"`
		AppSecret string `json:"app_secret"`
	}

	AppAccessToken struct {
		AppAccessToken string `json:"app_access_token"`
		Expire         int64  `json:"expire"`
	}

	AppAccessTokenResponse struct {
		CommonResponse
		*AppAccessToken
	}

	UserAccessTokenRequest struct {
		GrantType string `json:"grant_type"`
		Code      string `json:"code"`
	}

	UserTokenData struct {
		AccessToken      string `json:"access_token"`       // 字段access_token即user_access_token，用于获取用户资源和访问某些open api
		RefreshToken     string `json:"refresh_token"`      // 刷新user_access_token时使用的 refresh_token
		TokenType        string `json:"token_type"`         // token 类型，固定值
		ExpiresIn        int64  `json:"expires_in"`         // user_access_token有效期，单位: 秒，有效时间两个小时左右，需要以返回结果为准
		RefreshExpiresIn int64  `json:"refresh_expires_in"` // refresh_token有效期，单位: 秒，一般是30天左右，需要以返回结果为准
		Scope            string `json:"scope"`              // 用户授予app的权限全集
	}
	UserAccessTokenResponse struct {
		CommonResponse
		Data *UserTokenData `json:"data"`
	}

	GetUserInfoRequest struct {
		UserAccessToken string `json:"user_access_token"`
	}
	UserInfoData struct {
		Name            string `json:"name"`             // 用户姓名
		EnName          string `json:"en_name"`          // 用户英文名称
		AvatarURL       string `json:"avatar_url"`       // 用户头像
		AvatarThumb     string `json:"avatar_thumb"`     // 用户头像 72x72
		AvatarMiddle    string `json:"avatar_middle"`    // 用户头像 240x240
		AvatarBig       string `json:"avatar_big"`       // 用户头像 640x640
		OpenID          string `json:"open_id"`          // 用户在应用内的唯一标识
		UnionID         string `json:"union_id"`         // 用户对ISV的唯一标识，对于同一个ISV，用户在其名下所有应用的union_id相同
		Email           string `json:"email"`            // 用户邮箱（权限）
		EnterpriseEmail string `json:"enterprise_email"` // 企业邮箱，请先确保已在管理后台启用飞书邮箱服务（权限）
		UserID          string `json:"user_id"`          // 用户 user_id（权限）
		Mobile          string `json:"mobile"`           // 用户手机号（权限）
		TenantKey       string `json:"tenant_key"`       // 当前企业标识
		EmployeeNo      string `json:"employee_no"`      // 用户工号
	}
	GetUserInfoResponse struct {
		CommonResponse
		Data *UserInfoData `json:"data"`
	}
)

func (s CommonResponse) Check() error {
	if s.Code != 0 {
		return fmt.Errorf("%v[%v]", s.Message, s.Code)
	}
	return nil
}
