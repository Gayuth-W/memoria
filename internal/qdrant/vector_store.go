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
		Port: 6334,
	})

	return &VectorStore{client: client}
}

func (v *VectorStore) Init() error {
	ctx := context.Background()
	exists, err := v.client.CollectionExists(ctx, "memories")
	if err != nil {
		return err
	}

	// If it exists, skip creation entirely to prevent crashing
	if exists {
		return nil
	}

	//Create the collection safely if it's missing
	return v.client.CreateCollection(ctx, &qdrant.CreateCollection{
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
				Payload: toQdrantPayload(payload),
			},
		},
	})

	return err
}

func (v *VectorStore) Search(vector []float32, limit uint64) ([]string, error) {

	res, err := v.client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "memories",
		Query:          qdrant.NewQuery(vector...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
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

func toQdrantPayload(m map[string]any) map[string]*qdrant.Value {
	out := make(map[string]*qdrant.Value)

	for k, v := range m {
		switch val := v.(type) {
		case string:
			out[k] = &qdrant.Value{
				Kind: &qdrant.Value_StringValue{StringValue: val},
			}
		case int:
			out[k] = &qdrant.Value{
				Kind: &qdrant.Value_IntegerValue{IntegerValue: int64(val)},
			}
		case int64:
			out[k] = &qdrant.Value{
				Kind: &qdrant.Value_IntegerValue{IntegerValue: val},
			}
		}
	}

	return out
}
