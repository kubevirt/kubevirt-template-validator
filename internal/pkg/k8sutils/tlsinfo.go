package k8sutils

type TLSInfo struct {
	CertFilePath   string
	KeyFilePath    string
}

func (ti *TLSInfo) IsEnabled() bool {
	return ti.CertFilePath != "" && ti.KeyFilePath != ""
}
