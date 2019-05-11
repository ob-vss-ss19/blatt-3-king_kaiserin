pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'cd messages && make regenerate'
                sh 'cd tree && go build node.go'
                sh 'cd treeservice && go build main.go token.go'
                sh 'cd treecli && go build main.go'
            }
        }
        stage('Test') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'cd tree && go get -v -d -t ./...'
                sh 'go get github.com/t-yuki/gocover-cobertura' // install Code Coverage Tool
                sh 'cd tree && go test -v -coverprofile=cover.out' // save coverage info to file
                sh 'gocover-cobertura < tree/cover.out > coverage.xml' // transform coverage info to jenkins readable format
                publishCoverage adapters: [coberturaAdapter('coverage.xml')] // publish report on Jenkins
            }
        }
        stage('Lint') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }   
            steps {
                sh 'golangci-lint run --deadline 20m --enable-all'
            }
        }
        stage('Build Docker Image') {
            agent any
            steps {
                sh "docker-build-and-push -b ${BRANCH_NAME} -s treeservice -f treeservice.dockerfile"
                sh "docker-build-and-push -b ${BRANCH_NAME} -s treecli -f treecli.dockerfile"
            }
        }
    }
}
