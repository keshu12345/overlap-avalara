pipeline {
  agent any

  environment {
     PATH = "/usr/local/bin:/usr/local/go/bin:${env.PATH}"
  }

  stages {
    stage('Code') {
      steps {
        echo "→ Cloning code"
        git url: 'https://github.com/keshu12345/overlap-avalara', branch: 'main'
      }
    }


    stage('Docker Health Check') {
      steps {
        echo "→ Checking Docker status"
        script {
          def dockerRunning = false
          try {
            sh 'docker info > /dev/null 2>&1'
            dockerRunning = true
            echo "docker daemon is running"
          } catch (Exception e) {
            echo "Docker daemon is not running"
            echo "Please start Docker Desktop and wait for it to fully initialize"
            echo "You can verify Docker is ready by running: docker info"
            error "Docker daemon is not accessible. Pipeline stopped."
          }
        }
      }
    }
    stage('Build') {
      steps {
        echo "→ Building Docker image"
        sh 'which docker'
        sh 'docker version'
        
        // Login to Docker Hub to avoid 429 rate limit errors
        withCredentials([usernamePassword(credentialsId: 'credID', 
                                          passwordVariable: 'DOCKER_PASSWORD', 
                                          usernameVariable: 'DOCKER_USERNAME')]) {
          sh '''
            echo "→ Logging into Docker Hub"
            echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
            
            echo "→ Pulling base images"
            docker pull golang:1.24.2-alpine
            docker pull alpine:latest
            
            echo "→ Building Docker image"
            docker build -t overlap-avalara:latest .
          '''
        }
      }
    }

  stage('Push Image') {
      steps {
        echo "→ Pushing Docker image to Docker Hub"
        withCredentials([usernamePassword(credentialsId: 'credID', 
                                          passwordVariable: 'DOCKER_PASSWORD', 
                                          usernameVariable: 'DOCKER_USERNAME')]) {
          sh '''
            echo "→ Logging into Docker Hub for push"
            echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
            
            echo "→ Tagging image for Docker Hub"
            docker tag overlap-avalara:latest $DOCKER_USERNAME/overlap-avalara:latest
            
            echo "→ Pushing images to Docker Hub"
            docker image tag overlap-avalara overlap-avalara:latest
            docker push $DOCKER_USERNAME/overlap-avalara:latest
            
            echo " successfully pushed images:"
            echo "  - $DOCKER_USERNAME/overlap-avalara:latest"
            
            echo "→ Listing pushed images"
            docker images | grep overlap-avalara
          '''
        }
      }
    }
    
        stage('Test') {
          steps {
            echo "→ Running Go unit tests with coverage"
            // Verify Go is available
            sh 'go version'
            // Download modules
            sh 'go mod download'
            // Run all tests, output verbose logs and record coverage
            sh 'go test -v -coverprofile=coverage.out ./...'
          }
          post {
            always {
              // Archive the coverage report so you can download it from the build
              archiveArtifacts artifacts: 'coverage.out', fingerprint: true
            }
          }
        }

    stage('Deploy') {
      steps {
        echo "→ Deploying container"
        sh 'docker compose down && docker compose up -d --remove-orphans'
        echo "docker compose running and up in port :8081"
      }
    }
  }
}
