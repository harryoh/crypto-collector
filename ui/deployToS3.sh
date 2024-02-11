#!/bin/bash
NODE_ENV=production yarn build
aws s3 rm s3://crypto.5004.pe.kr/ --recursive
aws s3 cp dist/ s3://crypto.5004.pe.kr/ --recursive
