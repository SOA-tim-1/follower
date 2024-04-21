package repo

import (
	"context"
	"follower/model"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository struct {
	Driver neo4j.DriverWithContext
	Logger *log.Logger
}

func (repo *UserRepository) CheckConnection() {
	ctx := context.Background()
	err := repo.Driver.VerifyConnectivity(ctx)
	if err != nil {
		repo.Logger.Panic(err)
		return
	}
	// repoint Neo4J server address
	repo.Logger.Printf(`Neo4J server address: %s`, repo.Driver.Target().Host)
}

func (repo *UserRepository) WriteUser(user *model.User) error {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	//savedUser, err := ...
	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (u:User) SET u.id = $id RETURN u.id",
				map[string]any{"id": user.ID})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		repo.Logger.Println("Error inserting Person:", err)
		return err
	}
	//repo.Logger.Println(savedUser.(int64))
	return nil
}

func (repo *UserRepository) FindById(id int64) (model.User, error) {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	savedUser, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (u:User {id: $id}) RETURN u.id",
				map[string]any{"id": id})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		repo.Logger.Println("Error inserting Person:", err)
		return model.User{}, err
	}
	//repo.Logger.Println(savedUser.(int64))
	foundUser := model.User{ID: savedUser.(int64)}
	return foundUser, nil
}

func (repo *UserRepository) CreateFollowConnection(firstId int64, secondId int64) error {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// ExecuteWrite for write transactions (Create/Update/Delete)
	//savedPerson,err := ...
	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (a:User), (b:User) WHERE a.id = $idOne AND b.id = $idTwo CREATE (a)-[r:FOLLOW]->(b) RETURN type(r)",
				map[string]any{"idOne": firstId, "idTwo": secondId})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		repo.Logger.Println("Error inserting Person:", err)
		return err
	}
	//repo.Logger.Println(savedPerson.(string))
	return nil
}

func (repo *UserRepository) GetFollows(id int64) (*[]int64, error) {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	var followedIDs []int64

	// ExecuteWrite for write transactions (Create/Update/Delete)
	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (interface{}, error) {
			result, err := transaction.Run(ctx,
				"MATCH (a:User)-[:FOLLOW]->(b:User) WHERE a.id = $id RETURN b.id",
				map[string]interface{}{"id": id})
			if err != nil {
				return nil, err
			}

			for result.Next(ctx) {
				record := result.Record()
				idValue, found := record.Get("b.id")
				if found {
					followedID, ok := idValue.(int64)
					if ok {
						followedIDs = append(followedIDs, followedID)
					}
				}
			}

			if err := result.Err(); err != nil {
				return nil, err
			}

			return nil, nil
		})
	if err != nil {
		repo.Logger.Println("Error fetching follows:", err)
		return nil, err
	}

	//repo.Logger.Println("Followed IDs:", followedIDs)

	return &followedIDs, nil
}

func (repo *UserRepository) GetFollowers(id int64) (*[]int64, error) {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	var followedIDs []int64

	// ExecuteWrite for write transactions (Create/Update/Delete)
	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (interface{}, error) {
			result, err := transaction.Run(ctx,
				"MATCH (a:User)-[:FOLLOW]->(b:User) WHERE b.id = $id RETURN a.id",
				map[string]interface{}{"id": id})
			if err != nil {
				return nil, err
			}

			for result.Next(ctx) {
				record := result.Record()
				idValue, found := record.Get("a.id")
				if found {
					followedID, ok := idValue.(int64)
					if ok {
						followedIDs = append(followedIDs, followedID)
					}
				}
			}

			if err := result.Err(); err != nil {
				return nil, err
			}

			return nil, nil
		})
	if err != nil {
		repo.Logger.Println("Error fetching follows:", err)
		return nil, err
	}

	//repo.Logger.Println("Followers IDs:", followedIDs)

	return &followedIDs, nil
}

func (repo *UserRepository) GetSuggestionsForUser(id int64) (*[]int64, error) {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	var suggestedIDs []int64

	// ExecuteWrite for write transactions (Create/Update/Delete)
	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (interface{}, error) {
			result, err := transaction.Run(ctx,
				"MATCH (a:User)-[:FOLLOW]->(b:User)-[:FOLLOW]->(c:User) WHERE a.id = $id RETURN c.id",
				map[string]interface{}{"id": id})
			if err != nil {
				return nil, err
			}

			for result.Next(ctx) {
				record := result.Record()
				idValue, found := record.Get("c.id")
				if found {
					followedID, ok := idValue.(int64)
					if ok {
						suggestedIDs = append(suggestedIDs, followedID)
					}
				}
			}

			if err := result.Err(); err != nil {
				return nil, err
			}

			return nil, nil
		})
	if err != nil {
		repo.Logger.Println("Error fetching follows:", err)
		return nil, err
	}

	invalidIDs, err := repo.GetFollows(id)

	invalidIDsSlice := *invalidIDs
	invalidIDsSlice = append(invalidIDsSlice, id)
	*invalidIDs = invalidIDsSlice

	suggestedIDs = removeDuplicates(suggestedIDs, *invalidIDs)

	if len(suggestedIDs) == 0 {
		suggestedIDsPointer, err := repo.GetRandomUsers()

		invalidIDs = &[]int64{}
		invalidIDsSlice := *invalidIDs
		invalidIDsSlice = append(invalidIDsSlice, id)
		*invalidIDs = invalidIDsSlice

		suggestedIDs = *suggestedIDsPointer

		suggestedIDs = removeDuplicates(suggestedIDs, *invalidIDs)

		if err != nil {
			return nil, err
		}

	}

	//repo.Logger.Println("Suggested IDs:", suggestedIDs)

	return &suggestedIDs, nil
}

func (repo UserRepository) GetRandomUsers() (*[]int64, error) {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	var randomIDs []int64

	// ExecuteWrite for write transactions (Create/Update/Delete)
	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (interface{}, error) {
			result, err := transaction.Run(ctx,
				"MATCH (u:User) RETURN u.id LIMIT 10",
				nil)
			if err != nil {
				return nil, err
			}

			for result.Next(ctx) {
				record := result.Record()
				idValue, found := record.Get("u.id")
				if found {
					randomID, ok := idValue.(int64)
					if ok {
						randomIDs = append(randomIDs, randomID)
					}
				}
			}

			if err := result.Err(); err != nil {
				return nil, err
			}

			return nil, nil
		})
	if err != nil {
		//repo.Logger.Println("Error fetching random users:", err)
		return nil, err
	}

	return &randomIDs, nil
}

func (repo *UserRepository) CloseDriverConnection(ctx context.Context) {
	repo.Driver.Close(ctx)
}

func (repo *UserRepository) DropAll() error {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (n) DETACH DELETE n",
				nil)
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		repo.Logger.Println("Error dropping daatbase:", err)
		return err
	}
	repo.Logger.Println("Dropped")
	return nil
}

func removeDuplicates(slice1, slice2 []int64) []int64 {
	// Create a map to store elements of slice2 for fast lookup
	set := make(map[int64]bool)
	for _, num := range slice2 {
		set[num] = true
	}

	// Iterate over slice1, keep only elements not present in slice2
	result := make([]int64, 0)
	for _, num := range slice1 {
		if !set[num] {
			result = append(result, num)
		}
	}
	return result
}
