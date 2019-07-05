/*
Copyright 2017 The Kubernetes Authors.

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

package priorities

import (
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/algorithm"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

// NormalizeReduce generates a PriorityReduceFunction that can normalize the result
// scores to [0, maxPriority]. If reverse is set to true, it reverses the scores by
// subtracting it from maxPriority.
func NormalizeReduce(maxPriority int, reverse bool) algorithm.PriorityReduceFunction {
	return func(
		_ *v1.Pod,
		_ interface{},
		_ map[string]*schedulercache.NodeInfo,
		result schedulerapi.HostPriorityList) error { // 这个result是HostPriorityList，对应1个算法N个node的结果集

		var maxCount int
		for i := range result { // 遍历result将最高的Score赋值给maxCount
			if result[i].Score > maxCount {
				maxCount = result[i].Score
			}
		}

		if maxCount == 0 {
			if reverse {
				for i := range result {
					result[i].Score = maxPriority
				}
			}
			return nil
		}

		for i := range result {
			score := result[i].Score

			score = maxPriority * score / maxCount // maxPriority用来将得分限制在0-10
			if reverse {
				score = maxPriority - score
			}

			result[i].Score = score
		}
		return nil
	}
}
