package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func SaveImageToCloud(imageData []byte) (string, error) {
  var body bytes.Buffer
  writer := multipart.NewWriter(&body)

  part, err := writer.CreateFormFile("image", "upload.jpg")
  if err != nil {
      return "", err
  }
  part.Write(imageData)
  writer.Close()

  req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", &body)
  if err != nil {
    return "", err
  }
  req.Header.Set("Content-Type", writer.FormDataContentType())
  req.Header.Set("Authorization", "Client-ID "+os.Getenv("IMGUR_CLIENT_ID"))

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()

  respBody, err := io.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }

  if resp.StatusCode != http.StatusOK {
    return "", fmt.Errorf("upload failed: %s", respBody)
  }

  var result struct {
      Data struct {
          Link string `json:"link"`
      } `json:"data"`
  }

  if err := json.Unmarshal(respBody, &result); err != nil {
      return "", err
  }

  return result.Data.Link, nil
}
