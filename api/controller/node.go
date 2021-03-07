package controller

import (
	"context"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/adanos-scheduler/service"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
)

type NodeController struct {
	cc container.Container
}

func NewNodeController(cc container.Container) web.Controller {
	return &NodeController{cc: cc}
}

func (ctl NodeController) Register(router *web.Router) {
	router.Group("/nodes/", func(router *web.Router) {
		router.Get("/", ctl.All)
		router.Get("/basic/", ctl.AllNodesBasic)

		router.Post("/", ctl.CreateNode)
		router.Get("/{node_id}/", ctl.Node)
		router.Put("/{node_id}/", ctl.UpdateNode)
	})
}

// All 返回所有的节点信息
func (ctl NodeController) All(ctx context.Context, nodeSrv service.NodeService) ([]repo.Node, error) {
	return nodeSrv.AllNodes(ctx)
}

// NodeBasicResp 节点查询相应（基础信息）
type NodeBasicResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// AllNodesBasic 返回所有节点的基本信息
func (ctl NodeController) AllNodesBasic(ctx context.Context, nodeSrv service.NodeService) ([]NodeBasicResp, error) {
	nodes, err := nodeSrv.AllNodes(ctx)
	if err != nil {
		return nil, err
	}

	var results []NodeBasicResp
	err = coll.MustNew(nodes).Map(func(node repo.Node) NodeBasicResp {
		return NodeBasicResp{
			ID:   node.ID,
			Name: node.Name,
			IP:   node.IP,
		}
	}).All(&results)
	return results, err
}

// NodeUpdateReq 节点更新创建请求
type NodeUpdateReq struct {
	Name          string            `json:"name"`
	Labels        map[string]string `json:"labels"`
	IP            string            `json:"ip"`
	Port          int               `json:"port"`
	Type          string            `json:"type"`
	Status        string            `json:"status"`
	CredentialIDs []string          `json:"credential_ids"`
}

func (nodeUpdateReq NodeUpdateReq) Validate(request web.Request) error {
	// TODO 表单校验
	return nil
}

func (nodeUpdateReq NodeUpdateReq) Transform() (repo.Node, []string) {
	return repo.Node{
		Name:   nodeUpdateReq.Name,
		Labels: nodeUpdateReq.Labels,
		IP:     nodeUpdateReq.IP,
		Port:   nodeUpdateReq.Port,
		Type:   repo.NodeType(nodeUpdateReq.Type),
		Status: repo.NodeStatus(nodeUpdateReq.Status),
	}, nodeUpdateReq.CredentialIDs
}

// CreateNode 创建一个新的服务器节点
func (ctl NodeController) CreateNode(ctx context.Context, req web.Request, nodeSrv service.NodeService) (IDResponse, error) {
	var nodeReq NodeUpdateReq
	if err := req.Unmarshal(&nodeReq); err != nil {
		return IDResponse{}, err
	}

	req.Validate(nodeReq, true)

	node, credentialIDs := nodeReq.Transform()
	nodeID, err := nodeSrv.CreateNode(ctx, node, credentialIDs...)
	return IDResponse{ID: nodeID}, err
}

// UpdateNode 更新一个新的服务器部署节点
func (ctl NodeController) UpdateNode(ctx context.Context, req web.Request, nodeSrv service.NodeService) error {
	var nodeReq NodeUpdateReq
	if err := req.Unmarshal(&nodeReq); err != nil {
		return err
	}

	req.Validate(nodeReq, true)

	node, credentialIDs := nodeReq.Transform()
	nodeID := req.PathVar("node_id")

	return nodeSrv.UpdateNode(ctx, nodeID, node, credentialIDs...)
}

// NodeDetailResp 节点详细信息响应
type NodeDetailResp struct {
	Node        repo.Node                  `json:"node"`
	Credentials []service.CredentialMasked `json:"credentials"`
}

// Node 返回节点详细信息，包含关联的凭据列表
func (ctl NodeController) Node(ctx context.Context, req web.Request, nodeSrv service.NodeService, credService service.CredentialService) (NodeDetailResp, error) {
	nodeID := req.PathVar("node_id")

	nodeWithCredIDs, err := nodeSrv.Node(ctx, nodeID)
	if err != nil {
		return NodeDetailResp{}, err
	}

	if len(nodeWithCredIDs.CredentialIDs) == 0 {
		return NodeDetailResp{Node: nodeWithCredIDs.Node, Credentials: nil}, nil
	}

	creds, err := credService.CredentialsMasked(ctx, nodeWithCredIDs.CredentialIDs...)
	if err != nil {
		log.With(nodeWithCredIDs).Errorf("query credentials failed: %v", err)
	}

	return NodeDetailResp{Node: nodeWithCredIDs.Node, Credentials: creds}, nil
}

// RemoveNode 删除一个服务器节点
func (ctl NodeController) RemoveNode(ctx context.Context, req web.Request, nodeSrv service.NodeService) error {
	nodeID := req.PathVar("node_id")
	return nodeSrv.RemoveNode(ctx, nodeID)
}
