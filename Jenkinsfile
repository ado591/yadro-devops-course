@Library('jenkins-shared-lib') _

pipeline {
    agent none

    environment {
        DOCKER_IMAGE = 'atsova15/weather'
    }

    options {
        gitLabConnection('gitlab-yadro')
    }

    stages {
        stage('Lint & SAST') {
            agent any
            steps {
                checkout scm
                script {
                    env.DOCKER_TAG = GIT_COMMIT[0..6]
                }
                lintAndSast()
            }
        }

        stage('Test') {
            agent any
            steps {
                checkout scm
                runTests()
            }
        }

        stage('Build') {
            when {
                changeRequest()
            }
            agent { label 'staging' }
            steps {
                checkout scm
                sh 'docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .'
            }
        }

        stage('Push') {
            when {
                anyOf {
                    branch 'master'
                    branch 'main'
                    buildingTag()
                }
            }
            agent { label 'staging' }
            steps {
                buildAndPush(env.DOCKER_IMAGE, env.DOCKER_TAG)
            }
        }

        stage('Deploy to Stage') {
            when {
                anyOf {
                    branch 'master'
                    branch 'main'
                }
            }
            steps {
                build job: 'deploy-app',
                    parameters: [
                        string(name: 'IMAGE_TAG',   value: env.DOCKER_TAG),
                        string(name: 'ENVIRONMENT', value: 'staging')
                    ],
                    wait: true
            }
        }

        stage('Deploy to Prod') {
            when {
                buildingTag()
                tag pattern: '^v.*', comparator: 'REGEXP'
            }
            steps {
                build job: 'deploy-app',
                    parameters: [
                        string(name: 'IMAGE_TAG',   value: env.DOCKER_TAG),
                        string(name: 'ENVIRONMENT', value: 'production')
                    ],
                    wait: true
            }
        }
    }

    post {
        success {
            updateGitlabCommitStatus name: 'build', state: 'success'
        }
        failure {
            updateGitlabCommitStatus name: 'build', state: 'failed'
        }
    }
}
