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
                        withSonarQubeEnv('SonarQube') {
                            sh('sonar-scanner')
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
