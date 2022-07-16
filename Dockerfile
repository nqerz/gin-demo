FROM golang:1.18-stretch As builder

# ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ARG PROJECT_NAME="gin-demo"
ARG PROJECT_PATH="/go/src"

WORKDIR $PROJECT_PATH

RUN sed -i 's/http/https/g' /etc/apt/sources.list && \
    apt update && \
    apt install --no-install-recommends -y git tzdata ca-certificates && \
    rm /etc/localtime && \
    ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY ./ .

RUN go build -v -i -o ${PROJECT_NAME} ./

FROM	scratch
LABEL	maintainer="eric"

ENV BIN_PATH "/usr/local/bin"
ENV TZ Asia/Shanghai

ARG PROJECT_NAME="gin-demo"
ARG PROJECT_PATH="/go/src"

COPY --from=builder $PROJECT_PATH/$PROJECT_NAME $BIN_PATH/$PROJECT_NAME
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /etc/timezone /etc/timezone
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

EXPOSE      8080
# USER        nobody
ENTRYPOINT  [ "/usr/local/bin/gin-demo" ]