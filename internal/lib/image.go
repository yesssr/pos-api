package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func saveImageToCloud(ctx context.Context, client *s3.Client, bucket string, dir string, key string, imageData []byte) error {
  k := fmt.Sprintf("%s/%s", dir, key);
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(k),
		Body:   bytes.NewReader(imageData),
	});

	return err;
}

func DeleteImageFromCloud(ctx context.Context, client *s3.Client, bucket string, dir string, key string) error {
	k := fmt.Sprintf("%s/%s", dir, key);
	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(k),
	});

	return err;
}

func validateImgType(b []byte) error {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/jfif": true,
	}

	contentType := http.DetectContentType(b);

	if !allowedTypes[contentType] {
		return &AppError{
			Message: fmt.Sprintf("invalid file type: %s. Only JPG, PNG, GIF allowed", contentType),
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil;
}
const MaxSize = 10 << 20
func UploadHandler(r *http.Request, client *s3.Client, dir string, key string, keyDel *string) (string, error) {
	r.ParseMultipartForm(MaxSize)
	defaultURL := "http://default.png";

	mf := r.MultipartForm
	if mf == nil || mf.File == nil {
	  return defaultURL, nil
	}

	files := mf.File["image"]
	if len(files) == 0 {
	  return defaultURL, nil
	}

	fh := files[0]
	file, err := fh.Open()
	if err != nil {
	  return "", &AppError{
	  	Message: "Could not open file",
	   	StatusCode: http.StatusBadRequest,
	  }
	}
	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
	  return "", &AppError{
	  	Message: "Could not read file",
	  	StatusCode: http.StatusInternalServerError,
	  }
	}

	if len(data) > MaxSize {
    return "", &AppError{
    	Message: "file too large",
     	StatusCode: http.StatusBadRequest,
    }
}

	if err := validateImgType(data); err != nil {
	  return "", err
	}

	filename := key + "_" + GenerateUniqueNumber();
	bucket := os.Getenv("R2_BUCKET_NAME");
	if err := saveImageToCloud(r.Context(), client, bucket, dir, filename, data); err != nil {
	  return "", &AppError{
	  	Message: "Could not upload image to cloud",
	   	StatusCode: http.StatusInternalServerError,
	  }
	}

	if keyDel != nil {
		err := DeleteImageFromCloud(r.Context(), client, bucket, dir, *keyDel);
		if err != nil {
			fmt.Println("Warning: failed to delete old image:", err)
		}
	}

	accountId := os.Getenv("R2_ACCOUNT_ID");
	url := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", accountId, dir, filename)
	return url, nil
}
