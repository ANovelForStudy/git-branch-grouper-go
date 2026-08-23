package main

const repoPath = `G:\05-Projects\01-Programming\01-Pet-Projects\06-Git-branch-grouper-plugin\99-Test-repo`

var defaultBranchNames = []string{"develop", "main", "master"}

func main() {
	validateRepoPath(repoPath)

	gitRepo := openGitRepo(repoPath)
	branchData := collectBranches(gitRepo)

	printResults(branchData)
}
