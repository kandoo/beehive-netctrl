package discovery

import (
	"container/heap"
	"encoding/gob"
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
	if v, err := dict.Get(k); err == nil {
		links = v.(map[nom.UID][]nom.Link)
	}
	links[nt.UID()] = append(links[nt.UID()], link)
	return dict.Put(k, links)
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
		if v, err := dict.Get(string(nd.Node)); err == nil {
			nodeLinks = v.(map[nom.UID][]nom.Link)
		}
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

// LinksCentralized returns links of node.
//
// Note that this method should be used only when the GraphBuilderCentralized is
// in use.
func LinksCentralized(node nom.UID, ctx bh.RcvContext) (links []nom.Link) {
	dict := ctx.Dict(GraphDict)
	v, err := dict.Get(string(node))
	if err != nil {
		return nil
	}
	nodeLinks := v.(map[nom.UID][]nom.Link)

	for _, nl := range nodeLinks {
		links = append(links, nl...)
	}
	return links
}

// NodesCentralized returns the nodes with outgoing links so far.
//
// Note that this methods should be used only when the GraphBuilderCentralized
// is in use.
func NodesCentralized(ctx bh.RcvContext) (nodes []nom.UID) {
	ctx.Dict(GraphDict).ForEach(func(k string, v interface{}) bool {
		nodes = append(nodes, nom.UID(k))
		return true
	})
	return nodes
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

func init() {
	gob.Register(distAndLinks{})
	gob.Register(distMatrix{})
	gob.Register(GraphBuilderCentralized{})
	gob.Register(mentry{})
	gob.Register(nodeAndDist{})
	gob.Register(nodeAndDistSlice{})
}
