pipeline {
  agent any

  environment {
    PATH = "/usr/local/go/bin:${env.PATH}"
  }

  stages {
    stage('Check Requirements') {
      steps {
        echo 'Pulling... ' + env.GIT_BRANCH

        script {
          // Retrieve the last commit message using git command
          env.GIT_COMMIT_MSG = sh(script: 'git log -1 --pretty=%B ${GIT_COMMIT}', returnStdout: true).trim()
        }

        echo "${GIT_COMMIT_MSG}"
        sh 'whoami'
        echo 'Installing dependencies'
        sh 'which go'
        sh 'go version'
      }
    }

    stage('Run Test') {
      steps {
        sh 'make test'
      }
    }

    stage('Build Binary') {
      steps {
        echo 'Compiling and building Binary'
        sh 'CGO_ENABLED=0 GOOS=linux go build -o nanda-api main.go'
        sh 'chmod +x vdi-meter'


      }
    }

    stage('Deploy Server Development') {
      steps {
        // Stop any existing containers
        sh 'docker compose down'

        // Build and start the containers
        sh 'docker compose up -d --build'

        // Remove unused images
        sh 'docker image prune -f'
      }
    }

    stage('Send Email Notification') {
      steps {
        script {
          def commitMessage = env.GIT_COMMIT_MSG
          def gitBranch = env.GIT_BRANCH
          def projectName = currentBuild.projectName
          def buildNumber = currentBuild.number
          def buildStatus = currentBuild.currentResult.toString()

          // Send email notification using emailext plugin
          emailext body: "Commit Message:\n${commitMessage}",
            subject: "${projectName} - ${gitBranch} - Build #${buildNumber} - ${buildStatus}!",
            to: "nandarusfikri@gmail.com",
            from: "Jenkins <nandarusfikri@gmail.com>"
        }
      }
    }


  }
}
