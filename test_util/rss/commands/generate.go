package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cloudfoundry/gorouter/route_service"
	"github.com/cloudfoundry/gorouter/test_util/rss/common"
	"github.com/codegangsta/cli"
)

func GenerateSignature(c *cli.Context) {
	crypto, err := common.CreateCrypto(c)
	if err != nil {
		os.Exit(1)
	}

	signature, err := createSigFromArgs(c)
	if err != nil {
		os.Exit(1)
	}

	sigEncoded, metaEncoded, err := route_service.BuildSignatureAndMetadata(crypto, &signature)
	if err != nil {
		fmt.Printf("Failed to create signature: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Encoded Signature:\n%s\n\n", sigEncoded)
	fmt.Printf("Encoded Metadata:\n%s\n\n", metaEncoded)
}

func createSigFromArgs(c *cli.Context) (route_service.Signature, error) {
	signature := route_service.Signature{}
	url := c.String("url")

	if url == "" {
		cli.ShowCommandHelp(c, "generate")
		return signature, errors.New("url is required")
	}

	var sigTime time.Time

	timeStr := c.String("time")

	if timeStr != "" {
		unix, err := strconv.ParseInt(timeStr, 10, 64)
		if err != nil {
			fmt.Printf("Invalid time format: %s", timeStr)
			return signature, err
		}

		sigTime = time.Unix(unix, 0)
	} else {
		sigTime = time.Now()
	}

	return route_service.Signature{
		RequestedTime: sigTime,
		ForwardedUrl:  url,
	}, nil
}
