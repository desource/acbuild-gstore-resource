package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/googleapi"
	storage "google.golang.org/api/storage/v1"
)

var OCTET_STREAM = googleapi.ContentType("application/octet-stream")

// TODO:
// - Add support for env variables

var (
	bucketName   = flag.String("bucket", "", "The name of an existing bucket within your project.")
	bucketPrefix = flag.String("prefix", "", "Optional bucket prefix.") // Fold into bucketName

	email      = flag.String("email", "", "service email")
	privateKey = flag.String("privateKey", "", "private key")

	compress = flag.Bool("compress", true, "Compress content")

	cacheControl = "private, max-age=0"
)

func parseFlags() {
	flag.Parse()

	if *bucketName == "" {
		log.Fatalf("Bucket argument is required. See --help.")
	}
	if *email == "" {
		log.Fatalf("Email argument is required. See --help.")
	}
	if *privateKey == "" {
		log.Fatalf("PrivateKey argument is required. See --help.")
	}
	if compress == nil {
		*compress = false
	}
	if flag.NArg() == 0 {
		log.Fatalf("Missing files. See --help.")
	}
}

func main() {

	parseFlags()

	service := loadService(context.Background())

	for _, filename := range flag.Args() {
		stat, err := os.Stat(filename)
		if err != nil {
			log.Fatalf("Error stating file: %v", filename, err)
		}

		if !stat.IsDir() {
			if *compress {
				uploadCompressedFile(service, filename, stat.Size())
			} else {
				uploadFile(service, filename, stat.Size())
			}
		} else {
			err = filepath.Walk(filename, func(filename string, stat os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !stat.IsDir() {
					if *compress {
						uploadCompressedFile(service, filename, stat.Size())
					} else {
						uploadFile(service, filename, stat.Size())
					}
				}
				return nil
			})
			if err != nil {
				log.Fatalf("Error processing path %q: %v", filename, err)
			}
		}
	}
}

func loadService(ctx context.Context) *storage.Service {

	key, err := ioutil.ReadFile(*privateKey)
	if err != nil {
		log.Fatalf("Unable to load creds: %v", err)
	}

	conf := &jwt.Config{
		Email:      *email,
		PrivateKey: []byte(key),
		Scopes:     []string{storage.DevstorageFullControlScope},
		TokenURL:   google.JWTTokenURL,
	}

	service, err := storage.New(conf.Client(ctx))
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	return service
}

func uploadFile(service *storage.Service, filename string, size int64) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %q: %v", filename, err)
	}

	if *bucketPrefix != "" {
		filename = path.Join(*bucketPrefix, filename)
	}
	p := path.Clean(filename)

	_, err = service.Objects.Insert(*bucketName, &storage.Object{
		Name:         p,
		Size:         uint64(size),
		CacheControl: cacheControl,
	}).Media(file, OCTET_STREAM).Do()

	if err != nil {
		log.Fatalf("Objects.Insert failed: %v", err)
	}
	fmt.Printf("↝ %s (%s)\n", p, formatSize(size))
}

func uploadCompressedFile(service *storage.Service, filename string, size int64) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %q: %v", filename, err)
	}

	if *bucketPrefix != "" {
		filename = path.Join(*bucketPrefix, filename)
	}
	p := path.Clean(filename)

	pr, pw := io.Pipe()

	go func() {
		gz := gzip.NewWriter(pw)
		defer pw.Close()

		if _, err := io.Copy(gz, file); err != nil {
			log.Fatalf("Compress failed: %v", err)
		}
		if err := gz.Flush(); err != nil {
			log.Fatalf("Compress failed: %v", err)
		}
		if err := gz.Close(); err != nil {
			log.Fatalf("Compress failed: %v", err)
		}
	}()

	_, err = service.Objects.Insert(*bucketName, &storage.Object{
		Name:            p,
		ContentEncoding: "gzip",
		CacheControl:    cacheControl,
	}).Media(pr, OCTET_STREAM).Do()

	if err != nil {
		log.Fatalf("Objects.Insert failed: %v", err)
	}
	fmt.Printf("↝ %s (%s)\n", p, formatSize(size))
}

func formatSize(s int64) string {
	const (
		_          = iota // ignore first value by assigning to blank identifier
		kb float64 = 1 << (10 * iota)
		mb
		gb
		tb
	)
	b := float64(s)
	switch {
	case b >= tb:
		return fmt.Sprintf("%.2fTB", b/tb)
	case b >= gb:
		return fmt.Sprintf("%.2fGB", b/gb)
	case b >= mb:
		return fmt.Sprintf("%.2fMB", b/mb)
	case b >= kb:
		return fmt.Sprintf("%.2fKB", b/kb)
	default:
		return fmt.Sprintf("%dB", s)
	}
}
