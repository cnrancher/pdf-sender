package types

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var emailLock = &sync.Once{}

type kindDocs struct {
	kind
	docs []document
}

var cache = map[string]*kindDocs{}

func InitEmailBody(ctx *cli.Context) error {
	if OSSBucketEndpoint != "" {
		logrus.Info("Using oss as file host url")
	}
	for kind := range Config.Kinds {
		cache[kind] = &kindDocs{kind: Config.Kinds[kind]}
	}
	for i, doc := range Config.Documents {
		if err := doc.Validate(); err != nil {
			return err
		}
		for _, kind := range doc.Kind {
			content, ok := cache[kind]
			if !ok {
				logrus.Warnf("kind %s for document %s configuration is not validated, going to ignore", kind, doc.Name)
				continue
			}
			content.docs = append(content.docs, Config.Documents[i])
		}
	}

	return nil
}

func GetBodyByKind(kind string) (string, error) {
	isPublic := true

	content, ok := cache[kind]
	if !ok {
		return "", fmt.Errorf("kind %s is not supported", kind)
	}

	if OSSBucketEndpoint != "" {
		info, err := GetBucketInfo()
		if err != nil {
			logrus.Warnf("failed to get bucket info, assuming public access, error: %v", err)
		} else {
			isPublic = IsBucketPrivate(info.ACL)
		}
	}

	writer := bytes.NewBufferString(content.Header)
	for _, d := range content.docs {
		line := d.GetLine(isPublic) + "\n"
		writer.WriteString(line)
	}
	writer.WriteString(content.Footer)
	return writer.String(), nil
}

func GetSenderNameAndSubjectByKind(kind string) (string, string, error) {
	content, ok := cache[kind]
	if !ok {
		return "", "", fmt.Errorf("kind %s is not supported", kind)
	}
	return content.Subject, content.SenderName, nil
}
