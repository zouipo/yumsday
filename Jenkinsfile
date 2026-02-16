pipeline {
    agent any

    stages {
        stage('Run tests') {
            steps {
                script {
                    def img = docker.build('zouipo/yumsday:base', '--target base')
                    img.inside {
                        make test-cicd
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
