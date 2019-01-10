package diffs

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/equality"
	prowconfig "k8s.io/test-infra/prow/config"
)

// GetChangedPresubmits returns a mapping of repo to presubmits to execute.
func GetChangedPresubmits(configPath, jobConfigPath string) (map[string][]prowconfig.Presubmit, error) {
	ret := make(map[string][]prowconfig.Presubmit)

	prowMasterConfig, err := prowconfig.Load(configPath, jobConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load master's Prow config: %v", err)
	}

	prowPRConfig, err := prowconfig.Load("./cluster/ci/config/prow/config.yaml", "./ci-operator/jobs/")
	if err != nil {
		return nil, fmt.Errorf("failed to load PR's Prow config: %v", err)
	}

	masterJobs := getJobsByRepoAndName(prowMasterConfig.JobConfig.Presubmits)
	for repo, jobs := range prowPRConfig.JobConfig.Presubmits {
		for _, job := range jobs {
			if !equality.Semantic.DeepEqual(masterJobs[repo][job.Name].Spec, job.Spec) {
				ret[repo] = append(ret[repo], job)
			}
		}
	}
	return ret, nil
}

// To compare two maps of slices, instead of iterating through the slice
// and compare the same key and index of the other map of slices,
// we convert them as `repo-> jobName-> Presubmit` to be able to
// access any specific elements of the Presubmits without the need to iterate in slices.
func getJobsByRepoAndName(presubmits map[string][]prowconfig.Presubmit) map[string]map[string]prowconfig.Presubmit {
	jobsByRepo := make(map[string]map[string]prowconfig.Presubmit)

	for repo, preSubmitList := range presubmits {
		pm := make(map[string]prowconfig.Presubmit)
		for _, p := range preSubmitList {
			pm[p.Name] = p
		}
		jobsByRepo[repo] = pm
	}
	return jobsByRepo
}