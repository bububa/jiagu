package cluster

const (
	NOISE     = false
	CLUSTERED = true
)

// Dbscan DBSCAN 算法实现
type Dbscan struct {
	eps    float64
	minPts int
}

// NewDbscan 新建Dbscan
// eps: 邻域距离
// minPts: 核心对象中的最少样本数量
func NewDbscan(eps float64, minPts int) *Dbscan {
	return &Dbscan{
		eps:    eps,
		minPts: minPts,
	}
}

// Clusterize 输入数据，完成聚类
func (d *Dbscan) Clusterize(objects Observations) (Clusters, error) {
	var clusters Clusters
	visited := map[string]bool{}
	for _, point := range objects {
		neighbours := findNeighbours(point, objects, d.eps)
		if len(neighbours)+1 >= d.minPts {
			visited[point.GetID()] = CLUSTERED
			observations := make(Observations, 1)
			observations[0] = point
			observations = expandCluster(observations, neighbours, visited, d.minPts, d.eps)

			if len(observations) >= d.minPts {
				cluster := Cluster{
					Observations: observations,
				}
				clusters = append(clusters, cluster)
			}
		} else {
			visited[point.GetID()] = NOISE
		}
	}
	return clusters, nil
}

//Finds the neighbours from given array
//depends on Eps variable, which determines
//the distance limit from the point
func findNeighbours(point Observation, points Observations, eps float64) Observations {
	neighbours := make(Observations, 0)
	for _, potNeigb := range points {
		if point.GetID() != potNeigb.GetID() && potNeigb.Distance(point.Coordinates()) <= eps {
			neighbours = append(neighbours, potNeigb)
		}
	}
	return neighbours
}

//Try to expand existing clutser
func expandCluster(observations Observations, neighbours Observations, visited map[string]bool, minPts int, eps float64) Observations {
	seed := make(Observations, len(neighbours))
	copy(seed, neighbours)
	for _, point := range seed {
		pointState, isVisited := visited[point.GetID()]
		if !isVisited {
			currentNeighbours := findNeighbours(point, seed, eps)
			if len(currentNeighbours)+1 >= minPts {
				visited[point.GetID()] = CLUSTERED
				observations = merge(observations, currentNeighbours)
			}
		}

		if isVisited && pointState == NOISE {
			visited[point.GetID()] = CLUSTERED
			observations = append(observations, point)
		}
	}

	return observations
}

func merge(one Observations, two Observations) Observations {
	mergeMap := make(map[string]Observation)
	putAll(mergeMap, one)
	putAll(mergeMap, two)
	merged := make(Observations, 0)
	for _, val := range mergeMap {
		merged = append(merged, val)
	}

	return merged
}

//Function to add all values from list to map
//map keys is then the unique collecton from list
func putAll(m map[string]Observation, list Observations) {
	for _, val := range list {
		m[val.GetID()] = val
	}
}
