properties([disableConcurrentBuilds()])

node('linux') {

  def image_name = 'statusteam/geth_exporter'
  def commit
  def image

  stage('Git Prep') {
    checkout scm
    commit = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
  }

  stage('Build') {
    image = docker.build(image_name + ':' + commit)
  }

  stage('Publish') {
    withDockerRegistry([
      credentialsId: "dockerhub-statusteam-auto", url: ""
    ]) {
      image.push()
      image.push('latest')
    }
  }
}
