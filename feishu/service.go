package feishu

import (
	"context"
	"errors"
	"fmt"
	"github.com/airunny/wiki-go-tools/locker"
	"time"

	"github.com/go-kratos/kratos/v2/log" // nolint
	goCache "github.com/liyanbing/go-cache"
	goCacheError "github.com/liyanbing/go-cache/errors"
	goCacheTool "github.com/liyanbing/go-cache/tools"
)

var ErrTimeout = errors.New("timeout")

const (
	appAccessTokenURL  = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
	userAccessTokenURL = "https://open.feishu.cn/open-apis/authen/v1/oidc/access_token"
	userInfoURL        = "https://open.feishu.cn/open-apis/authen/v1/user_info"
)

const (
	tokenLockKey            = "FSAppTokenLocker"
	tokenKeyFormat          = "FSAppToken:%v"
	tokenExpireEarlySeconds = int64(5 * 60)
	lockExpire              = time.Second * 10
)

// AppAccessToken 获取自建应用的app_access_token（https://open.feishu.cn/document/server-docs/authentication-management/access-token/app_access_token_internal）
func (c *FSClient) AppAccessToken(ctx context.Context) (string, error) {
	key := fmt.Sprintf(tokenKeyFormat, c.Config.AppID)
	// 这里token存在redis中，设置有过期时间，过期之后删除；需要时重新获取即可，不需要主动刷新
	token, err := goCache.FetchWithString(ctx, c.cache, key, func() (value interface{}, expiration time.Duration, err error) {
		// 分布式锁
		release, err := c.locker.TryLock(ctx, tokenLockKey, lockExpire)
		if err != nil {
			return "", 0, err
		}
		defer func() {
			err1 := release()
			if err1 != nil {
				log.Context(ctx).Errorf("Release Err:%v", err1)
			}
		}()

		newToken, err := c.getAccessTokenFromFeiShu(ctx)
		if err != nil {
			return "", 0, err
		}

		// 提前五分钟过期，接口文档是2小时过期时间
		expire := newToken.Expire - tokenExpireEarlySeconds
		// 防止意外，这里多做一层判断保护好自己
		if expire <= 0 {
			expire = newToken.Expire
		}
		return newToken.AppAccessToken, time.Duration(expire) * time.Second, nil
	})
	// 如果发生不可预期的错误
	if err != nil && err != locker.ErrAlreadyLocked {
		return "", err
	}

	// 到这里说明是并发情况下没有拿到锁
	if err != nil {
		timer := time.After(5 * time.Second)
		for {
			select {
			case <-timer: // 超时
				return "", ErrTimeout
			case <-time.After(time.Millisecond * 5):
				var cacheData interface{}
				cacheData, err = c.cache.Get(ctx, key)
				if err != nil && err != goCacheError.ErrEmptyCache {
					return "", err
				}
				return goCacheTool.ToString(cacheData)
			}
		}
	}

	return token, nil
}

func (c *FSClient) getAccessTokenFromFeiShu(ctx context.Context) (*AppAccessToken, error) {
	req := AppAccessTokenRequest{
		AppID:     c.Config.AppID,
		AppSecret: c.Config.AppSecret,
	}

	var out AppAccessTokenResponse
	err := c.httpPost(ctx, &doRequest{
		domain: appAccessTokenURL,
		req:    req,
		out:    &out,
	})
	if err != nil {
		return nil, err
	}
	return out.AppAccessToken, nil
}

// UserAccessToken 获取用户的 access_token（
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/oidc-access_token/create?appId=cli_a59e637f58bdd00b）
func (c *FSClient) UserAccessToken(ctx context.Context, in *UserAccessTokenRequest) (*UserTokenData, error) {
	if in.Code == "" {
		return nil, fmt.Errorf("empty code")
	}

	if in.GrantType == "" {
		in.GrantType = "authorization_code"
	}

	token, err := c.AppAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	var out UserAccessTokenResponse
	err = c.httpPost(ctx, &doRequest{
		domain:        userAccessTokenURL,
		req:           in,
		out:           &out,
		authorization: token,
	})
	if err != nil {
		return nil, err
	}

	return out.Data, nil
}

// GetUserInfo 获取登录用户信息（https://open.feishu.cn/document/server-docs/authentication-management/login-state-management/get）
func (c *FSClient) GetUserInfo(ctx context.Context, in *GetUserInfoRequest) (*UserInfoData, error) {
	if in.UserAccessToken == "" {
		return nil, fmt.Errorf("empty user_access_token")
	}

	var out GetUserInfoResponse
	err := c.httpGet(ctx, &doRequest{
		domain:        userInfoURL,
		out:           &out,
		authorization: in.UserAccessToken,
	})
	if err != nil {
		return nil, err
	}
	return out.Data, nil
}
