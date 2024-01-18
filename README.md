# msync
A utility to sync files from disk or http or ftp to minio bucket. `msync` using redis list for managing sync source and destination.

## Usage

### Command line to run
```shell
msync -redis tcp://localhost:6379/3/SYNC_LIST_KEY -minio http://ACCESS_ID:ACCESS_SECRET@minio.myexample:9900/BUCKET_NAME
```

### Redis list lines
```redis
LPUSH SYNC_LIST_KEY SOURCE|DESTINATION
```
For example:
```redis
LPUSH SYNC_LIST_KEY /mnt/myvideos/12345.mp4|/myvideos/12345678.mp4
```
Or:
```redis
LPUSH SYNC_LIST_KEY http://v.myserver.com/videos/12345.mp4|/myvideos/12345678.mp4
```
Or:
```redis
LPUSH SYNC_LIST_KEY ftp://username:password@v.myserver.com/videos/12345.mp4|/myvideos/12345678.mp4
```