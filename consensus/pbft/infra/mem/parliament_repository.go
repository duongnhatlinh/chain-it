/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mem

import (
	"sync"

	"it-chain/consensus/pbft"
)

type ParliamentRepository struct {
	parliament pbft.Parliament
	sync.RWMutex
}

func NewParliamentRepositoryWithParliament(parliament pbft.Parliament) *ParliamentRepository {
	return &ParliamentRepository{
		parliament: parliament,
		RWMutex:    sync.RWMutex{},
	}
}

func NewParliamentRepository() *ParliamentRepository {
	return &ParliamentRepository{
		parliament: pbft.NewParliament(),
		RWMutex:    sync.RWMutex{},
	}
}

func (p *ParliamentRepository) Save(parliament pbft.Parliament) {

	p.Lock()
	defer p.Unlock()

	p.parliament = parliament
}

func (p *ParliamentRepository) Load() pbft.Parliament {

	p.Lock()
	defer p.Unlock()

	return p.parliament
}
