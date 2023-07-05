package db

const (
	StmtCreateTable = `CREATE TABLE IF NOT EXISTS filesContent
                       (
                           path       text PRIMARY KEY,
                           content    jsonb,
                           created_at timestamp DEFAULT now(),
                           updated_at timestamp 
                       );`

	StmtInsert = `INSERT INTO filesContent (path, content) VALUES ($1, $2) 
                  ON CONFLICT (path) DO UPDATE 
                  SET content = EXCLUDED.content, 
                      updated_at = now();`
)
