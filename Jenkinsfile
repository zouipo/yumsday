pipeline {
    agent any

    environment {
        GOCACHE="${WORKSPACE}"
    }

    stages {
        stage('Run tests') {
            steps {
                script {
                    def img = docker.build('zouipo/yumsday:base', '--target base .')
                    img.inside {
                        sh("make test-cicd")
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
