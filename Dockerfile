FROM registry.access.redhat.com/ubi8/ubi
MAINTAINER Pavel Macík <pavel.macik@gmail.com>
ADD aws-rds /aws-rds
ENTRYPOINT ["/aws-rds"]
