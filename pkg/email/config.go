package email

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var emailLock = &sync.Once{}

type kindDocs struct {
	kind types.Kind
	docs []types.Document
}

var cache = map[string]*kindDocs{}

func InitEmailBody(ctx *cli.Context) error {
	if types.OSSBucketEndpoint != "" {
		logrus.Info("Using oss as file host url")
	}
	for kind := range types.Config.Kinds {
		cache[kind] = &kindDocs{kind: types.Config.Kinds[kind]}
	}
	for i, doc := range types.Config.Documents {
		if err := doc.Validate(); err != nil {
			return err
		}
		for _, kind := range doc.Kind {
			content, ok := cache[kind]
			if !ok {
				logrus.Warnf("kind %s for document %s configuration is not validated, going to ignore", kind, doc.Name)
				continue
			}
			content.docs = append(content.docs, types.Config.Documents[i])
		}
	}
	if types.Register.Template != "" {
		if _, err := RegisterTemplate.Parse(types.Register.Template); err != nil {
			return err
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

	if types.OSSBucketEndpoint != "" {
		info, err := types.GetBucketInfo()
		if err != nil {
			logrus.Warnf("failed to get bucket info, assuming public access, error: %v", err)
		} else {
			isPublic = types.IsBucketPrivate(info.ACL)
		}
	}

	writer := bytes.NewBufferString(content.kind.Header)
	for _, d := range content.docs {
		line := d.GetLine(isPublic) + "\n"
		writer.WriteString(line)
	}
	writer.WriteString(content.kind.Footer)
	return writer.String(), nil
}

func GetSenderNameAndSubjectByKind(kind string) (string, string, error) {
	content, ok := cache[kind]
	if !ok {
		return "", "", fmt.Errorf("kind %s is not supported", kind)
	}
	return content.kind.Subject, content.kind.SenderName, nil
}
