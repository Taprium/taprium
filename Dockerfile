FROM golang:alpine
WORKDIR /src
COPY main.go go.mod go.sum  ./
COPY migrations ./migrations
# COPY helpers ./helpers
COPY helper ./helper
# COPY scheduledjob ./scheduledjob
RUN go build -o /bin/pocketbase


FROM alpine:latest

# ARG PB_VERSION=0.21.3

WORKDIR /pb
COPY --from=0 /bin/pocketbase ./
# COPY *.sh ./
# COPY crontab.txt ./
# RUN /usr/bin/crontab /pb/crontab.txt
# RUN touch /var/log/pb-job.log

# RUN apk add --no-cache \
#     unzip \
#     ca-certificates

# # download and unzip PocketBase
# ADD https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_amd64.zip /tmp/pb.zip
# RUN unzip /tmp/pb.zip -d /pb/

# uncomment to copy the local pb_migrations dir into the image
# COPY ./pb_migrations /pb/pb_migrations

# uncomment to copy the local pb_hooks dir into the image
# COPY ./pb_hooks /pb/pb_hooks
# ENV PB_ENABLE_WEB_UI=0

EXPOSE 8080

# start PocketBase
CMD ["/pb/pocketbase", "serve", "--http=0.0.0.0:8080"]
# ENTRYPOINT [ "./entrypoint.sh" ]
# CMD [ "sh","entrypoint.sh" ]