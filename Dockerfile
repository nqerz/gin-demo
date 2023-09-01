FROM golang:1.21-bookworm As build-stage

# ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Create appuser.
ENV USER=appuser
ENV UID=10001 
ENV USERGROUP=appgroup
ENV GID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN addgroup \
    --gid "${GID}" \
    "${USERGROUP}" && \
    adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \  
    --gid "${GID}" \
    "${USER}"


WORKDIR /app

RUN sed -i 's/http/https/g' /etc/apt/sources.list.d/debian.sources && \
    apt update && \
    apt install --no-install-recommends -y git tzdata ca-certificates && \
    rm /etc/localtime && \
    ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY ./ .

RUN go build -v -o /gin-demo

FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /gin-demo /gin-demo
COPY --from=build-stage /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build-stage /etc/localtime /etc/localtime
COPY --from=build-stage /etc/timezone /etc/timezone
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /etc/passwd /etc/passwd
COPY --from=build-stage /etc/group /etc/group

EXPOSE 8080

USER 10001:10001

ENTRYPOINT  [ "/gin-demo" ]