NAME	:= esp32-request-url-go
BIN	:= my-lambda-binary
ROLE	:= arn:aws:iam::497939524935:role/esp32-request-role
ZIP	:= function.zip
OUT	:= response.json
PAYLOAD	:= '{"key": "volume", "val": 100 }'


all: clean build update run


build:
	GOOS=linux go build -o ${BIN} main.go
	zip ${ZIP} ${BIN}

clean:
	-rm *~
	-rm *zip
	-rm ${BIN} 
	-rm ${OUT}

create:
	aws lambda create-function --function-name ${NAME} --runtime go1.x --zip-file fileb://${ZIP} --handler ${BIN} --role ${ROLE}

update:
	aws lambda update-function-code --function-name sample-event-handle --zip-file fileb://function.zip

run:
	aws lambda invoke --function-name ${NAME}  --payload ${PAYLOAD} ${OUT}
	cat ${OUT}
	echo



