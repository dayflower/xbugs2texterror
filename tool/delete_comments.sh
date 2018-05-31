#!/bin/sh

curl -s -H "Authorization: token ${REVIEWDOG_GITHUB_API_TOKEN}" "https://api.github.com/repos/${CI_REPO_OWNER}/${CI_REPO_NAME}/pulls/comments" | jq '.[] | .id' | while read line
do
	echo $line
	curl -s -H "Authorization: token ${REVIEWDOG_GITHUB_API_TOKEN}" -X DELETE "https://api.github.com/repos/${CI_REPO_OWNER}/${CI_REPO_NAME}/pulls/comments/${line}"
done
