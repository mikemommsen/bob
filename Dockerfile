FROM golang:alpine as builder
RUN apk update && apk add --no-cache git

RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get github.com/mikemommsen/cuesheet
RUN go build -o main .
# this is the second image build
FROM alpine
ARG manifest

RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
COPY --from=builder /build/*.json /app/
COPY --from=builder /build/web/* /app/web/
WORKDIR /app
EXPOSE 5000
CMD ["./main"]
