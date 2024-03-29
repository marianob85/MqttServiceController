properties(
	[
		buildDiscarder(logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '5')),
	]
)
pipeline
{
	agent any
	options {
		skipDefaultCheckout true
	}
	environment {
		GITHUB_TOKEN = credentials('marianob85-github-jenkins')
		NEXUS_CREDS = credentials('Nexus-Mariano-HTPC')
	}
	
	stages
	{
		stage('Build package') 
		{
			agent{ label "linux/u18.04/go:1.20.1" }
			steps
			{
				checkout scm
				script {
					env.GITHUB_REPO = sh(script: 'basename $(git remote get-url origin) .git', returnStdout: true).trim()
				}
				sh '''
					make package
				'''
				archiveArtifacts artifacts: 'build/**', onlyIfSuccessful: true,  fingerprint: true
				stash includes: 'build/dist/**', name: 'dist'
			}
		}

		stage('Release') {
			when {
				buildingTag()
			}
			agent{ label "linux/u18.04/go:1.20.1" }
			steps {
				unstash 'dist'
				sh '''
					export GOPATH=${PWD}
					go install github.com/github-release/github-release@v0.10.0
					bin/github-release release --user marianob85 --repo ${GITHUB_REPO} --tag ${TAG_NAME} --name ${TAG_NAME}
					sleep 2m
					for filename in build/dist/*; do
						[ -e "$filename" ] || continue
						basefilename=$(basename "$filename")
						bin/github-release upload --user marianob85 --repo ${GITHUB_REPO} --tag ${TAG_NAME} --name ${basefilename} --file ${filename}
					done
				'''
			}
		}
		
		stage('Nexus upload') {
			agent{ label "linux/u18.04/base" }
			when {
				buildingTag()
			}
			steps {
				unstash 'dist'
				sh '''
					for f in  build/dist/*.deb; do
						[ -e "$f" ] || continue
						STATUS=$(curl -s -o /dev/null -w '%{http_code}' --insecure -u ${NEXUS_CREDS_USR}:${NEXUS_CREDS_PSW} -H "Content-Type: multipart/form-nedata" --data-binary @$f ${NEXUS_SERVER}/repository/ubuntu/)
						if [ $STATUS -ne 201 ]; then
							exit $STATUS
						fi
					done		
				'''
			}
		}

	}
	post { 
        changed { 
            emailext body: "Please go to ${env.BUILD_URL}", to: '${DEFAULT_RECIPIENTS}', subject: "Job ${env.JOB_NAME} (${env.BUILD_NUMBER}) ${currentBuild.currentResult}".replaceAll("%2F", "/")
        }
    }
}