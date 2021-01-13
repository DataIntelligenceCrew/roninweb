package navigation

import (
	"math/rand"

	"github.com/DataIntelligenceCrew/OpenDataLink/internal/database"
	indexpkg "github.com/DataIntelligenceCrew/OpenDataLink/internal/index"
	"github.com/ekzhu/go-fasttext"
)

func (O *TableGraph) labelNodes(db *database.DB, ft *fasttext.FastText) error {
	idx, err := indexpkg.BuildCategoryEmbeddingIndex(db, ft)
	if err != nil {
		return err
	}

	for it := O.Nodes(); it.Next(); {
		node := it.Node().(*Node)
		names, _, err := idx.Query(node.vector, 1)
		if err != nil {
			return err
		}
		if O.isLeafNode(node) {
			var name string
			row := db.QueryRow("SELECT name FROM metadata WHERE dataset_id='" + node.name + "';")
			err := row.Scan(&name)
			if err == nil {
				node.name = name
			} else {
				println(err)
			}
		}
		if node.name == "" {
			token := make([]byte, 4)
			rand.Read(token)
			node.name = names[0] //+ " " + hex.EncodeToString(token)
		}
		println(node.name)
	}
	return nil
}
