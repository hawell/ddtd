FROM golang:1.22-alpine as gobuilder

WORKDIR /app
COPY . ./

RUN go build -o ./ddtd ./cmd

FROM node:latest as nodebuilder

WORKDIR /panel
COPY ./panel ./

RUN npm install && npm run generate

FROM alpine as agent
COPY --from=gobuilder /app/ddtd /bin/ddtd
COPY --from=nodebuilder /panel/.output/public /var/www

ENV HTTP_ADDRESS=:8866
ENV HTTP_PANEL_ROOT=/var/www
ENTRYPOINT ["/bin/ddtd"]
