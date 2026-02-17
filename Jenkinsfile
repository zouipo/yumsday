def img

pipeline {
    agent any

    environment {
        GOCACHE="${WORKSPACE}"
    }

    stages {
        stage('Build Docker Image') {
            steps {
                script {
                    img = docker.build('zouipo/yumsday:base', '--target base .')
                }
            }
        }
        stage('Run tests') {
            steps {
                script {
                    img.inside {
                        sh('make test-cicd')
                    }
                }
            }
        }
        stage('SonarQube Analysis') {
            steps {
                script {
                    img.inside {
                        scannerHome = tool('SonarQube Scanner')
                        withSonarQubeEnv('SonarQube') {
                            sh("${scannerHome}/bin/sonar-scanner")
                        }
                    }
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
    }
}
