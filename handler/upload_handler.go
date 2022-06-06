package handler

import (
	"awesomeProject/service"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"sync"
)

type csvUploadInput struct {
	CsvFile *multipart.FileHeader `form:"file" binding:"required"`
}

type UploadHandler struct {
	hcService service.HealthCheckService
}

func NewUploadHandler(h service.HealthCheckService) *UploadHandler {
	return &UploadHandler{
		hcService: h,
	}
}

func (u *UploadHandler) Upload(c *gin.Context) {
	var wg sync.WaitGroup

	var req csvUploadInput
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "file not support",
		})
		return
	}

	file, err := req.CsvFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "file cannot open",
		})
	}

	defer file.Close()
	records, err := u.hcService.ReadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "file cannot read",
		})
		return
	}

	chInput := make(chan string)
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go u.hcService.Checker(i, chInput, &wg) // consumer -> queue chInput
	}

	for _, v := range records {
		chInput <- v[0] // publish queue
	}
	close(chInput)
	wg.Wait()

	c.JSON(200, gin.H{
		"data": u.hcService.Result(),
	})
}
