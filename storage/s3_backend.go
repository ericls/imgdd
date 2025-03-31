package storage

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"io"

	"github.com/ericls/imgdd/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	lru "github.com/hashicorp/golang-lru/v2"
)

const defaultRegion = "us-east-1"

type S3Client struct {
	uploader *s3manager.Uploader
	s3       *s3.S3
}

// JSON config of S3 storage
type S3StorageConfig struct {
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
	Access   string `json:"access"`
	Secret   string `json:"secret"`
}

func (conf S3StorageConfig) Hash() uint32 {
	h := fnv.New32a()
	s := conf.Endpoint + "|" + conf.Bucket + "|" + conf.Access + "|" + conf.Secret
	h.Write([]byte(s))
	return h.Sum32()
}

func (conf S3StorageConfig) ToJSON() []byte {
	data, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	return data
}

type S3StorageBackend struct {
	cache *lru.TwoQueueCache[uint32, *S3Storage]
}

func (s *S3StorageBackend) FromJSONConfig(config []byte) (Storage, error) {
	var conf S3StorageConfig
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return nil, err
	}
	hash := conf.Hash()
	if storage, ok := s.cache.Get(hash); ok {
		return storage, nil
	}
	store := S3Storage{
		endpoint: conf.Endpoint,
		bucket:   conf.Bucket,
		access:   conf.Access,
		secret:   conf.Secret,
	}
	store.ensureClient()
	s.cache.Add(hash, &store)
	return &store, nil
}

func (s *S3StorageBackend) ValidateJSONConfig(config []byte) error {
	var conf S3StorageConfig
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return err
	}
	if conf.Endpoint == "" || conf.Bucket == "" || conf.Access == "" || conf.Secret == "" {
		return errors.New("invalid S3 storage config")
	}
	return nil
}

type S3Storage struct {
	endpoint string
	bucket   string
	access   string
	secret   string
	client   *utils.Lazy[*S3Client]
}

func (s *S3Storage) ensureClient() {
	if s.client != nil {
		return
	}
	s.client = utils.NewLazy(func() *S3Client {
		var staticResolver endpoints.ResolverFunc = func(service, region string, opts ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
			return endpoints.ResolvedEndpoint{
				URL:         s.endpoint,
				PartitionID: "",
			}, nil
		}
		cred := credentials.NewStaticCredentials(s.access, s.secret, "")
		cfg := aws.Config{
			Region:           aws.String(defaultRegion),
			EndpointResolver: staticResolver,
			Credentials:      cred,
		}
		cfg.WithS3ForcePathStyle(true)

		sess := session.Must(session.NewSession(&cfg))

		uploader := s3manager.NewUploader(sess)
		s3Client := s3.New(sess)
		return &S3Client{
			uploader: uploader,
			s3:       s3Client,
		}
	})
}

func (s *S3Storage) GetReader(filename string) io.ReadCloser {
	service := s.client.Value().s3
	res, err := service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil
	}
	return res.Body
}

func (s *S3Storage) Save(file utils.SeekerReader, filename string, mimeType string) error {
	_, err := s.client.Value().uploader.Upload(&s3manager.UploadInput{
		Bucket:      &s.bucket,
		Body:        file,
		Key:         aws.String(filename),
		ContentType: aws.String(mimeType),
	})
	return err
}

func (s *S3Storage) GetMeta(filename string) FileMeta {
	meta, err := s.client.Value().s3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return FileMeta{
			ByteSize:    0,
			ContentType: "",
		}
	}
	return FileMeta{
		ByteSize:    *meta.ContentLength,
		ContentType: *meta.ContentType,
		ETag:        *meta.ETag,
	}
}

func (s *S3Storage) Delete(filename string) error {
	_, err := s.client.Value().s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	return err
}

func (s *S3Storage) CheckConnection() error {
	client := s.client.Value()
	_, err := client.s3.ListBuckets(&s3.ListBucketsInput{})
	return err
}

func (s *S3Storage) CreateBucket(name string) error {
	_, err := s.client.Value().s3.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	return err
}
