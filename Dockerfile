# Use an official Python runtime as a parent image
FROM golang

# Set the working directory to /app
WORKDIR /go/src/github.com/katenicoletti/buoy-api

# Copy the current directory contents into the container at /app
ADD . /go/src/github.com/katenicoletti/buoy-api

# Install any needed packages specified in requirements.txt
RUN cd $GOPATH/src/github.com/katenicoletti/buoy-api
RUN go install

# Make port 80 available to the world outside this container
EXPOSE 80

# Define environment variable
ENV NAME buoy-api
ENV PGPORT 5432
ENV PGUSER postgres
ENV PGHOST docker.for.mac.localhost

# Run main.go when the container launches
CMD ["buoy-api"]
