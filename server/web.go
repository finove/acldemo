package server

import (
	"errors"

	"github.com/finove/acldemo/model/device"
	"github.com/finove/acldemo/model/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var validate = validator.New() // Create Validate for using.
var db *gorm.DB
var app *fiber.App

func Run() {
	var err error
	if err = connectDB(); err != nil {
		log.Error().Err(err).Msg("failed to connect database")
		return
	}
	app = fiber.New()
	app.Use(
		logger.New(), // add simple logger
		cors.New(),
	)
	app.Use(favicon.New())
	app.Static("/", "./static")

	app.Post("/login", Login)
	app.Post("/logout", Login)
	apiV1 := app.Group("/v1")
	apiV1.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code":   0,
			"status": "success",
			"versions": fiber.Map{
				"cloudcms": "0.0.3",
			},
		})
	})
	apiUser := apiV1.Group("/user", Protected())
	apiUser.Get("/", FindUser)
	apiUser.Post("/", NewUser)
	apiUser.Get("/session", SessionInfo)
	if err = app.Listen(":3000"); err != nil {
		log.Fatal().Err(err).Msg("run web server")
	}
}

func connectDB() (err error) {
	var usr user.User
	if db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{}); err != nil {
		return
	}
	db.AutoMigrate(&device.Device{}, &user.User{})
	res := db.Where("name = ?", "sysadmin").First(&usr)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		usr.Role = user.RoleSysadmin
		usr.NickName = "Super"
		usr.Profile.Remark = "系统管理员"
		usr.UpdatePassword("helloadmin")
		db.Create(&usr)
		log.Info().Msg("init sysadmin user")
	}
	return
}
