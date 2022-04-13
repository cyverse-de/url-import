#!groovy
milestone 0
node('docker') {
    try {
        slackJobDescription = "job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})"
        def commitHash = checkout(scm).GIT_COMMIT
        def descriptive_version = sh(returnStdout: true, script: 'git describe --long --tags --dirty --always').trim()
        docker.withRegistry('https://harbor.cyverse.org', 'jenkins-harbor-credentials') {
            def dockerImage
	    stage('Build') {
                milestone 50
                dockerImage = docker.build("harbor.cyverse.org/de/url-import:${env.BUILD_TAG}", "--build-arg git_commit=${commitHash} --build-arg descriptive_version=${descriptiveVersion}")
                milestone 51
                dockerImage.push()
            }
            stage('Test') {
               try {
                 sh "docker run --rm --entrypoint 'sh' ${dockerImage.imageName()} -c \"go test -v github.com/cyverse-de/url-import | tee /dev/stderr | go-junit-report\" > test-results.xml"
               } finally {
                   junit 'test-results.xml'

                   sh "docker run --rm -v \$(pwd):/build -w /build alpine rm -r test-results.xml"
               }
            }
            stage('Docker Push'){
                milestone 100
                dockerImage.push("${env.BRANCH_NAME}")
                // Retag to 'qa' if this is master/main (keep both so when it switches this keeps working)
                if ( "${env.BRANCH_NAME}" == "master" || "${env.BRANCH_NAME}" == "main" ) {
                    dockerImage.push("qa")
                }
                milestone 101
            }
        }
   } catch (InterruptedException e) {
       currentBuild.result = "ABORTED"
       slackSend color: 'warning', message: "ABORTED: ${slackJobDescription}"
       throw e
   } catch (e) {
       currentBuild.result = "FAILED"
       sh "echo ${e}"
       slackSend color: 'danger', message: "FAILED: ${slackJobDescription}"
       throw e
   }
}
