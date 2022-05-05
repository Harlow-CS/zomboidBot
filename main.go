/*
Copyright © 2022 Jordan Harlow <harlowjordancs@gmail.com>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"

	"github.com/Harlow-CS/zomboidBot/bot"
	"github.com/Harlow-CS/zomboidBot/zomboid"
	"go.uber.org/zap"
)

func main() {

	fmt.Println("Hello, World!")

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	// sugar := logger.Sugar()

	// intialize server settings
	factorio.SetServerSettings()

	// start the bot
	bot.Start()

}
