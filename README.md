# Welcome to your CDK Go project!

This is a blank project for CDK development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests

https://static.frontendmasters.com/assets/courses/2024-04-23-go-aws/go-aws-slides.pdf

##
```
cdk init app --language go
go get
go get github.com/aws/aws-lambda-go/lambda
in lambda dir:
GOOS=linux GOARCH=amd64 go build -o bootstrap
zip function.zip bootstrap
```
GOOS=linux: This environment variable tells the Go compiler to build the program for Linux operating systems. GOOS stands for "Go Operating System"
GOARCH=amd64: This specifies the target processor architecture as AMD64 (also known as x86-64), which is the 64-bit version of the x86 instruction set used by most modern processors from both Intel and AMD
go build: This is the Go compiler command to build an executable from your source code
-o bootstrap: The -o flag specifies the output filename. In this case, the compiled executable will be named "bootstrap"


cdk interaction 
```
cdk bootstrap aws://ACCOUNT-NUMBER/REGION
cdk diff
cdk deploy
```