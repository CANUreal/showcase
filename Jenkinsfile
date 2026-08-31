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
                        for i in $(seq 1 30); do
                            docker info >/dev/null 2>&1 && break
                            sleep 1
                        done

                        docker info
                        docker build -t "$IMAGE_NAME:$IMAGE_TAG" .
                    ''' 
                }
            }
        }

        stage('Push to Registry') {
            steps {
                container('docker') {
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
