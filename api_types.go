package main

type endpoints struct {
	Kind       string   `json:"kind"`
	ApiVersion string   `json:"apiVersion"`
	Metadata   metadata `json:"metadata"`
	Subsets    []subset `json:"subsets"`
	Message    string   `json:"message"`
}

type services struct {
	Kind  string `json:"kind"`
	Items []item `json:"items"`
}

type metadata struct {
	Name string `json:"name"`
}

type subset struct {
	Addresses []address `json:"addresses"`
	Ports     []port    `json:"ports"`
}

type item struct {
	Metadata metadata `json:"metadata"`
	Spec     spec     `json:"spec"`
}

type spec struct {
	Ports []port `json:"ports"`
}

type address struct {
	IP string `json:"ip"`
}

type port struct {
	Name string `json:"name"`
	Port int32  `json:"port"`
}

type status struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
