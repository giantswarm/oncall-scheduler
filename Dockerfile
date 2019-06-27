FROM alpine:3.8

 RUN apk update && apk --no-cache add ca-certificates && \
  update-ca-certificates

 ADD ./oncall-scheduler /usr/local/bin/oncall-scheduler
ENTRYPOINT ["/usr/local/bin/oncall-scheduler"]