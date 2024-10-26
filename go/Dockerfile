FROM golang:1.23-alpine 

# installing bash and ffmpeg 
RUN apk add --no-cache bash ffmpeg

WORKDIR /app

COPY . .

CMD ["bash"]

#  production
# COPY go.mod go.sum ./

# RUN go mod download