package postgreSQL

import (
	"OzonTestTask/internal/model"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"time"
)

type Storage struct {
	db       *sqlx.DB
	squirrel squirrel.StatementBuilderType
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db:       db,
		squirrel: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// CreatePost Создание поста
func (s *Storage) CreatePost(ctx context.Context, post *model.Post) error {
	req, args, err := s.squirrel.
		Insert("posts").
		Columns("title", "content", "author", "are_comments_allowed", "created_at").
		Values(post.Title, post.Content, post.Author, post.AreCommentsAllowed, time.Now().UTC()).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return fmt.Errorf("ошибка построения SQL-запроса: %v", err)
	}

	err = s.db.QueryRowxContext(ctx, req, args...).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка при создании поста: %v", err)
	}
	return nil
}

func (s *Storage) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	req, args, err := s.squirrel.
		Select("id", "title", "content", "author", "are_comments_allowed", "created_at").
		From("posts").
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении постов: %v", err)
	}

	var posts []model.Post
	if err = s.db.SelectContext(ctx, &posts, req, args...); err != nil {
		return nil, fmt.Errorf("ошибка при получении постов: %v", err)
	}
	return posts, nil
}

func (s *Storage) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	req, args, err := s.squirrel.
		Select("id", "title", "content", "author", "are_comments_allowed", "created_at").
		From("posts").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении поста: %v", err)
	}

	var post model.Post
	if err = s.db.GetContext(ctx, &post, req, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пост не найден")
		}
		return nil, fmt.Errorf("ошибка при получении поста: %v", err)
	}
	return &post, nil
}

func (s *Storage) CreateComment(ctx context.Context, comment *model.Comment) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback()

	var allowed bool
	commentsAllowedReq, args, err := (s.squirrel.
		Select("are_comments_allowed").
		From("posts").
		Where(squirrel.Eq{"id": comment.PostID})).
		ToSql()

	if err = tx.GetContext(ctx, &allowed, commentsAllowedReq, args...); err != nil {
		return fmt.Errorf("пост не найден: %v", err)
	}

	// вставляю комментарий без path, чтобы получить id коммента и сформировать правильный путь
	req, args, err := s.squirrel.
		Insert("comments").
		Columns("post_id", "author", "content", "parent_comment_id", "path").
		Values(comment.PostID, comment.Author, comment.Content, comment.ParentCommentID, "").
		Suffix("RETURNING id, created_at").
		ToSql()

	if err = tx.QueryRowxContext(ctx, req, args...).Scan(&comment.ID, &comment.CreatedAt); err != nil {
		return fmt.Errorf("ошибка при вставке комментария: %v", err)
	}

	// не использую здесь squirrel, потому что работа с ltree
	// более читаема и удобна в написании с raw sql-запросом
	if comment.ParentCommentID != nil {
		rawReq := `
			UPDATE comments
			SET path = (SELECT path || text2ltree($1::text) FROM comments WHERE id = $2)
			WHERE id = $3
			RETURNING path`
		err = tx.QueryRowxContext(ctx, rawReq, comment.ID, *comment.ParentCommentID, comment.ID).Scan(&comment.Path)
	} else {
		rawReq := `UPDATE comments SET path = text2ltree($1::text) WHERE id = $2 RETURNING path`
		err = tx.QueryRowxContext(ctx, rawReq, comment.ID, comment.ID).Scan(&comment.Path)
	}
	if err != nil {
		return fmt.Errorf("ошибка при обновлении path: %v", err)
	}

	tx.Commit()
	return nil
}

func (s *Storage) GetCommentsByPost(ctx context.Context, postID, limit, offset int) ([]model.Comment, int, error) {
	req, args, err := s.squirrel.
		Select("id", "post_id", "author", "content", "parent_comment_id", "path::text AS path", "created_at").
		From("comments").
		Where(squirrel.Eq{"post_id": postID}).
		Where("parent_comment_id IS NULL").
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("ошибка формирования запроса на получение корневых комментариев: %v", err)
	}

	var comments []model.Comment
	if err = s.db.SelectContext(ctx, &comments, req, args...); err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении корневых комментариев: %v", err)
	}

	rootCommentsAmountReq, args, err := s.squirrel.
		Select("COUNT(*)").
		From("comments").
		Where(squirrel.Eq{"post_id": postID}).
		Where("parent_comment_id IS NULL").
		ToSql()

	var amount int
	if err = s.db.GetContext(ctx, &amount, rootCommentsAmountReq, args...); err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении количества всех корневых комментариев: %v", err)
	}

	return comments, amount, nil
}

func (s *Storage) GetReplies(ctx context.Context, parentID int) ([]model.Comment, error) {
	// не использую здесь squirrel, потому что работа с ltree
	// более читаема и удобна в написании с raw sql-запросом
	sqlStr := `
		SELECT c2.id, c2.post_id, c2.author, c2.content, c2.parent_comment_id, c2.path::text, c2.created_at
		FROM comments AS c1
		JOIN comments AS c2 ON c2.path <@ c1.path AND c2.id != c1.id
		WHERE c1.id = $1
		ORDER BY c2.path`

	var comments []model.Comment
	if err := s.db.SelectContext(ctx, &comments, sqlStr, parentID); err != nil {
		return nil, fmt.Errorf("ошибка при получении вложенных комментариев: %v", err)
	}

	return comments, nil
}
