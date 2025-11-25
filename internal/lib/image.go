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

func saveImageToCloud(ctx context.Context, client *s3.Client, bucket, dir, key, contentType string, imageData []byte) error {
  k := fmt.Sprintf("%s/%s", dir, key);
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(k),
		Body:   bytes.NewReader(imageData),
		ContentType: aws.String(contentType),
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

func validateImgType(b []byte) (string, error) {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/jfif": true,
	}

	contentType := http.DetectContentType(b);

	if !allowedTypes[contentType] {
		return "", &AppError{
			Message: fmt.Sprintf("invalid file type: %s. Only JPG, PNG, GIF allowed", contentType),
			StatusCode: http.StatusBadRequest,
		}
	}

	return contentType, nil;
}
const MaxSize = 10 << 20
func UploadHandler(r *http.Request, client *s3.Client, bucket, dir, key string, keyDel *string) (string, error) {
	r.ParseMultipartForm(MaxSize)
	mf := r.MultipartForm
	if mf == nil || mf.File == nil {
	  return "", nil
	}

	files := mf.File["image"]
	if len(files) == 0 {
	  return "", nil
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

	cType, err := validateImgType(data);
	if err != nil {
	  return "", err;
	}
	filename := key + "_" + GenerateUniqueNumber();
	if err := saveImageToCloud(r.Context(), client, bucket, dir, filename, cType, data); err != nil {
	  return "", &AppError{
	  	Message: err.Error(),
	   	StatusCode: http.StatusInternalServerError,
	  }
	}

	if keyDel != nil {
		err := DeleteImageFromCloud(r.Context(), client, bucket, dir, *keyDel);
		if err != nil {
			fmt.Println("Warning: failed to delete old image:", err)
		}
	}

	baseUrl := os.Getenv("PUBLIC_ENDPOINT_URL");
	url := fmt.Sprintf("%s/%s/%s", baseUrl, dir, filename);
	return url, nil
}
