# External DNS Connector source example implementation
This implementation is basically https://github.com/kubernetes-sigs/external-dns/blob/master/source/connector_test.go wrapped in an application for helm deployment. 

The port numbers are hardcoded in both the application and the helm chart, which should be okay because this is just a proof of concept.
