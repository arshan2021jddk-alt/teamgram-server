// Copyright 2022 Teamgram Authors
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: teamgramio (teamgram.io@gmail.com)
//

package me

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/teamgram/marmota/pkg/hack"
	"github.com/teamgram/teamgram-server/pkg/code/conf"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	_smsURL      = "https://api.sms.ir/v1/send/verify"
	_smsAPIKey   = "2vTYcfmR3PPybWQyXjz1wsL7tDvSzTrdkWE44d7eaXFeMl4r"
	_smsTemplate = "986819"
)

type verifySendParameterModel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type verifySendModel struct {
	Mobile     string                     `json:"mobile"`
	TemplateID int                        `json:"templateId"`
	Parameters []verifySendParameterModel `json:"parameters"`
}

func New(c *conf.SmsVerifyCodeConfig) *meVerifyCode {
	return &meVerifyCode{
		code: c,
	}
}

type meVerifyCode struct {
	code *conf.SmsVerifyCodeConfig
}

func (m *meVerifyCode) SendSmsVerifyCode(ctx context.Context, phoneNumber, code, codeHash string) (string, error) {
	apiKey := _smsAPIKey
	if m.code != nil && m.code.Key != "" {
		apiKey = m.code.Key
	}

	templateID := 986819
	if m.code != nil && m.code.Secret != "" {
		if _, err := fmt.Sscanf(m.code.Secret, "%d", &templateID); err != nil {
			logx.Errorf("invalid template id in config secret=%q: %v", m.code.Secret, err)
		}
	}

	bodyData, err := json.Marshal(&verifySendModel{
		Mobile:     phoneNumber,
		TemplateID: templateID,
		Parameters: []verifySendParameterModel{{Name: "CODE", Value: code}},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, _smsURL, bytes.NewReader(bodyData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	logx.Infof("send me sms via sms.ir: mobile=%s template=%d", phoneNumber, templateID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)
	body := hack.String(rawBody)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logx.Errorf("sms.ir request failed: status=%d body=%s", resp.StatusCode, body)
		return "", fmt.Errorf("sms.ir send failed: status=%d", resp.StatusCode)
	}

	logx.Infof("sms.ir response: %s", body)
	return "", nil
}

func (m *meVerifyCode) VerifySmsCode(ctx context.Context, codeHash, code, extraData string) error {
	if len(code) != 5 {
		return fmt.Errorf("code invalid")
	}

	//
	if code != extraData {
		return fmt.Errorf("code invalid")
	}

	// ...
	return nil
}
