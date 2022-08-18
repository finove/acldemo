package server

import (
	"errors"
	"time"

	"github.com/finove/acldemo/model/user"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

/*
  登录
  授权中间件
  用户管理
*/

var (
	jwtTokenSecret        = []byte("hellosecret")
	jwtContextKey  string = "jwt_user_token"
)

type loginRequest struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

// Login 帐号登录
func Login(c *fiber.Ctx) (err error) {
	var req loginRequest
	var usr user.User
	var claims jwt.MapClaims = make(jwt.MapClaims)
	var token *jwt.Token
	var resp struct {
		commResponse
		Token string `json:"token,omitempty"`
	}
	// 解析参数
	if err = c.BodyParser(&req); err != nil {
		resp.SetErrorCode(2001, err)
		return c.Status(fiber.StatusUnauthorized).JSON(&resp)
	}
	if err = validate.Struct(&req); err != nil {
		resp.SetErrorCode(2001, err)
		return c.Status(fiber.StatusUnauthorized).JSON(&resp)
	}
	// 查找用户，验证密码
	if res := db.Where("name = ?", req.Username).First(&usr); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		resp.SetErrorCode(1003, res.Error)
		return c.Status(fiber.StatusForbidden).JSON(&resp)
	}
	if usr.Status == user.StatusDisabled {
		resp.SetErrorCode(1004)
		return c.Status(fiber.StatusForbidden).JSON(&resp)
	} else if err = usr.VerifyPassword(req.Password); err != nil {
		resp.SetErrorCode(1003, err)
		return c.Status(fiber.StatusForbidden).JSON(&resp)
	}
	// 生成令牌
	claims["id"] = usr.ID
	claims["user"] = usr.Name
	claims["nick_name"] = usr.NickName
	claims["expire"] = time.Now().Add(30 * time.Minute).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if resp.Token, err = token.SignedString(jwtTokenSecret); err != nil {
		resp.SetErrorCode(1005, err)
		err = c.Status(fiber.StatusInternalServerError).JSON(&resp)
		return
	}
	err = c.JSON(&resp)
	return
}

type newUserRequest struct {
	UserName string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func SessionInfo(c *fiber.Ctx) (err error) {
	var resp struct {
		commResponse
		User user.User `json:"user"`
	}
	resp.User.ID = c.Locals("login_uid").(uint)
	db.First(&resp.User)
	resp.User.Password = ""
	err = c.JSON(&resp)
	return
}

func FindUser(c *fiber.Ctx) (err error) {
	var usrs []user.User
	db.Find(&usrs)
	for k := range usrs {
		usrs[k].Password = ""
	}
	err = c.JSON(usrs)
	return
}

func NewUser(c *fiber.Ctx) error {
	var err error
	var req newUserRequest
	var resp commResponse
	var u *user.User
	if err = c.BodyParser(&req); err != nil {
		resp.SetErrorCode(2001, err)
		return c.JSON(&resp)
	}
	if err = validate.Struct(req); err != nil {
		resp.SetErrorCode(2001, err)
		return c.Status(fiber.StatusBadRequest).JSON(&resp)
	}
	u = user.NewUser(req.UserName, req.Password)
	if tx := db.Create(u); tx.Error != nil {
		resp.SetErrorCode(2201, tx.Error)
	}
	return c.JSON(&resp)
}

func GetUser(c *fiber.Ctx) error {
	var err error
	var resp struct {
		commResponse
		User *user.User `json:"user,omitempty"`
	}
	resp.User = new(user.User)
	if err = c.ParamsParser(resp.User); err != nil {
		resp.SetErrorCode(2001, err)
	} else {
		if tx := db.First(resp.User); tx.Error != nil {
			resp.SetErrorCode(1002, tx.Error)
		}
		resp.User.Password = ""
	}
	if resp.Code != 0 {
		resp.User = nil
	}
	return c.Status(resp.StatusCode()).JSON(&resp)
}

type updateUserRequest struct {
	NickName string `json:"nick_name"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
	Status   int    `json:"status"`
	Remark   string `json:"remark,omitempty"`
	Address  string `json:"address,omitempty"`
	Email    string `json:"email,omitempty"`
}

func UpdateUser(c *fiber.Ctx) error {
	var err error
	var req updateUserRequest
	var resp commResponse
	var u user.User
	if err = c.BodyParser(&req); err != nil {
		resp.SetErrorCode(2001, err)
	} else if err = c.ParamsParser(&u); err != nil {
		resp.SetErrorCode(2001, err)
	} else if tx := db.First(&u); tx.Error != nil {
		resp.SetErrorCode(1002, tx.Error)
	} else {
		if req.Password != "" {
			u.UpdatePassword(req.Password)
		}
		if req.NickName != "" {
			u.NickName = req.NickName
		}
		if req.Role != "" {
			u.Role = req.Role
		}
		if tx := db.Updates(&u); tx.Error != nil {
			resp.SetErrorCode(2203, tx.Error)
		}
	}

	return c.Status(resp.StatusCode()).JSON(&resp)
}

func RemoveUser(c *fiber.Ctx) error {
	var err error
	var resp commResponse
	var u user.User
	if err = c.ParamsParser(&u); err != nil {
		resp.SetErrorCode(2001, err)
	} else if tx := db.First(&u); tx.Error != nil {
		resp.SetErrorCode(1002, tx.Error)
	} else if tx := db.Delete(&u); tx.Error != nil {
		resp.SetErrorCode(2202, tx.Error)
	}
	return c.Status(resp.StatusCode()).JSON(&resp)
}

/* ---- 授权中间件 ---- */

// SetupJwtSecret 配置jwt加密密钥
func SetupJwtSecret(key string) {
	jwtTokenSecret = []byte(key)
}

// SetupJwtContextKey 配置jwt上下文字段名
func SetupJwtContextKey(key string) {
	jwtContextKey = key
}

// Protected protect routes
func Protected() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtTokenSecret,
		ContextKey:     jwtContextKey,
		SuccessHandler: jwtSuccess,
		ErrorHandler:   jwtError,
	})
}

func jwtSuccess(c *fiber.Ctx) (err error) {
	var ok bool
	var tk *jwt.Token
	var claims jwt.MapClaims
	if tk, ok = c.Locals(jwtContextKey).(*jwt.Token); !ok {
		return c.Next()
	}
	if claims, ok = tk.Claims.(jwt.MapClaims); !ok {
		return c.Next()
	}
	c.Locals("login_user", claims["user"])
	c.Locals("login_nickname", claims["nick_name"])
	if v, ex := claims["id"]; ex {
		var uid float64
		if uid, ok = v.(float64); ok {
			c.Locals("login_uid", uint(uid))
		}
	}
	if v, ex := claims["expire"]; ex {
		var expired float64
		if expired, ok = v.(float64); ok {
			c.Locals("login_expire", time.Unix(int64(expired), 0))
		}
		c.Locals("login_expire_unix", claims["expire"])
	}
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	var resp commResponse
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		resp.SetErrorCode(1001, err)

	} else {
		c.Status(fiber.StatusUnauthorized)
		resp.SetErrorCode(1001, err)
	}
	return c.JSON(&resp)
}
