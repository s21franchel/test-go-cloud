FROM golang:1.22.3-alpine

#WORKDIR /app
#COPY backends.go .
#RUN go build -o backends .
#EXPOSE $PORT
#CMD ["./backends"]


WORKDIR /app
COPY . .


RUN go build -o backends .

EXPOSE $PORT
CMD ["./backends"]