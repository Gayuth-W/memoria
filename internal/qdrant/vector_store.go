package vector

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

type VectorStore struct {
	client *qdrant.Client
}

func NewVectorStore() *VectorStore {
	client, _ := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6333,
	})

	return &VectorStore{client: client}
}

func (v *VectorStore) Init() error {
	return v.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "memories",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     768, // IMPORTANT: match your Ollama model output
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

func (v *VectorStore) Upsert(id string, vector []float32, payload map[string]any) error {

	_, err := v.client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "memories",
		Points: []*qdrant.PointStruct{
			{
				Id: &qdrant.PointId{
					PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
				},
				Vectors: &qdrant.Vectors{
					VectorsOptions: &qdrant.Vectors_Vector{
						Vector: &qdrant.Vector{Data: vector},
					},
				},
				Payload: payload,
			},
		},
	})

	return err
}

func (v *VectorStore) Search(vector []float32, limit uint64) ([]string, error) {

	res, err := v.client.Search(context.Background(), &qdrant.SearchPoints{
		CollectionName: "memories",
		Vector:         vector,
		Limit:          limit,
		WithPayload:    true,
	})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, r := range res {
		ids = append(ids, r.Id.GetUuid())
	}

	return ids, nil
}
