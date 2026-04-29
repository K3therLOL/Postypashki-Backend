package manager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	repository "task/repository/rai"
	"time"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"

	"github.com/gen2brain/avif"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type API struct {
	repo     *repository.Repository
	client   *http.Client
	logger   *log.Logger
	s3Client *s3.Client
}

type imgUrl struct {
	Url string `json:"url"`
}

var (
	ErrWrongTaskID  = errors.New("Invalid task_id.")
	ErrStatusAccess = errors.New("Could not get status.")
	ErrSaveTaskID   = errors.New("Could not save task_id.")
	ErrNoConnection = errors.New("Check connection to internet.")
	ErrWithResp     = errors.New("Response ended with error.")
	ErrNotImage     = errors.New("No image returned from response.")
	ErrImgProc      = errors.New("Image processing failed.")
	ErrReqCreation  = errors.New("Could not create a request handler.")
)

func confS3Client() *s3.Client {
	secretKey := os.Getenv("SECRET_KEY")
	accessKey := os.Getenv("ACCESS_KEY")
	endpoint := os.Getenv("ENDPOINT")
	region := os.Getenv("REGION")

	log.Println(region)
	staticCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(staticCreds),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
		//o.ClientLogMode = aws.LogSigning | aws.LogRetries | aws.LogRequestWithBody
	})

	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range result.Buckets {
		log.Printf("* %s (created %s)\n", aws.ToString(bucket.Name), bucket.CreationDate)
	}

	return client
}

func NewAPI() *API {
	api := &API{
		repo: repository.NewRepository(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:   log.New(os.Stdout, "api: ", log.Ldate|log.Ltime),
		s3Client: confS3Client(),
	}
	return api
}

func getDomainName(url string) string {
	domainString, _ := strings.CutPrefix(url, "https://")
	index := strings.Index(domainString, "/")
	return domainString[:index]
}

func (api *API) downloadAndDecodeImg(url string) (image.Image, string, error) {
	api.logger.Println(url)
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_144),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	req, err := fhttp.NewRequest(fhttp.MethodGet, url, nil)
	if err != nil {
		return nil, "", ErrReqCreation
	}

	api.logger.Println(getDomainName(url))
	//req.Header.Set("Host", "cdn.plus.unsplash.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:149.0) Gecko/20100101 Firefox/149.0")
	req.Header.Set("Accept", "image/avif,image/webp,image/png,image/svg+xml,image/*;q=0.8,*/*;q=0.5")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	//req.Header.Set("Referer", "https://plus.unsplash.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", ErrNoConnection
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", ErrWithResp
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, "", ErrNotImage
	}
	api.logger.Println(contentType)

	img, format, err := image.Decode(resp.Body)
	api.logger.Printf("Image format: %v\n", format)

	if err != nil {
		return nil, "", fmt.Errorf("Cannot decode image with %v format.", format)
	}

	api.logger.Println("Image successfully decoded")
	return img, format, nil
}

func imageToReader(img image.Image, format string) io.Reader {
	buf := new(bytes.Buffer)
	if format == "png" {
		png.Encode(buf, img)
	} else if format == "jpeg" {
		jpeg.Encode(buf, img, nil)
	} else if format == "gif" {
		gif.Encode(buf, img, nil)
	} else if format == "avif" {
		avif.Encode(buf, img)
	} else {
		log.Fatal("No valid encoding found.")
	}

	return buf
}

func generatePublicURL(bucket, key string) string {
	return fmt.Sprintf("https://gateway.storjshare.io/%s/%s", bucket, key)
}

func (api *API) generateTemporaryUrl(bucket, key string) (string, error) {
	presignClient := s3.NewPresignClient(api.s3Client)

	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Generating url for 24 hours
	presignedRequest, err := presignClient.PresignGetObject(context.TODO(), presignParams, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(24 * time.Hour)
	})

	if err != nil {
		return "", err
	}

	return presignedRequest.URL, nil
}

func (api *API) uploadImgForClient(img image.Image, format string) (string, error) {
	transferClient := transfermanager.New(api.s3Client)

	key := fmt.Sprintf("image.%v", format)
	bucket := os.Getenv("BUCKET")
	_, err := transferClient.UploadObject(context.TODO(), &transfermanager.UploadObjectInput{
		Bucket:            aws.String(bucket),
		Key:               aws.String(key),
		Body:              imageToReader(img, format),
		ContentType:       aws.String(fmt.Sprintf("image/%v", format)),
		ChecksumAlgorithm: "",
	})
	if err != nil {
		return "", err
	}

	tempImgUrl, err := api.generateTemporaryUrl(bucket, key)
	if err != nil {
		api.logger.Println("Cannot generate temporary url, returning fixed.")
		return generatePublicURL(bucket, key), nil
	}
	return tempImgUrl, nil
}

func (api *API) processTask(uuid uuid.UUID, url string) {
	img, format, err := api.downloadAndDecodeImg(url)
	if err != nil {
		api.logger.Println(err)
		return // ? should i do change that
	}

	imgUrl, err := api.uploadImgForClient(img, format)
	if err != nil {
		api.logger.Println("Hello:")
		api.logger.Println(err)
		return // ? should i do change that
	}

	api.logger.Printf("Image uploaded to %v\n.", imgUrl)

	api.repo.Update(uuid)
}

// Image processing
func (api *API) ExecuteTask(c *gin.Context) {
	taskID := uuid.New()

	req := imgUrl{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	go api.processTask(taskID, req.Url)

	if err := api.repo.Save(taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": taskID.String(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID.String(),
	})
}

func (api *API) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	uuid, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrWrongTaskID.Error(),
		})
		return
	}
	taskobj, err := api.repo.Get(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrStatusAccess.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": taskobj.Status,
	})
}

/* TO DO */
func (api *API) GetTaskResult(c *gin.Context) {
	taskID := c.Param("task_id")
	c.JSON(http.StatusOK, gin.H{
		"result":  "soon here will be result",
		"task_id": taskID,
	})
}
