## Prerequisite of this project
- php 1.21 or above 
- postgresql 14.8

**Bold**

## Installation
#### Option 1: Clone from Git (With HTTPS)
- Enjoy

#### Option 2: Clone and install using docker
- Enjoy!

#### (Important) Steps For Development With Git Flow

1. start a feature with command: `git flow feature start <feature_name>`
2. start development and stage and commit your changes.
3. when feature development finished, run `git flow feature publish <feature_name>`
4. create a merge request from `<feture_name>` to `develop` branch. 
5. if your merge request got rejected by your head team member, go to step 2.
6. if merge request accepted and your feature got merged
   1. run: `git checkout develop` to change working branch to `develop`
   2. pull latest results from remote repo: `git pull origin develop`

#### (very important) note: You Should Always run step 1 when you are in your local `develop` branch. 
---

