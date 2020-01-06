package test

import (
	"crypto/tls"
	"fmt"
	"net"
	"testing"
	"time"

	// Testify
	"github.com/stretchr/testify/suite"

	// Terratest
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const instanceText = "Hello, Terratest!"

type TerraformTestSuite struct {
	suite.Suite
	terraformOptions *terraform.Options
}

func TestTerraformTestSuite(t *testing.T) {
	suite.Run(t, new(TerraformTestSuite))
}

func (suite *TerraformTestSuite) SetupSuite() {

	suite.terraformOptions = &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"response_text": instanceText,
		},
	}

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(suite.T(), suite.terraformOptions)
}

func (suite *TerraformTestSuite) TearDownSuite() {
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	terraform.Destroy(suite.T(), suite.terraformOptions)
}

func (suite *TerraformTestSuite) TestTerraformHttpExample() {

	// Run `terraform output` to get the value of an output variable
	url := terraform.Output(suite.T(), suite.terraformOptions, "url")

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second

	// Verify that we get back a 200 OK with the expected instanceText
	http_helper.HttpGetWithRetry(
		suite.T(),
		url,
		&tlsConfig,
		200,
		instanceText,
		maxRetries,
		timeBetweenRetries)
}

func (suite *TerraformTestSuite) TestDial() {
	hostname := terraform.Output(suite.T(), suite.terraformOptions, "hostname")
	address := net.JoinHostPort(hostname, "8080") // Whatever port makes sense

	timeOut := time.Duration(5) * time.Second
	maxRetries := 12
	timeBetweenRetries := 5 * time.Second

	retry.DoWithRetry(
		suite.T(),
		fmt.Sprintf("Verify %s is open", address),
		maxRetries, timeBetweenRetries,
		func() (string, error) {
			conn, err := net.DialTimeout("tcp", address, timeOut)

			if err != nil {
				return "", err
			}
			if conn == nil {
				return "", fmt.Errorf("conn is nil")
			}
			defer conn.Close()
			return "", nil
		})

}
