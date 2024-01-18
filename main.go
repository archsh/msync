package main

import (
	"context"
	"flag"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"

	"github.com/minio/minio-go/v7"
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
		debug    bool
		serve    bool
	)
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.BoolVar(&serve, "serve", false, "serve as a server")
	flag.StringVar(&minioUrl, "minio", "", "minio url, eg: http://localhost:9000/bucketName")
	flag.StringVar(&redisUrl, "redis", "", "redis url, eg: tcp://localhost:6379/1/listKey")
	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	if redisUrl == "" {
		log.Fatal("redis Url can not be empty")
	}
	opts, e := ParseRedisURL(redisUrl)
	if nil != e {
		log.Fatalln(e)
	}
	var redisClient = redis.NewClient(&opts.Options)

	optm, e := ParseMinioURL(minioUrl)
	if nil != e {
		log.Fatalln(e)
	}

	s3Client, err := minio.New(optm.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(optm.Access, optm.Secret, ""),
		Secure: optm.Secure,
	})
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	if b, e := s3Client.BucketExists(ctx, optm.Bucket); nil != e {
		log.Fatalln("Check BucketExists failed:>", e)
	} else if !b {
		log.Fatalln("Bucket does not exists!", optm.Bucket)
	}
	for {
		llen := redisClient.LLen(opts.Key).Val()
		if llen < 1 {
			if serve {
				log.Infoln("No more lines in redis list, wait for a while ...")
				time.Sleep(time.Second * 3)
				continue
			} else {
				log.Infoln("End of redis list reached. Done.")
				return
			}
		}
		line, e := redisClient.LPop(opts.Key).Result()
		if nil != e {
			log.Infoln("End of redis list reached. Done.")
			return
		}
		if line == "" {
			continue
		}
		log.Debugln("Get Line via Redis:>", line)
		ss := strings.Split(line, "|")
		if len(ss) != 2 {
			log.Errorln("invalid line:>", line)
			continue
		}
		orig, dst := ss[0], ss[1]
		var contentType = "application/octet-stream"
		if size, e := FileSize(orig); nil != e {
			log.Errorln("Get FileSize failed:", orig, e)
			redisClient.LPush(opts.Key+"_FAILED", line)
		} else if info, e := s3Client.StatObject(ctx, optm.Bucket, dst, minio.StatObjectOptions{}); nil == e && info.Size == size {
			log.Infof("Object '%s' exists and match size (%d). Skipped.\n", path.Join(optm.Bucket, dst), info.Size)
			if info.ContentType != contentType {

			}
			continue
		} else if src, e := FetchFile(orig); nil != e {
			log.Errorln("Fetch original file failed:", src, e)
			redisClient.LPush(opts.Key+"_FAILED", line)
		} else if n, e := s3Client.FPutObject(ctx, optm.Bucket, dst, src, minio.PutObjectOptions{
			ContentType: contentType,
		}); nil != e {
			log.Errorln("FPutObject failed:", src, dst, e)
			redisClient.LPush(opts.Key+"_FAILED", line)
		} else {
			log.Infoln("Uploaded:>", orig, " to ", path.Join(n.Bucket, n.Key), " >>", n.Size, " bytes.")
		}
	}
}
