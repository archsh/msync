package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"

	"github.com/minio/minio-go/v6"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-bucketname, my-objectname
	// and my-filename.csv are dummy values, please replace them with original values.

	// Requests are always secure (HTTPS) by default. Set secure=false to enable insecure (HTTP) access.
	// This boolean value is the last argument for New().

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	var (
		minioUrl string
		redisUrl string
	)
	flag.StringVar(&minioUrl, "minio", "http://localhost:9000", "minio")
	flag.StringVar(&redisUrl, "redis", "tcp://localhost:6379/1", "redis")
	flag.Parse()
	opts, e := ParseRedisURL(redisUrl)
	if nil != e {
		log.Fatalln(e)
	}
	var redisClient = redis.NewClient(opts)
	redisClient.LPop("").Result()
	optm, e := ParseMinioURL(minioUrl)
	if nil != e {
		log.Fatalln(e)
	}
	s3Client, err := minio.New(optm.Endpoint, optm.Access, optm.Secret, optm.Secure)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := s3Client.FPutObject("my-bucketname", "my-objectname", "my-filename.csv", minio.PutObjectOptions{
		ContentType: "application/csv",
	}); err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully uploaded my-filename.csv")
}
