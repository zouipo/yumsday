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
            when {
                not {
                    buildingTag()
                }
            }
            steps {
                script {
                    img.inside {
                        sh('make test-cicd')
                    }
                }
            }
        }
        stage('SonarQube Analysis') {
            when {
                branch 'main'
            }
            steps {
                withSonarQubeEnv('SonarQube') {
                    sh("${tool('SonarQube Scanner')}/bin/sonar-scanner")
                }
            }
        }
        stage('Image build and push') {
            when {
                buildingTag()
            }
            steps {
                script {
                    img = docker.build("zouipo/yumsday:${env.TAG_NAME}", '--target runtime .')
                    docker.withRegistry('', 'docker-zouipo') {
                        img.push()
                        img.push('latest')
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
