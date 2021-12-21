FROM registry.access.redhat.com/ubi7-dev-preview/ubi-minimal:7.6

ENV OPERATOR=/usr/local/bin/balancer \
    USER_UID=1001 \
    USER_NAME=balancer

# copy operator binary
COPY build/_output/bin/balancer ${OPERATOR}
COPY build/bin /usr/local/bin

# set up user and run with it
RUN /usr/local/bin/user_setup
ENTRYPOINT ["/usr/local/bin/entrypoint"]
USER $USER_UID
