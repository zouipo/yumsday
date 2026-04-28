pipeline {
    agent any

    environment {
        GOCACHE="${WORKSPACE}"
    }

    stages {
        stage('Run Tests') {
            when {
                not {
                    buildingTag()
                }
            }
            steps {
                script {
                    docker.build('zouipo/yumsday:base', '--target base .').inside {
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
                script {
                    withSonarQubeEnv('SonarQube') {
                        docker.image("sonarsource/sonar-scanner-cli").inside {
                            sh("sonar-scanner")
                        }
                    }
                }
            }
        }
        stage('Release Docker Image') {
            when {
                buildingTag()
            }
            steps {
                script {
                    def img = docker.build("zouipo/yumsday:${env.TAG_NAME}", '--target runtime .')
                    docker.withRegistry('', 'docker-zouipo') {
                        sh("echo ${img.id}")
                        //img.push()
                        //img.push('latest')
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
