package discovery

import (
	"container/heap"
	"fmt"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

const (
	GraphDict = "NetGraph"
)

// GraphBuilderCentralized is a handler that builds a centralized graph of the
// network topology. This handler is only useful for centralized applications.
type GraphBuilderCentralized struct{}

func (b GraphBuilderCentralized) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	dict := ctx.Dict(GraphDict)
	var link nom.Link
	switch dm := msg.Data().(type) {
	case nom.LinkAdded:
		link = nom.Link(dm)
	case nom.LinkDeleted:
		link = nom.Link(dm)
	default:
		return fmt.Errorf("GraphBuilderCentralized: unsupported message type %v",
			msg.Type())
	}

	nf, _ := nom.ParsePortUID(link.From)
	nt, _ := nom.ParsePortUID(link.To)

	if nf == nt {
		return fmt.Errorf("%v is a loop", link)
	}

	k := string(nf)
	links := make(map[nom.UID][]nom.Link)
	dict.GetGob(k, &links)
	links[nt.UID()] = append(links[nt.UID()], link)
	return dict.PutGob(k, &links)
}

func (b GraphBuilderCentralized) Map(msg bh.Msg,
	ctx bh.MapContext) bh.MappedCells {

	var from nom.UID
	switch dm := msg.Data().(type) {
	case nom.LinkAdded:
		from = dm.From
	case nom.LinkDeleted:
		from = dm.From
	default:
		return nil
	}
	// TODO(soheil): maybe store and update the matrix directly here.
	n, _ := nom.ParsePortUID(from)
	return bh.MappedCells{{GraphDict, string(n)}}
}

// ShortestPathCentralized calculates the shortest path from node "from" to node
// "to" according to the state stored in GraphDict by the
// GraphBuilderCentralized.
//
// This method is not go-routine safe and must be called within a handler of the
// application that uses the GraphBuilderCentralized as a handler. Otherwise,
// the user needs to synchronize the two.
func ShortestPathCentralized(from, to nom.UID, ctx bh.RcvContext) (
	paths [][]nom.Link, length int) {

	if from == to {
		return nil, 0
	}

	visited := make(map[nom.UID]distAndLinks)
	visited[from] = distAndLinks{Dist: 0}

	pq := nodeAndDistSlice{{Dist: 0, Node: from}}
	heap.Init(&pq)

	dict := ctx.Dict(GraphDict)
	for len(pq) != 0 {
		nd := heap.Pop(&pq).(nodeAndDist)
		if nd.Node == to {
			continue
		}
		nodeLinks := make(map[nom.UID][]nom.Link)
		dict.GetGob(string(nd.Node), &nodeLinks)
		nd.Dist = visited[nd.Node].Dist
		for _, links := range nodeLinks {
			for _, l := range links {
				nid, _ := nom.ParsePortUID(l.To)
				ton := nom.UID(nid)
				if dl, ok := visited[ton]; ok {
					switch {
					case nd.Dist+1 < dl.Dist:
						glog.Fatalf("invalid distance in BFS")
					case nd.Dist+1 == dl.Dist:
						dl.BackLinks = append(dl.BackLinks, l)
						visited[ton] = dl
					}
					continue
				}

				visited[ton] = distAndLinks{
					Dist:      nd.Dist + 1,
					BackLinks: []nom.Link{l},
				}
				ndto := nodeAndDist{
					Dist: nd.Dist + 1,
					Node: ton,
				}
				heap.Push(&pq, ndto)
			}
		}
	}
	return allPaths(from, to, visited)
}

func allPaths(from, to nom.UID, visited map[nom.UID]distAndLinks) (
	[][]nom.Link, int) {

	if from == to {
		return nil, 0
	}

	todl, ok := visited[to]
	if !ok {
		return nil, -1
	}

	var paths [][]nom.Link
	for _, l := range todl.BackLinks {
		lfn, _ := nom.ParsePortUID(l.From)
		prevpaths, _ := allPaths(from, nom.UID(lfn), visited)
		if len(prevpaths) == 0 {
			paths = append(paths, []nom.Link{l})
			continue
		}
		for _, p := range prevpaths {
			paths = append(paths, append(p, l))
		}
	}
	return paths, todl.Dist
}

type nodeAndDist struct {
	Dist int
	Node nom.UID
}

type nodeAndDistSlice []nodeAndDist

func (s nodeAndDistSlice) Len() int           { return len(s) }
func (s nodeAndDistSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s nodeAndDistSlice) Less(i, j int) bool { return s[i].Dist < s[j].Dist }

func (s *nodeAndDistSlice) Push(x interface{}) {
	*s = append(*s, x.(nodeAndDist))
}

func (s *nodeAndDistSlice) Pop() interface{} {
	l := len(*s) - 1
	nd := (*s)[l]
	*s = (*s)[0:l]
	return nd
}

type distMatrix map[mentry]distAndLinks

type mentry struct {
	From nom.UID
	To   nom.UID
}

type distAndLinks struct {
	Dist      int
	BackLinks []nom.Link
}

//func calcSPCentralized(from, to nom.UID, dict state.Dict, length int,
//blacklist map[nom.UID]int) [][]nom.UID {

//var paths [][]nom.UID
//links := make(map[nom.UID]nom.Link)
//dict.GetGob(string(from), &links)
//sp := math.MaxInt64
//for _, l := range links {
//if prevlen, ok := blacklist[l.To]; ok && prevlen < length+1 {
//continue
//}

//if l.To == to {
//if sp != 1 {
//paths = [][]nom.UID{[]nom.UID{l}}
//sp = 1
//} else {
//paths = append(paths, []nom.UID{l})
//}
//blackList[l.To] = length + 1
//continue
//}

//blackList[l.To]
//calcSPCentralized(l.To, to, dict, blacklist)
//}
//return paths
//}
