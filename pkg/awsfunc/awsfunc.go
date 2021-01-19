package awsfunc

import (
	"context"
	"fmt"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type AWSsvcs struct {
	elbv2 *elasticloadbalancingv2.Client
	auto  *autoscaling.Client
}

var Logflg bool

func New() *AWSsvcs {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//fmt.Println("logflg=", Logflg)
	if Logflg {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)

	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	elbv2 := elasticloadbalancingv2.NewFromConfig(cfg)
	auto := autoscaling.NewFromConfig(cfg)

	return &AWSsvcs{
		elbv2: elbv2,
		auto:  auto,
	}
}

func (a *AWSsvcs) describeTags() {

	svc := a.auto
	param := &autoscaling.DescribeTagsInput{}

	p := autoscaling.NewDescribeTagsPaginator(svc, param)
	var i int
	asgs := []string{}
	for p.HasMorePages() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			fmt.Println(err)
		}

		for _, v := range page.Tags {
			log.WithFields(log.Fields{
				"asg": aws.ToString(v.ResourceId),
			}).Info("A walrus appears")
			asgs = append(asgs, aws.ToString(v.ResourceId))

		}
	}
}

func (a *AWSsvcs) GetAsgsFromTags(tag string) []string {
	svc := a.auto
	param := &autoscaling.DescribeTagsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("value"),
				Values: []string{tag},
			},
		},
	}

	p := autoscaling.NewDescribeTagsPaginator(svc, param)
	r := regexp.MustCompile(tag)
	var i int
	asgs := []string{}
	for p.HasMorePages() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			fmt.Println(err)
		}

		for _, v := range page.Tags {
			if r.MatchString(*v.ResourceId) {
				log.WithFields(log.Fields{
					"asg": aws.ToString(v.ResourceId),
				}).Info("A walrus appears")
				asgs = append(asgs, aws.ToString(v.ResourceId))

			}
		}
	}
	return uniq(asgs)
}

func uniq(data []string) []string {

	m := make(map[string]struct{})
	// make data to unique
	for _, v := range data {
		m[v] = struct{}{}
	}

	uniq := []string{}
	for i := range m {
		uniq = append(uniq, i)
	}
	return uniq
}
