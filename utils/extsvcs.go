package utils

// func sendSimpleHttpReq(method string, url string) {
// 	newReq, err := http.NewRequest(method, url, nil)
// 	if err != nil {
// 		Logger.Debug(err.Error())
// 	}
// 	resp, err := http.DefaultClient.Do(newReq)
// 	if err == nil {
// 		io.Copy(io.Discard, resp.Body)
// 		resp.Body.Close()
// 	}
// }

// func StopExtSvcs() {
// 	//Pause/sleep disable because envoy/istio-proxy does not intercept any traffic from consul-sidecar
// 	//time.Sleep(5 * time.Second)

// 	//Send to envoy
// 	//logger.Debug("Entrypoint sending quitquitquit signal to envoy-sidecar")
// 	sendSimpleHttpReq(http.MethodPost, "http://127.0.0.1:15000/quitquitquit")
// }
