FROM golang
LABEL authors="abdul"

ENTRYPOINT ["top", "-b"]