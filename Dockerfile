FROM alpine
ADD collector-srv /collector-srv
ENTRYPOINT [ "/collector-srv" ]
