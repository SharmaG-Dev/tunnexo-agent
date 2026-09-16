package agent

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (a *Agent) handleHTTPRequest(
	payload map[string]any,
) {

	requestData :=
		parseHTTPRequest(payload)

	if requestData.RequestID == "" {

		fmt.Println(
			"✗ Request ignored: missing requestId",
		)

		return
	}

	if requestData.Method == "" {
		requestData.Method = http.MethodGet
	}

	if requestData.Path == "" {
		requestData.Path = "/"
	}

	targetURL :=
		buildTargetURL(
			a.config.LocalTarget,
			requestData.Path,
		)

	fmt.Println("→ Forwarding request")

	var requestBody io.Reader

	if requestData.Method != http.MethodGet &&
		requestData.Method != http.MethodHead &&
		requestData.Body != "" {

		requestBody =
			bytes.NewBufferString(
				requestData.Body,
			)
	}

	localRequest, err :=
		http.NewRequest(
			requestData.Method,
			targetURL,
			requestBody,
		)

	if err != nil {

		a.sendBadGateway(
			requestData.RequestID,
			err,
		)

		return
	}

	copyRequestHeaders(
		localRequest,
		requestData.Headers,
	)

	response, err :=
		a.httpClient.Do(
			localRequest,
		)

	if err != nil {

		fmt.Println("✗ Local target request failed")

		a.sendBadGateway(
			requestData.RequestID,
			err,
		)

		return
	}

	defer response.Body.Close()

	responseBody, err :=
		io.ReadAll(response.Body)

	if err != nil {

		a.sendBadGateway(
			requestData.RequestID,
			err,
		)

		return
	}

	responseHeaders :=
		extractResponseHeaders(
			response,
		)

	contentType :=
		response.Header.Get(
			"Content-Type",
		)

	binary :=
		isBinaryContentType(
			contentType,
		)

	var body string

	if binary {

		body =
			base64.StdEncoding.EncodeToString(
				responseBody,
			)

	} else {

		body =
			string(responseBody)
	}

	tunnexoResponse :=
		HTTPResponse{
			RequestID: requestData.RequestID,

			StatusCode: response.StatusCode,

			Headers: responseHeaders,

			Body: body,

			IsBase64Encoded: binary,
		}

	if err :=
		a.socket.Emit(
			"http:response",
			tunnexoResponse,
		); err != nil {

		fmt.Println("✗ Failed sending response")

		return
	}

	fmt.Printf(
		"← %d %s\n",
		response.StatusCode,
		http.StatusText(response.StatusCode),
	)
}

func parseHTTPRequest(
	payload map[string]any,
) HTTPRequest {

	request :=
		HTTPRequest{
			Headers: map[string]string{},
		}

	if value,
		ok := payload["requestId"].(string); ok {

		request.RequestID = value
	}

	if value,
		ok := payload["subdomain"].(string); ok {

		request.Subdomain = value
	}

	if value,
		ok := payload["method"].(string); ok {

		request.Method =
			strings.ToUpper(value)
	}

	if value,
		ok := payload["path"].(string); ok {

		request.Path = value
	}

	if value,
		ok := payload["body"].(string); ok {

		request.Body = value
	}

	if headers,
		ok :=
		payload["headers"].(map[string]any); ok {

		for key, value := range headers {

			request.Headers[key] =
				fmt.Sprint(value)
		}
	}

	return request
}

func buildTargetURL(
	target string,
	path string,
) string {

	target =
		strings.TrimRight(
			target,
			"/",
		)

	if path == "" {
		path = "/"
	}

	if !strings.HasPrefix(
		path,
		"/",
	) {

		path = "/" + path
	}

	return target + path
}

func copyRequestHeaders(
	request *http.Request,
	headers map[string]string,
) {

	for key, value := range headers {

		if shouldSkipRequestHeader(key) {
			continue
		}

		request.Header.Set(
			key,
			value,
		)
	}
}

func shouldSkipRequestHeader(
	header string,
) bool {

	header =
		strings.ToLower(
			strings.TrimSpace(header),
		)

	switch header {

	case "host":
		return true

	case "connection":
		return true

	case "proxy-connection":
		return true

	case "keep-alive":
		return true

	case "proxy-authenticate":
		return true

	case "proxy-authorization":
		return true

	case "te":
		return true

	case "trailers":
		return true

	case "transfer-encoding":
		return true

	case "upgrade":
		return true

	default:
		return false
	}
}

func extractResponseHeaders(
	response *http.Response,
) map[string]string {

	headers :=
		make(map[string]string)

	for key, values := range response.Header {

		if shouldSkipResponseHeader(key) {
			continue
		}

		headers[key] =
			strings.Join(
				values,
				", ",
			)
	}

	return headers
}

func shouldSkipResponseHeader(
	header string,
) bool {

	header =
		strings.ToLower(
			strings.TrimSpace(header),
		)

	switch header {

	case "connection":
		return true

	case "proxy-connection":
		return true

	case "keep-alive":
		return true

	case "transfer-encoding":
		return true

	case "upgrade":
		return true

	default:
		return false
	}
}

func isBinaryContentType(
	contentType string,
) bool {

	contentType =
		strings.ToLower(
			contentType,
		)

	if strings.HasPrefix(
		contentType,
		"image/",
	) {
		return true
	}

	if strings.HasPrefix(
		contentType,
		"audio/",
	) {
		return true
	}

	if strings.HasPrefix(
		contentType,
		"video/",
	) {
		return true
	}

	if strings.Contains(
		contentType,
		"application/octet-stream",
	) {
		return true
	}

	if strings.Contains(
		contentType,
		"application/pdf",
	) {
		return true
	}

	if strings.Contains(
		contentType,
		"application/zip",
	) {
		return true
	}

	return false
}

func (a *Agent) sendBadGateway(
	requestID string,
	proxyErr error,
) {

	errorBody,
		_ :=
		json.Marshal(
			map[string]string{
				"error": "Local Target Unreachable",

				"message": proxyErr.Error(),
			},
		)

	response :=
		HTTPResponse{
			RequestID: requestID,

			StatusCode: http.StatusBadGateway,

			Headers: map[string]string{
				"content-type": "application/json",
			},

			Body: string(errorBody),

			IsBase64Encoded: false,
		}

	if err :=
		a.socket.Emit(
			"http:response",
			response,
		); err != nil {

		fmt.Println(
			"Failed to send 502 response",
		)
	}
}
