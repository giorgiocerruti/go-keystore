FROM golang:1.17 as build
LABEL maintainer="giorgio.cerruti@edesk.com"

RUN mkdir /application
COPY . /application
WORKDIR /application
RUN ls && mkdir ./build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./build/gks ./cmd/*.go

FROM alpine as app

COPY --from=build /application/build/gks . 

ENTRYPOINT [ "/gks" ]