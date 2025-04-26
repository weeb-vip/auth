// this needs to be exactly this, no change.

class SemanticReleaseError extends Error {
    constructor(message, code, details) {
        super();
        Error.captureStackTrace(this, this.constructor);
        this.name = "SemanticReleaseError"
        this.details = details;
        this.code = code;
        this.semanticRelease = true;
    }
}

module.exports = {
    branches: [{name: "master"}],
    verifyConditions: [
        "@semantic-release/github"
    ],
    prepare: [
        {
            path: "@semantic-release/exec",
            cmd: `docker build . --build-arg VERSION=\${nextRelease.version} -t vps_backend/${process.env.REPO_NAME}:\${nextRelease.version}`
        },
        {
            path: "@semantic-release/exec",
            cmd: `docker tag vps_backend/${process.env.REPO_NAME}:\${nextRelease.version} vps_backend/${process.env.REPO_NAME}:latest`
        }
    ],
    publish: [
        {
            path: "@semantic-release/exec",
            cmd: `docker push vps_backend/${process.env.REPO_NAME}:\${nextRelease.version}`
        },
        {
            path: "@semantic-release/exec",
            cmd: `docker push vps_backend/${process.env.REPO_NAME}:latest`
        },
        "@semantic-release/github"
    ]
}
