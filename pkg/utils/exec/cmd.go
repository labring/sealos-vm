// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/labring/sealvm/pkg/utils/logger"
	strutil "github.com/labring/sealvm/pkg/utils/strings"
)

func Cmd(cmd string, args ...string) error {
	logger.Debug("cmd for pipe in host: ", fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
	cmder := exec.Command(cmd, args...)
	cmder.Stdout = os.Stdout
	cmder.Stderr = os.Stderr
	return cmder.Run()
}

func RunSimpleCmd(cmd string) (string, error) {
	logger.Debug("cmd for sh in host: ", cmd)
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	result, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput() // #nosec
	return string(result), err
}

func RunBashCmd(cmd string) (string, error) {
	logger.Debug("cmd for bash in host: ", cmd)
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	result, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput() // #nosec
	return string(result), err
}

func BashEval(cmd string) string {
	out, _ := RunBashCmd(cmd)
	return strutil.TrimWS(out)
}

func Eval(cmd string) string {
	out, _ := RunSimpleCmd(cmd)
	return strutil.TrimWS(out)
}

func CheckCmdIsExist(cmd string) (string, bool) {
	cmd = fmt.Sprintf("type %s", cmd)
	out, err := RunSimpleCmd(cmd)
	if err != nil {
		return "", false
	}

	outSlice := strings.Split(out, "is")
	last := outSlice[len(outSlice)-1]

	if last != "" && !strings.Contains(last, "not found") {
		return strings.TrimSpace(last), true
	}
	return "", false
}
