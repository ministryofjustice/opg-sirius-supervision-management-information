# 1. Create a new MI Back End

Date: 2026-01-13

## Status

Accepted

## Context

We have created the management information front end, however the decision was made to add a service between S3 and the front end to upload files to S3 without making AWS credentials accessible to the user.

## Decision

A new back end will be created in this repo, similarly to the structure of finance hub and finance admin. This will communicate with S3 to perform uploads using Go's AWS SDK. 

## Consequences

This ensures we can interact with S3 whilst keeping our AWS credentials away from the client. 

