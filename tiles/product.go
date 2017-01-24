package tiles

type Product struct {
	Metadata Metadata
	Releases []Release
}

func (p Product) UnusedReleaseJobs() []Release {
	var releases []Release

	jobTemplateUsageCounts := map[string]map[string]int{}
	for _, release := range p.Releases {
		jobTemplateUsageCounts[release.Name] = map[string]int{}
		for _, job := range release.Jobs {
			jobTemplateUsageCounts[release.Name][job.Name] = 0
		}
	}

	for _, job := range p.Metadata.Jobs {
		for _, template := range job.Templates {
			jobTemplateUsageCounts[template.Release][template.Name]++
		}
	}

	for _, release := range jobTemplateUsageCounts {
		for jobName, count := range release {
			if count != 0 {
				delete(release, jobName)
			}
		}
	}

	for releaseName, releaseJobs := range jobTemplateUsageCounts {
		if len(releaseJobs) > 0 {
			release := Release{Name: releaseName}

			for jobName, _ := range releaseJobs {
				release.Jobs = append(release.Jobs, ReleaseJob{Name: jobName})
			}

			releases = append(releases, release)
		}
	}

	return releases
}

func (p Product) UnusedReleasePackages() []Release {
	var releases []Release

	packageUsageCounts := map[string]map[string]int{}
	for _, release := range p.Releases {
		packageUsageCounts[release.Name] = map[string]int{}
		for _, pkg := range append(release.Packages, release.CompiledPackages...) {
			packageUsageCounts[release.Name][pkg.Name] = 0
		}
	}

	for _, metadataJob := range p.Metadata.Jobs {
		for _, metadataJobTemplate := range metadataJob.Templates {
			for _, release := range p.Releases {
				if metadataJobTemplate.Release == release.Name {
					for _, releaseJob := range release.Jobs {
						for _, pkg := range releaseJob.AllPackages(release) {
							packageUsageCounts[release.Name][pkg]++
						}
					}
				}
			}
		}
	}

	for _, releasePackages := range packageUsageCounts {
		for packageName, count := range releasePackages {
			if count != 0 {
				delete(releasePackages, packageName)
			}
		}
	}

	for releaseName, releasePackage := range packageUsageCounts {
		if len(releasePackage) > 0 {
			release := Release{Name: releaseName}

			for packageName, _ := range releasePackage {
				release.Packages = append(release.Packages, ReleasePackage{Name: packageName})
			}

			releases = append(releases, release)
		}
	}

	return releases
}
