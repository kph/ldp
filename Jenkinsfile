#!groovy

import groovy.transform.Field

@Field String email_to = 'sw@platinasystems.com'
@Field String email_from = 'jenkins-bot@platinasystems.com'
@Field String email_reply_to = 'no-reply@platinasystems.com'

pipeline {
    agent any
    environment {
	GOPATH = "$WORKSPACE/go-pkg"
	HOME = "$WORKSPACE"
    }
    stages {
	stage('Checkout') {
	    steps {
		echo "Running build #${env.BUILD_ID} on ${env.JENKINS_URL} GOPATH ${GOPATH}"
		dir('ldp') {
		    git([
			url: 'https://github.com/platinasystems/ldp.git',
			branch: 'master'
		    ])
		}
	    }
	}
	stage('Build') {
	    steps {
		dir('ldp') {
		    echo "Building ldp..."
		    sh 'set +x; export PATH=/usr/local/go/bin:${PATH}; go build -x; go test -x'
		}
	    }
	}
    }

    post {
	success {
	    mail body: "ldp build ok: ${env.BUILD_URL}",
		from: email_from,
		replyTo: email_reply_to,
		subject: 'ldp build ok',
		to: email_to
	}
	failure {
	    cleanWs()
	    mail body: "ldp build error: ${env.BUILD_URL}",
		from: email_from,
		replyTo: email_reply_to,
		subject: 'ldp BUILD FAILED',
		to: email_to
	}
    }
}
