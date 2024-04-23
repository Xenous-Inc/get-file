package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

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

	resp, err := http.Get(s3Link)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Ошибка при загрузке файла: " + err.Error())
	}
	defer resp.Body.Close()

	chunkSize := int64(1024 * 1024) // 1 МБ
	reader := resp.Body

	var totalBytes int64
	for {
		buffer := make([]byte, chunkSize)
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		totalBytes += int64(n)
		// progress := float64(totalBytes) / float64(resp.ContentLength) * 100

		c.Write(buffer[:n])

		// progressJSON := fmt.Sprintf(`{"progress": %.2f}`, progress)
		// c.Write([]byte(progressJSON))

		buffer = make([]byte, chunkSize)
	}
	c.Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	c.Set("Transfer-Encoding", "chunked")
	return nil
}
