load('ext://configmap', 'configmap_create')
configmap_create('azure-resources-export-config', namespace='default', from_file='config.json')
docker_build('localhost:5000/azure-pricing-exporter', context='.')
k8s_yaml('k8s.yaml')
k8s_resource(workload='azure-pricing-exporter', port_forwards='8080:8080')