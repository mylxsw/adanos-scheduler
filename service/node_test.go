package service_test

import (
	"context"
	"testing"

	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/adanos-scheduler/repo/mock"
	"github.com/mylxsw/adanos-scheduler/service"
	"github.com/mylxsw/container"
	"github.com/stretchr/testify/assert"
)

func createMockNodeService(ctx context.Context) service.NodeService {
	cc := container.NewWithContext(ctx)
	mock.Provider{}.Register(cc)

	return service.NewNodeService(cc)
}

func TestNodeService(t *testing.T) {

	srv := createMockNodeService(context.TODO())

	for _, node := range loadNodeData() {
		nodeID, err := srv.CreateNode(context.TODO(), node, "12")
		assert.NoError(t, err)
		assert.NotEmpty(t, nodeID)
	}

	nodes, err := srv.AllNodes(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, len(loadNodeData()), len(nodes))

	for _, nd := range nodes {
		node, err := srv.Node(context.TODO(), nd.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, node.Node.Name)
		assert.NotEmpty(t, node.CredentialIDs)
	}

	{
		selectedNodes, err := srv.SelectNodes(context.TODO(), repo.Labels{"deploy_app": "php", "priority": "medium"})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(selectedNodes))
	}

	{
		selectedNodes, err := srv.SelectNodes(context.TODO(), repo.Labels{"deploy_app": "php"})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(selectedNodes))
	}

	{
		selectedNodes, err := srv.SelectNodes(context.TODO(), repo.Labels{"priority": "high"})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(selectedNodes))
	}

	node0, err := srv.Node(context.TODO(), nodes[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, "node-1", node0.Node.Name)
	assert.Equal(t, "10.100.100.1", node0.Node.IP)
	assert.Equal(t, 22, node0.Node.Port)
	assert.Equal(t, "php", node0.Node.Labels["deploy_app"])
	assert.Equal(t, 1, len(node0.CredentialIDs))

	{
		nodeUpdate := node0.Node
		nodeUpdate.IP = "10.100.100.99"
		nodeUpdate.Port = 2222
		nodeUpdate.Labels["deploy_app"] = "java"
		err := srv.UpdateNode(context.TODO(), nodes[0].ID, nodeUpdate, "15", "17")
		assert.NoError(t, err)

		node0updated, err := srv.Node(context.TODO(), nodes[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, nodeUpdate.IP, node0updated.Node.IP)
		assert.Equal(t, nodeUpdate.Port, node0updated.Node.Port)
		assert.Equal(t, nodeUpdate.Labels["deploy_app"], node0updated.Node.Labels["deploy_app"])
		assert.Equal(t, 2, len(node0updated.CredentialIDs))
		assert.Contains(t, node0updated.CredentialIDs, "17")
		assert.Contains(t, node0updated.CredentialIDs, "15")
	}

	assert.NoError(t, srv.RemoveNode(context.TODO(), nodes[0].ID))
	nodes, err = srv.AllNodes(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, nodes, len(loadNodeData())-1)

}

func loadNodeData() []repo.Node {
	return []repo.Node{
		{
			Name: "node-1",
			IP:   "10.100.100.1",
			Port: 22,
			Type: repo.NodeTypeSSH,
			Labels: repo.Labels{
				"server":     "node-1",
				"deploy_app": "php",
			},
			Status: repo.NodeStatusEnabled,
		},
		{
			Name: "node-2",
			IP:   "10.100.100.2",
			Port: 22,
			Type: repo.NodeTypeSSH,
			Labels: repo.Labels{
				"server":     "node-2",
				"deploy_app": "java",
				"priority":   "high",
			},
			Status: repo.NodeStatusEnabled,
		},
		{
			Name: "node-3",
			IP:   "10.100.100.3",
			Port: 22,
			Type: repo.NodeTypeAgent,
			Labels: repo.Labels{
				"server":     "node-3",
				"deploy_app": "php",
				"priority":   "medium",
			},
			Status: repo.NodeStatusEnabled,
		},
		{
			Name: "node-4",
			IP:   "10.100.100.4",
			Port: 22,
			Type: repo.NodeTypeAgent,
			Labels: repo.Labels{
				"server":     "node-4",
				"deploy_app": "golang",
				"priority":   "high",
			},
			Status: repo.NodeStatusEnabled,
		},
	}
}
