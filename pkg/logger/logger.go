package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ILogger interface {
	Print() ILogger
	Save()
	setQuery(c *fiber.Ctx)
	setBody(c *fiber.Ctx)
}

type logger struct {
	Time       string
	Ip         string
	Path       string
	Method     string
	StatusCode int
	Query      any
	Body       any
	Response   any
}

func NewLogger(c *fiber.Ctx, code int, res any) ILogger {
	log := &logger{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Path:       c.Path(),
		Method:     c.Method(),
		StatusCode: code,
		Response:   res,
	}
	log.setQuery(c)
	log.setBody(c)

	return log
}

func (l *logger) Print() ILogger {
	utils.Debug(l)
	return l
}

func (l *logger) Save() {
	data := utils.Output(l)

	filename := fmt.Sprintf("./assets/logs/%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("open log file failed: ", err)
	}
	defer file.Close()

	file.Write(data)
	file.WriteString("\n")
}

func (l *logger) setQuery(c *fiber.Ctx) {
	query := new(any)
	if err := c.QueryParser(query); err != nil {
		log.Println("query parser error: ", err)
	}
	l.Query = query
}

func (l *logger) setBody(c *fiber.Ctx) {
	body := new(any)
	if err := c.BodyParser(body); err != nil {
		log.Println("body parser error: ", err)
	}
	l.Body = body

}
