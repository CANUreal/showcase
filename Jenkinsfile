pipeline {
    agent any
    
    environment {
        IMAGE_NAME = "quay.io/denizdayan/showcase"
        IMAGE_TAG = "${env.BUILD_NUMBER}"
        QUAY_CREDS = credentials('quay-io-creds')
    }

    stages {
        stage('build') {
            steps {
                sh "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} ."
            }
        }

        stage('push to registry') {
            steps {
                // do this below this comment to read from shell's environment variables
                sh '''
                    echo "$QUAY_CREDS_PSW" | docker login quay.io -u "$QUAY_CREDS_USR" --password-stdin
                    docker push "$IMAGE_NAME:$IMAGE_TAG"
                '''
            }
        }
    }

    post {
        always {
            sh "docker logout quay.io || true"
            cleanWs()
        }
    }
}
