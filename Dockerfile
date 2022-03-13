# build backend
FROM golang:1.17-alpine

RUN mkdir /app
COPY . /app
WORKDIR /app/backend

RUN go build -o server .


#build front end
FROM node:16

RUN mkdir /app 
COPY . /app
WORKDIR /app/webapp
RUN npm install
RUN npm run build 


# start server
FROM alpine:latest
RUN mkdir /app 
WORKDIR /app
RUN mkdir ./static/
COPY --from=0 /app/backend/server ./
COPY --from=1 /app/webapp/build ./static
EXPOSE 80
CMD ["./server"]