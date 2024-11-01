package wabot

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/pkg/errors"
)

type WabotApiClient struct {
    ClientID       string
    ClientSecret   string
    AccessToken    string
    RefreshToken   string
    TokenExpiresAt time.Time
    ApiBaseURL     string
    HttpClient     *http.Client
}

func NewWabotApiClient(clientID, clientSecret string) *WabotApiClient {
    return &WabotApiClient{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        ApiBaseURL:   "https://api.wabot.shop/v1",
        HttpClient:   &http.Client{Timeout: 30 * time.Second},
    }
}

// Authenticate and obtain access token
func (client *WabotApiClient) Authenticate() error {
    url := fmt.Sprintf("%s/authenticate", client.ApiBaseURL)

    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return errors.Wrap(err, "failed to create request")
    }

    req.Header.Set("clientId", client.ClientID)
    req.Header.Set("clientSecret", client.ClientSecret)

    resp, err := client.HttpClient.Do(req)
    if err != nil {
        return errors.Wrap(err, "request failed")
    }
    defer resp.Body.Close()

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return errors.Wrap(err, "failed to read response body")
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("authentication failed: %s", string(bodyBytes))
    }

    var data struct {
        Token        string `json:"token"`
        RefreshToken string `json:"refreshToken"`
    }

    if err := json.Unmarshal(bodyBytes, &data); err != nil {
        return errors.Wrap(err, "failed to parse response")
    }

    client.AccessToken = data.Token
    client.RefreshToken = data.RefreshToken
    client.TokenExpiresAt = client.getTokenExpiration(data.Token)

    return nil
}

// Refresh access token
func (client *WabotApiClient) RefreshTokenMethod() error {
    url := fmt.Sprintf("%s/refreshToken", client.ApiBaseURL)

    body := map[string]string{
        "refreshToken": client.RefreshToken,
    }
    bodyBytes, err := json.Marshal(body)
    if err != nil {
        return errors.Wrap(err, "failed to marshal request body")
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
    if err != nil {
        return errors.Wrap(err, "failed to create request")
    }

    req.Header.Set("clientId", client.ClientID)
    req.Header.Set("clientSecret", client.ClientSecret)
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.HttpClient.Do(req)
    if err != nil {
        return errors.Wrap(err, "request failed")
    }
    defer resp.Body.Close()

    respBodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return errors.Wrap(err, "failed to read response body")
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("token refresh failed: %s", string(respBodyBytes))
    }

    var data struct {
        Token        string `json:"token"`
        RefreshToken string `json:"refreshToken"`
    }

    if err := json.Unmarshal(respBodyBytes, &data); err != nil {
        return errors.Wrap(err, "failed to parse response")
    }

    client.AccessToken = data.Token
    client.RefreshToken = data.RefreshToken
    client.TokenExpiresAt = client.getTokenExpiration(data.Token)

    return nil
}

// Get templates
func (client *WabotApiClient) GetTemplates() ([]map[string]interface{}, error) {
    if err := client.ensureAuthenticated(); err != nil {
        return nil, err
    }

    url := fmt.Sprintf("%s/get-templates", client.ApiBaseURL)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, errors.Wrap(err, "failed to create request")
    }

    req.Header.Set("Authorization", client.AccessToken)

    resp, err := client.HttpClient.Do(req)
    if err != nil {
        return nil, errors.Wrap(err, "request failed")
    }
    defer resp.Body.Close()

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, errors.Wrap(err, "failed to read response body")
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get templates: %s", string(bodyBytes))
    }

    var data struct {
        Data []map[string]interface{} `json:"data"`
    }

    if err := json.Unmarshal(bodyBytes, &data); err != nil {
        return nil, errors.Wrap(err, "failed to parse response")
    }

    return data.Data, nil
}

// Send message
func (client *WabotApiClient) SendMessage(to string, templateID string, params []string) error {
    if err := client.ensureAuthenticated(); err != nil {
        return err
    }

    url := fmt.Sprintf("%s/send-message", client.ApiBaseURL)

    body := map[string]interface{}{
        "to":         to,
        "templateId": templateID,
        "params":     params,
    }
    bodyBytes, err := json.Marshal(body)
    if err != nil {
        return errors.Wrap(err, "failed to marshal request body")
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
    if err != nil {
        return errors.Wrap(err, "failed to create request")
    }

    req.Header.Set("Authorization", client.AccessToken)
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.HttpClient.Do(req)
    if err != nil {
        return errors.Wrap(err, "request failed")
    }
    defer resp.Body.Close()

    respBodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return errors.Wrap(err, "failed to read response body")
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send message: %s", string(respBodyBytes))
    }

    // Optionally, parse and return the response data
    // var responseData map[string]interface{}
    // if err := json.Unmarshal(respBodyBytes, &responseData); err != nil {
    //     return errors.Wrap(err, "failed to parse response")
    // }

    return nil
}

// Logout
func (client *WabotApiClient) Logout() error {
    url := fmt.Sprintf("%s/logout/%s", client.ApiBaseURL, client.RefreshToken)

    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return errors.Wrap(err, "failed to create request")
    }

    req.Header.Set("clientId", client.ClientID)
    req.Header.Set("clientSecret", client.ClientSecret)

    resp, err := client.HttpClient.Do(req)
    if err != nil {
        return errors.Wrap(err, "request failed")
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("logout failed: %s", string(bodyBytes))
    }

    // Clear tokens
    client.AccessToken = ""
    client.RefreshToken = ""
    client.TokenExpiresAt = time.Time{}

    return nil
}

// Utility methods

func (client *WabotApiClient) ensureAuthenticated() error {
    if client.AccessToken == "" || client.isTokenExpired() {
        if client.RefreshToken != "" {
            if err := client.RefreshTokenMethod(); err != nil {
                return client.Authenticate()
            }
        } else {
            return client.Authenticate()
        }
    }
    return nil
}

func (client *WabotApiClient) isTokenExpired() bool {
    return !client.TokenExpiresAt.IsZero() && time.Now().After(client.TokenExpiresAt)
}

func (client *WabotApiClient) getTokenExpiration(token string) time.Time {
    parser := &jwt.Parser{}

    tokenObj, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
    if err != nil {
        return time.Time{}
    }

    if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
        if exp, ok := claims["exp"].(float64); ok {
            return time.Unix(int64(exp), 0)
        }
    }

    return time.Time{}
}
