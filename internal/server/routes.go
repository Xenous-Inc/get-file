package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("b64/", s.B64Hanlder)
	s.App.Get("chunk/", s.GetFileHandler)

}

func (s *FiberServer) B64Hanlder(c *fiber.Ctx) error {
	s3Link := c.Query("link")
	if s3Link == "" {
		return c.Status(http.StatusBadRequest).SendString("Не указан URL файла")
	}

	resp, err := http.Get(s3Link)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Ошибка при загрузке файла: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении данных из ответа:", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	type response struct {
		Data string `json:"data"`
	}

	encodedData := base64.StdEncoding.EncodeToString(body)

	c.Set("Content-Type", "application/json")
	c.Status(http.StatusOK).JSON(response{
		Data: encodedData,
	})
	return nil
}
func (s *FiberServer) GetFileHandler(c *fiber.Ctx) error {
	s3Link := c.Query("link")
	if s3Link == "" {
		return c.Status(http.StatusBadRequest).SendString("Не указан URL файла")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(s3Link)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Ошибка при загрузке файла: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).SendString("Ошибка при загрузке файла: " + resp.Status)
	}

	chunkSize := int64(1024 * 1024) // 1 МБ
	reader := resp.Body

	c.Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	c.Set("Transfer-Encoding", "chunked")

	var totalBytes int64
	buffer := make([]byte, chunkSize)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		totalBytes += int64(n)
		if _, err := c.Write(buffer[:n]); err != nil {
			return err
		}
	}

	return nil
}
