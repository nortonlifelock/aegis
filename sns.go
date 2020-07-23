package log

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

const (
	regionUSEast1 = "us-east-1"
)

// SNSClient can be used for pushing messages to a topic which are then emailed to all users who are subscribed to that topic
type SNSClient struct {
	ctx           context.Context
	topicID       string
	client        *sns.SNS
	readMessages  <-chan string
	writeMessages chan<- string
}

const snsWaitTime = 5 * time.Minute

// NewSNSClient creates a client for creating emails VIA Amazon SNS using the [default] credentials stored in ~/.aws/credentials
func NewSNSClient(ctx context.Context, topicID string) (client *SNSClient, err error) {
	if len(topicID) > 0 {
		var sess *session.Session
		if sess, err = session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String(regionUSEast1)},
			Profile: "dev-account",
		}); err == nil {
			var messages = make(chan string, 100)
			client = &SNSClient{
				ctx:           ctx,
				topicID:       topicID,
				client:        sns.New(sess),
				readMessages:  messages,
				writeMessages: messages,
			}

			go func() {
				defer close(client.writeMessages)
				defer client.sendMessages()

				// TODO make configurable
				tic := time.Tick(snsWaitTime)

				for {

					select {
					case <-ctx.Done():
						return
					case <-tic:
						// proceed
					}

					client.sendMessages()
				}
			}()
		}
	} else {
		err = fmt.Errorf("empty topic ID")
	}

	return client, err
}

func (client *SNSClient) sendMessages() {
	messages := make([]string, 0)
	func() {
		for {
			select {
			case message := <-client.readMessages:
				messages = append(messages, message)
			default:
				return
			}
		}
	}()
	if len(messages) > 0 {
		concatenatedMessages := strings.Join(messages, "\n\n")

		params := &sns.PublishInput{
			Message:  &concatenatedMessages, // This is the message itself (can be XML / JSON / Text - anything you want)
			TopicArn: &client.topicID,       //Get this from the Topic in the AWS console.
			Subject:  aws.String("Aegis Alert"),
		}

		_, err := client.client.Publish(params)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// PushMessage pushes the message to sns
func (client *SNSClient) PushMessage(message string) {
	go func() {
		// the timeout is twice as long as the waiting period between SNS bursts
		tic := time.Tick(2 * snsWaitTime)

		select {
		case <-tic:
			return
		case client.writeMessages <- message:
		}
	}()
}
