#!/usr/bin/env bash

LOGGING_TARGET_DIR=/var/log/ict-flex-rss

if [ ! -d ${LOGGING_TARGET_DIR} ]; then
    mkdir -p ${LOGGING_TARGET_DIR};
    touch ${LOGGING_TARGET_DIR}/errors.log;
fi

if ! cmd=$(which ict-flex-rss); then
    echo \`ict-flex-rss\` executable not found in path. | tee -a ${LOGGING_TARGET_DIR}/errors.log;
    exit 1;
fi

if output=$(${cmd} 2>>${LOGGING_TARGET_DIR}/errors.log); then
    echo ${output} > /var/www/rss/ict-flex.rss;
fi
