package project

var (
	description = "Tool to schedule oncall shifts."
	gitSHA      = "n/a"
	name        = "oncall-scheduler"
	source      = "https://github.com/giantswarm/oncall-scheduler"
	version     = "n/a"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
