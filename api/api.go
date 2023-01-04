package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gen1us2k/analytics-writer/config"
	"github.com/gen1us2k/analytics-writer/event"
	"github.com/gin-gonic/gin"
)

type (
	Server struct {
		router *gin.Engine
		config *config.Config
		s3     *s3manager.Uploader
	}
	Message struct {
	}
)

func (m *Message) fromEvent(e *event.Event) {
}
func New(c *config.Config) (*Server, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	s := &Server{
		router: gin.Default(),
		config: c,
		s3:     s3manager.NewUploader(sess),
	}
	s.init()
	return s, nil
}
func (s *Server) init() {
	s.router.POST("/event", s.handleEvent)
}
func (s *Server) handleEvent(c *gin.Context) {
	var e event.Event
	if err := c.BindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse input"})
		return
	}
	m := &Message{}
	m.fromEvent(&e)
	payload, err := json.Marshal(m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed encoding message"})
		return
	}
	_, err = s.s3.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String("something"),
		Body:   bytes.NewBuffer(payload),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed storing the event"})
		return
	}
	c.Status(http.StatusNoContent)

}
