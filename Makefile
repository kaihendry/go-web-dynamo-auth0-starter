STACK = zauth
VERSION = 0.2

.PHONY: build deploy validate destroy

DOMAINNAME = zero.dabase.com
ACMCERTIFICATEARN = arn:aws:acm:ap-southeast-1:407461997746:certificate/87b0fd84-fb44-4782-b7eb-d9c7f8714908

deploy:
	sam build
	SAM_CLI_TELEMETRY=0 sam deploy --resolve-s3 --stack-name $(STACK) --parameter-overrides DomainName=$(DOMAINNAME) ACMCertificateArn=$(ACMCERTIFICATEARN) --no-confirm-changeset --no-fail-on-empty-changeset --capabilities CAPABILITY_IAM --disable-rollback

build-MainFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o ${ARTIFACTS_DIR}/bootstrap

validate:
	aws cloudformation validate-template --template-body file://template.yml

destroy:
	aws cloudformation delete-stack --stack-name $(STACK)

sam-tail-logs:
	sam logs --stack-name $(STACK) --tail

clean:
	rm -rf main gin-bin

sync:
	sam sync --stack-name $(STACK) --watch

# fetch aws secretsmanager SecretString out from auth0/fhwuGmCUqUr5Sk3NxwUgPT4yJAEu1a7Z
secrets.json:
	aws secretsmanager get-secret-value --secret-id auth0/fhwuGmCUqUr5Sk3NxwUgPT4yJAEu1a7Z --query SecretString --output text > secrets.json
