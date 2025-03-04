/**
 * Tencent is pleased to support the open source community by making Polaris available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package config

import (
	"context"

	apiconfig "github.com/polarismesh/specification/source/go/api/v1/config_manage"

	"github.com/polarismesh/polaris/auth"
	"github.com/polarismesh/polaris/cache"
	"github.com/polarismesh/polaris/common/model"
	"github.com/polarismesh/polaris/namespace"
	"github.com/polarismesh/polaris/plugin"
	"github.com/polarismesh/polaris/store"
)

// Initialize 初始化配置中心模块
func TestInitialize(ctx context.Context, config Config, s store.Store, cacheMgn *cache.CacheManager,
	namespaceOperator namespace.NamespaceOperateServer, userMgn auth.UserServer,
	strategyMgn auth.StrategyServer) (ConfigCenterServer, ConfigCenterServer, error) {
	mockServer := &Server{}
	if err := mockServer.initialize(ctx, config, s, namespaceOperator, cacheMgn); err != nil {
		return nil, nil, err
	}
	return newServerAuthAbility(mockServer, userMgn, strategyMgn), mockServer, nil
}

func (s *Server) TestCheckClientConfigFile(ctx context.Context, files []*apiconfig.ClientConfigFileInfo,
	compartor compareFunction) (*apiconfig.ConfigClientResponse, bool) {
	return s.checkClientConfigFile(ctx, files, compartor)
}

func TestCompareByVersion(clientInfo *apiconfig.ClientConfigFileInfo, file *model.ConfigFileRelease) bool {
	return compareByVersion(clientInfo, file)
}

func TestCompareByMD5(clientInfo *apiconfig.ClientConfigFileInfo, file *model.ConfigFileRelease) bool {
	return compareByMD5(clientInfo, file)
}

// TestDecryptConfigFile 解密配置文件
func (s *Server) TestDecryptConfigFile(ctx context.Context, configFile *model.ConfigFile) (err error) {
	for i := range s.chains.chains {
		chain := s.chains.chains[i]
		if val, ok := chain.(*CryptoConfigFileChain); ok {
			if _, err := val.AfterGetFile(ctx, configFile); err != nil {
				return err
			}
		}
	}
	return nil
}

// TestEncryptConfigFile 解密配置文件
func (s *Server) TestEncryptConfigFile(ctx context.Context,
	configFile *model.ConfigFile, algorithm string, dataKey string) error {
	for i := range s.chains.chains {
		chain := s.chains.chains[i]
		if val, ok := chain.(*CryptoConfigFileChain); ok {
			return val.encryptConfigFile(ctx, configFile, algorithm, dataKey)
		}
	}
	return nil
}

// TestMockStore
func (s *Server) TestMockStore(ms store.Store) {
	s.storage = ms
}

// TestMockCryptoManager 获取加密管理
func (s *Server) TestMockCryptoManager(mgr plugin.CryptoManager) {
	s.cryptoManager = mgr
}
