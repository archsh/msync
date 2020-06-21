package main

import (
	"flag"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"

	"github.com/minio/minio-go/v6"

	"github.com/cavaliercoder/grab"
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
	if redisUrl == "" {
		log.Fatal("minio Url can not be empty")
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
	s3Client, err := minio.New(optm.Endpoint, optm.Access, optm.Secret, optm.Secure)
	if err != nil {
		log.Fatalln(err)
	}
	if b, e := s3Client.BucketExists(optm.Bucket); nil != e {
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
		var contentType = "application/octet-stream"
		if resp, e := http.Head(ss[0]); nil != e {
			log.Errorln("Get HEAD failed:>", ss[0], e)
			continue
		} else {
			if s := resp.Header.Get("Content-Type"); s != "" {
				contentType = s
			}
			if info, e := s3Client.StatObject(optm.Bucket, ss[1], minio.StatObjectOptions{});
				nil == e && info.Size == resp.ContentLength {
				log.Infof("Object '%s' exists and match size (%d). Skipped.\n", path.Join(optm.Bucket, ss[1]), info.Size)
				if info.ContentType != contentType {

				}
				continue
			}
		}
		log.Infoln("Downloading:>", ss[0], contentType, " ...")
		if resp, e := grab.Get(".", ss[0]); nil != e {
			log.Errorln("Download failed:>", ss[0], e)
		} else {
			log.Infoln("Downloaded:>", resp.Filename)
			if n, e := s3Client.FPutObject(optm.Bucket, ss[1], resp.Filename, minio.PutObjectOptions{
				ContentType: contentType,
			}); nil != e {
				log.Errorln("upload to minio failed:>", ss[1], e)
				log.Println(">>", line, "<<")
			} else {
				log.Infoln("Uploaded:>", resp.Filename, " to ", path.Join(optm.Bucket, ss[1]), " >>", n, " bytes.")
			}
			if e := os.Remove(resp.Filename); nil != e {
				log.Errorln("Delete file failed:", resp.Filename, e)
			}
		}
	}
}
