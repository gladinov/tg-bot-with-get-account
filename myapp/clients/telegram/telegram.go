package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"main.go/lib/e"
	"main.go/service/service_models"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod     = "getUpdates"
	sendUpdateMethod     = "sendMessage"
	sendPhotoMethod      = "sendPhoto"
	sendMediaGroupMethod = "sendMediaGroup"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can`t get updates", err) }()
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err

	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendUpdateMethod, q)
	if err != nil {
		return e.Wrap("can`t send message", err)
	}
	return nil

}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can`t do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}

func (c *Client) SendImageFromBuffer(chatID int, imageData []byte, caption string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("chat_id", strconv.Itoa(chatID))
	if caption != "" {
		writer.WriteField("caption", caption)
	}

	part, err := writer.CreateFormFile("photo", "image.png")
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, bytes.NewReader(imageData)); err != nil {
		return err
	}

	writer.Close()

	_, err = c.doMultipartRequest("sendPhoto", body, writer.FormDataContentType())
	return err
}

func (c *Client) SendMediaGroupFromBuffer(chatID int, images []*service_models.ImageData) error {
	if len(images) == 0 {
		return errors.New("no images to send")
	}

	// Ограничение Telegram
	if len(images) > 10 {
		images = images[:10]
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем chat_id
	writer.WriteField("chat_id", strconv.Itoa(chatID))

	// Подготавливаем медиа-группу
	media := make([]map[string]string, len(images))
	for i, img := range images {
		media[i] = map[string]string{
			"type":  "photo",
			"media": "attach://image_" + strconv.Itoa(i),
		}

		// Подпись только для первого изображения
		if i == 0 && img.Caption != "" {
			media[i]["caption"] = img.Caption
		}
	}

	mediaJSON, _ := json.Marshal(media)
	writer.WriteField("media", string(mediaJSON))

	// Добавляем изображения из буферов
	for i, img := range images {
		part, err := writer.CreateFormFile("image_"+strconv.Itoa(i), img.Name)
		if err != nil {
			return fmt.Errorf("can't create form file: %v", err)
		}

		// Копируем байты из буфера
		if _, err := io.Copy(part, bytes.NewReader(img.Data)); err != nil {
			return fmt.Errorf("can't copy image data: %v", err)
		}
	}

	writer.Close()

	_, err := c.doMultipartRequest("sendMediaGroup", body, writer.FormDataContentType())
	if err != nil {
		return fmt.Errorf("can't send media group: %v", err)
	}

	return nil
}

func (c *Client) doMultipartRequest(method string, body *bytes.Buffer, contentType string) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do multipart request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	bodyResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram API error: %s - %s", resp.Status, string(bodyResponse))
	}

	return bodyResponse, nil
}
