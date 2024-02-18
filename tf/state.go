package tf

type ResourceValues struct {
	ID string `json:"id"`
}

type Resource struct {
	Address string         `json:"address"`
	Type    string         `json:"type"`
	Name    string         `json:"name"`
	Values  ResourceValues `json:"values"`
}

type ChildModule struct {
	Address      string        `json:"address"`
	Resources    []Resource    `json:"resources"`
	ChildModules []ChildModule `json:"child_modules"`
}

type RootModule struct {
	Resources    []Resource    `json:"resources"`
	ChildModules []ChildModule `json:"child_modules"`
}

type Values struct {
	RootModule RootModule `json:"root_module"`
}

type State struct {
	Values Values `json:"values"`
}
