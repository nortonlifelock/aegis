package aqua

const (
	getRepositories     = "/api/v2/repositories?include_totals=true&order_by=name"                                                                                              // page=1&pagesize=50&
	getVulnerabilities  = "/api/v2/risks/vulnerabilities?include_vpatch_info=true&show_negligible=true&hide_base_image=false&image_name=$IMAGENAME&registry_name=$REGISTRYNAME" // &page=1&pagesize=50
	getImages           = "/api/v2/images?registry=$REGISTRYNAME&include_totals=true&order_by=name&repository=$REPOSITORYNAME"                                                  // page=1 pagesize=10
	postStartImageScan  = "/api/v1/scanner/registry/$REGISTRYNAME/image/$IMAGENAME/scan"
	postCreateException = "/api/v2/risks/acknowledge"
	getImageScanStatus  = "/scanner/registry/$REGISTRYNAME/image/$IMAGENAME/status"
	getExceptions       = "/api/v2/risks/vulnerabilities?include_vpatch_info=true&show_negligible=true&hide_base_image=false&acknowledge_status=true"

	getContainersFromNameSpace = "/api/v1/containers?groupby=containers&namespace=$NAMESPACE" // &status=running&page=1&pagesize=50
)
