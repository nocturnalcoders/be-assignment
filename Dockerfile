# Use an official MongoDB image as the base image
FROM mongo:latest

# Set the working directory in the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Expose the port used by MongoDB
EXPOSE 27017

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    wget \
    gnupg2 \
    lsb-release

# Install Go
RUN wget -q https://golang.org/dl/go1.17.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz && \
    rm go1.17.5.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOBIN="/go/bin"

# Install Gin-Gonic and Swagger
RUN go get -u github.com/swaggo/swag/cmd/swag \
    github.com/swaggo/gin-swagger \
    github.com/swaggo/files \
    github.com/gin-gonic/gin \
    github.com/supertokens/supertokens-golang/recipe/session \
    github.com/supertokens/supertokens-golang/supertokens \
    go.mongodb.org/mongo-driver/mongo \
    golang.org/x/crypto/bcrypt

# Generate Swagger documentation
RUN swag init

# Copy the Supabase config file (if available)
COPY supabase_config.json .

# Set the entry point
CMD ["go", "run", "main.go"]
