pipeline {
    agent any
    
    environment {
        IMAGE_NAME = "quay.io/denizdayan/showcase"
        IMAGE_TAG = "${env.BUILD_NUMBER}"
        QUAY_CREDS = credentials('76e364a6-e3ea-4839-9c2d-7acff14df0a8')
    }

    stages {
        stage('build') {
            steps {
                container('docker') {
                    sh '''
                        docker version
                        docker build -t "$IMAGE_NAME:$IMAGE_TAG" .
                    '''
                }
            }
        }

        stage('push to registry') {
            container('docker') {
                steps {
                    // do this below this comment to read from shell's environment variables
                    sh '''
                        echo "$QUAY_CREDS_PSW" | docker login quay.io -u "$QUAY_CREDS_USR" --password-stdin
                        docker push "$IMAGE_NAME:$IMAGE_TAG"
                    '''
                }
            }
        }
    }

    post {
        always {
            container('docker') {
                sh "docker logout quay.io || true"
            }
        }
    }
}
