#!/bin/bash
/root/cleaner/cleaner_linux /root/code/ '.*mysql.*.sql.gz$' 5  >> /root/cleaner/clean_log.txt  2>&1
/root/cleaner/cleaner_linux /root/code/ '.*mysql-.*.sql$'  5 >> /root/cleaner/clean_log.txt 2>&1