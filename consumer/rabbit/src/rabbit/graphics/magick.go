package graphics

import (
  "fmt"
  "path"
  "strings"
  "log"
  "os"
  "os/exec"
  "io"
  "context"
  "cloud.google.com/go/storage"
)

func Strcat(a string, b string) string {
  var buf strings.Builder
  fmt.Fprintf(&buf, "%s%s", a, b)
  return buf.String()
}

func Resize(srcId string) error {
  srcName := Strcat(srcId, ".jpeg")
  ogName := Strcat(srcId, "_og.jpeg")
  lgName := Strcat(srcId, "_lg.jpeg")
  mdName := Strcat(srcId, "_md.jpeg")
  smName := Strcat(srcId, "_sm.jpeg")

  srcPath := path.Join("/tmp", srcName)
  ogPath := path.Join("/tmp", ogName)
  lgPath := path.Join("/tmp", lgName)
  mdPath := path.Join("/tmp", mdName)
  smPath := path.Join("/tmp", smName)


  srcFile, err := os.OpenFile(srcPath, os.O_RDWR|os.O_CREATE, 0600)
  if err != nil {
    log.Printf("Failed to open '%s' for writing: ", srcPath, err)
    return err
  }
  defer srcFile.Close()
  defer os.Remove(srcPath)

  ctx := context.Background()
  client, err := storage.NewClient(ctx)
  if err != nil {
    log.Printf("Failed to create client: %v", err)
    return err
  }
  defer client.Close()

  bucketName := os.Getenv("STORAGE_BUCKET")
  bucket := client.Bucket(bucketName)

  srcReader, err := bucket.Object(srcName).NewReader(ctx)
  if err != nil {
    log.Printf("Failed to create reader for object: %v", err)
    return err
  }
  defer srcReader.Close()

  if _, err := io.Copy(srcFile, srcReader); err != nil {
    log.Printf("Failed to make local copy of source image: %v", err)
    return err
  }

  srcReader.Close()
  srcFile.Close()


  // TODO: Customize compression options could be added.
  cmd := exec.Command("gm", "convert",
                      /* [options...], path */
                      "+profile", "*", srcPath, /* strip EXIF data from src */
                      ogPath,
                      "-resize", "1080x1080", lgPath,
                      "-resize", "612x612", mdPath,
                      "-resize", "161x161",  smPath)
  if err := cmd.Run(); err != nil {
    log.Printf("Failed to process image: %v", err)
    return err
  }
  defer os.Remove(ogPath)
  defer os.Remove(lgPath)
  defer os.Remove(mdPath)
  defer os.Remove(smPath)


  ogFile, err := os.Open(ogPath)
  if err != nil {
    log.Printf("Failed to open file: %v", err)
    return err
  }
  defer ogFile.Close()

  ogWriter := bucket.Object(ogName).NewWriter(ctx)
  defer ogWriter.Close()

  if _, err := io.Copy(ogWriter, ogFile); err != nil {
    log.Printf("Failed to write object: %v", err)
    return err
  }

  ogWriter.Close()
  ogFile.Close()


  lgFile, err := os.Open(lgPath)
  if err != nil {
    log.Printf("Failed to open file: %v", err)
    return err
  }
  defer lgFile.Close()

  lgWriter := bucket.Object(lgName).NewWriter(ctx)
  defer lgWriter.Close()

  if _, err := io.Copy(lgWriter, lgFile); err != nil {
    log.Printf("Failed to write object: %v", err)
    return err
  }

  lgWriter.Close()
  lgFile.Close()


  mdFile, err := os.Open(mdPath)
  if err != nil {
    log.Printf("Failed to open file: %v", err)
    return err
  }
  defer mdFile.Close()

  mdWriter := bucket.Object(mdName).NewWriter(ctx)
  defer mdWriter.Close()

  if _, err := io.Copy(mdWriter, mdFile); err != nil {
    log.Printf("Failed to write object: %v", err)
    return err
  }

  mdWriter.Close()
  mdFile.Close()


  smFile, err := os.Open(smPath)
  if err != nil {
    log.Printf("Failed to open file: %v", err)
    return err
  }
  defer smFile.Close()

  smWriter := bucket.Object(smName).NewWriter(ctx)
  defer smWriter.Close()

  if _, err := io.Copy(smWriter, smFile); err != nil {
    log.Printf("Failed to write object: %v", err)
    return err
  }

  smWriter.Close()
  smFile.Close()

  return nil
}
