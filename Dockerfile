FROM golang:1.21

RUN apt-get update \
    &&  DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends tzdata \
    &&  apt-get install vim -y 
    
RUN TZ=Asia/Taipei \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && dpkg-reconfigure -f noninteractive tzdata 

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build

ENTRYPOINT [ "./vitalsign-publisher"]