pipeline {
    agent { label 'worker' }

    triggers {
        gitlab(
            triggerOnPush: true,
            triggerOnMergeRequest: true,
            triggerOnAcceptedMergeRequest: true,
            branchFilterType: 'All'
        )
    }

    environment {
        DOCKER_IMAGE = 'atsova15/weather'
        DOCKER_TAG = "${GIT_COMMIT[0..6]}"
    }

    options {
        gitLabConnection('gitlab-yadro')
    }

    stages {
        stage('Lint') {
            steps {
                sh 'golangci-lint run ./...'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./tests/... -v -coverprofile=coverage.out -coverpkg=./...'
                sh 'go tool cover -func=coverage.out'
            }
        }

        stage('Build') {
            steps {
                sh 'docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .'
                sh 'docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest'
            }
        }

        stage('Deploy') {
            when {
                beforeInput true
                branch 'master'
            }
            input {
                message 'Deploy to prod? (YES if it is Friday evening)'
                ok 'Deploy'
            }
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-creds',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
                  
                    sh 'docker push ${DOCKER_IMAGE}:${DOCKER_TAG}'
                    sh 'docker push ${DOCKER_IMAGE}:latest'
                }
            }
        }
    }

    post {
        always {
            sh 'docker logout'
        }
        success {
            updateGitlabCommitStatus name: 'build', state: 'success'
        }
        failure {
            updateGitlabCommitStatus name: 'build', state: 'failed'
        }
    }
}
