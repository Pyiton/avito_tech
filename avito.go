package tech_avito

import (
	"database/sql"
	"errors"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type MyModel struct {
	DB *sql.DB
}

// создает сегмент
func (m *MyModel) InsertSegment(SegmentName string) (int, error) {
	stmt := "INSERT INTO segments (segment_name) VALUES (?)"
	result, err := m.DB.Exec(stmt, SegmentName)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// создает юзера
func (m *MyModel) InsertUser(UserId int) (int, error) {
	stmt := "INSERT INTO users (id) VALUES (?)"
	result, err := m.DB.Exec(stmt, UserId)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// добавляет юзера в сегмент
func (m *MyModel) InsertUserAndSegment(UserId int, SegmentName string) (int, error) {
	stmt1 := "SELECT id FROM segments WHERE segment_name = ?"
	row := m.DB.QueryRow(stmt1, SegmentName)
	IdFind := 0
	err := row.Scan(&IdFind)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}
	stmt2 := "INSERT INTO users_and_segments (user_id, segment_id) VALUES (?, ?) "
	result, err := m.DB.Exec(stmt2, UserId, IdFind)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// добавляет юзера в множество сегментов
func (m *MyModel) InsertUserAndManySegments(UserId int, SegmentsNames []string) (int, error) {
	stmt1 := "SELECT id FROM users WHERE id = ?"
	row := m.DB.QueryRow(stmt1, UserId)
	IdFind := 0
	var err error
	err = row.Scan(&IdFind)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}
	var result int = 0
	for i := 0; i < len(SegmentsNames); i++ {
		result, err := m.InsertUserAndSegment(UserId, SegmentsNames[i])
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// удаляет сегмент
func (m *MyModel) DeleteSegment(SegmentName *string) (sql.Result, error) {
	stmt := "DELETE FROM segments WHERE segment_name = ?"
	result, err := m.DB.Exec(stmt, SegmentName)
	if err != nil {
		return result, err
	}
	return result, nil
}

// удаляет юзера
func (m *MyModel) DeleteUser(UserId int) (sql.Result, error) {
	stmt := "DELETE FROM users WHERE id = ?"
	result, err := m.DB.Exec(stmt, UserId)
	if err != nil {
		return result, err
	}
	return result, nil
}

// удаляет сегмент у юзера
func (m *MyModel) DeleteUserFromSegment(UserId int, SegmentName string) (sql.Result, error) {
	stmt1 := "SELECT id FROM segments WHERE segment_name = ?"
	row := m.DB.QueryRow(stmt1, SegmentName)
	IdFind := 0
	err := row.Scan(&IdFind)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	stmt := "DELETE FROM users_and_segments WHERE user_id = ? AND segment_id = ?"
	result, err := m.DB.Exec(stmt, UserId, IdFind)
	if err != nil {
		return result, err
	}
	return result, nil
}

// удаляет юзера из списка сегментов
func (m *MyModel) DeleteUserFromManySegments(UserId int, SegmentsNames []string) (sql.Result, error) {
	stmt1 := "SELECT id FROM users WHERE id = ?"
	row := m.DB.QueryRow(stmt1, UserId)
	IdFind := 0
	var err error
	err = row.Scan(&IdFind)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	var result sql.Result
	for i := 0; i < len(SegmentsNames); i++ {
		result, err = m.DeleteUserFromSegment(UserId, SegmentsNames[i])
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// находит все активные сегменты юзера
func (m *MyModel) FindUserSegments(UserId int) ([]string, error) {
	stmt := "SELECT segment_id FROM users_and_segments WHERE user_id = ?"
	rows, err := m.DB.Query(stmt, UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	stmt2 := "SELECT segment_name FROM segments WHERE id = ?"
	for rows.Next() {
		var SegIdNow int
		err = rows.Scan(&SegIdNow)
		if err != nil {
			return result, err
		}
		row := m.DB.QueryRow(stmt2, SegIdNow)
		var SegNameNow string
		err = row.Scan(&SegNameNow)
		if err != nil {
			return result, err
		}
		result = append(result, SegNameNow)
	}
	err = rows.Err()
	if err != nil {
		return result, err
	}
	return result, err
}
