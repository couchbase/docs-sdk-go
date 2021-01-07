#!/bin/env groovy

// Jenkins job configuration
// -------------------------
// Category: Multibranch Pipeline
// Pipeline name: docs-sdk-go
// Source Code Management: Github
// Branch Sources: Multiple branches (All branches)
// Discover pull requests from forks: Current pull request version
// All branches get the same properties
// Trust: Contributors
// Repository URL: https://github.com/couchbase/docs-sdk-go
// Credentials: - cb-docs-robot -
// Build Configuration:
// Mode: by Jenkinsfile
// Script Path: Jenkinsfile

pipeline {
  agent {
    dockerfile {
        filename 'Dockerfile.jenkins'
    }
  }

  stages {
    stage('Build Code Examples') {
      steps {
        sh 'make build-all-examples'
      }
    }
  }
}
