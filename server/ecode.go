package server

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var webErrorCode = map[int][]string{
	0: {"200", "success"},

	1001: {"403", "未登录", "没有登录或登录超时"}, // 10XX 用户或登录相关
	1002: {"400", "用户不存在", "未找到用户"},
	1003: {"400", "登录失败", "密码不对或帐号不存在"},
	1004: {"400", "登录失败", "用户被禁用"},
	1005: {"400", "登录失败", "系统异常"},

	1100: {"400", "权限不足", "用户没有权限"}, // 11xx 权限相关

	2001: {"400", "请求参数不对"},
	2002: {"400", "请求参数不对"},
	2100: {"400", "购买失败", "购买失败"},
	2101: {"400", "清理失败", "刷新过期购买设备失败"},
	2201: {"400", "添加用户失败"},
	2202: {"400", "删除用户失败"},
	2203: {"400", "更新用户信息失败"},
}

type commResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	DevMessage string `json:"dev_message,omitempty"`
}

func (r *commResponse) StatusCode() (code int) {
	if r.Code == 0 {
		return
	}
	if v, ok := webErrorCode[r.Code]; ok {
		if vv, err := strconv.Atoi(v[0]); err == nil {
			code = vv
		} else {
			code = fiber.StatusBadRequest
		}
	} else {
		code = fiber.StatusBadRequest
	}
	return
}

func (r *commResponse) SetErrorCode(code int, errs ...error) {
	var msgs []string
	r.Code = code
	if r.Code == 0 {
		r.Message = "success"
	}
	if v, ok := webErrorCode[r.Code]; ok {
		r.Message = strings.Join(v[1:], ",")
	}
	for _, err := range errs {
		if err == nil {
			continue
		}
		msgs = append(msgs, err.Error())
	}
	if len(msgs) > 0 {
		r.DevMessage = strings.Join(msgs, ",")
	}
}
